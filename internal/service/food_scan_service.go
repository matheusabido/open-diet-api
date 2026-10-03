package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/matheusabido/open-diet-api/internal/dto"
	"github.com/matheusabido/open-diet-api/internal/groq"
	"github.com/matheusabido/open-diet-api/internal/openfoodfacts"
	"github.com/matheusabido/open-diet-api/internal/unit"
)

const (
	maxNutrient = 100000
	maxAmount   = 100000
)

var ErrLabelNotFound = errors.New("nutrition label not found")

const labelPrompt = `You read the nutrition label ("Informação Nutricional") of a Brazilian food package from a photo. The label is either a table or a running text.

A label often shows two sets of values at the same time: one per 100 g or 100 ml and one per serving ("porção"). In a table they are side by side columns. In a running text they follow each other for every nutrient, one in regular text and the other in bold or inside parentheses, as in "Carboidratos 57 g (24 g, 8%)". The header says which is which, as in "Por 100 g (41,5 g, %VD*)". Never mix the sets: all the values of a set come from the same column or the same text style. The %VD numbers are percentages, never use them.

Answer only with a JSON object in this format:
{"found": true, "serving": {"amount": 30, "unit": "g", "kcal": 102, "carbohydrate": 17, "protein": 4.8, "fat": 1.8}, "per_100": {"amount": 100, "unit": "g", "kcal": 340, "carbohydrate": 57, "protein": 16, "fat": 6}}

- "serving": the declared serving and the values for it. "amount" and "unit" are the metric size of the serving, "Porção de 30 g (2 colheres de sopa)" is amount 30 and unit "g". Never use the household measure. null when the label has no values per serving.
- "per_100": the values per 100 g or 100 ml. null when the label does not have them. Never calculate a set from the other.
- "unit": "g" or "ml".
- "kcal": the energy value ("Valor energético") in kcal, never in kJ.
- "carbohydrate": the total carbohydrates ("Carboidratos"), never the sugars.
- "protein": "Proteínas".
- "fat": the total fat ("Gorduras totais"), never the saturated or trans fat.
- Numbers use a decimal point: "4,8 g" is 4.8. A nutrient declared as zero or "não contém" is 0. A value that cannot be read is null.
- "found": false when the photo has no nutrition label.`

type labelValues struct {
	Amount       *float64 `json:"amount"`
	Unit         string   `json:"unit"`
	Kcal         *float64 `json:"kcal"`
	Carbohydrate *float64 `json:"carbohydrate"`
	Protein      *float64 `json:"protein"`
	Fat          *float64 `json:"fat"`
}

type label struct {
	Found   bool         `json:"found"`
	Serving *labelValues `json:"serving"`
	Per100  *labelValues `json:"per_100"`
}

type FoodScanService struct {
	groq *groq.Client
	off  *openfoodfacts.Client
}

func NewFoodScanService(groq *groq.Client, off *openfoodfacts.Client) *FoodScanService {
	return &FoodScanService{groq: groq, off: off}
}

func (s *FoodScanService) FromPhoto(ctx context.Context, photo io.Reader) (*dto.FoodDraftDTO, error) {
	image, err := io.ReadAll(photo)
	if err != nil {
		return nil, err
	}

	answer, err := s.groq.Vision(ctx, labelPrompt, image)
	if err != nil {
		return nil, err
	}

	var read label
	if err := json.Unmarshal([]byte(answer), &read); err != nil {
		return nil, fmt.Errorf("read label: %w: %s", err, answer)
	}

	for _, values := range []*labelValues{read.Per100, read.Serving} {
		if !read.Found || values == nil {
			continue
		}

		itemUnit, ok := unit.Parse(values.Unit)
		amount := valid(values.Amount, maxAmount)
		if !ok || amount == nil || *amount == 0 {
			continue
		}

		return &dto.FoodDraftDTO{
			Kcal:         valid(values.Kcal, maxNutrient),
			Protein:      valid(values.Protein, maxNutrient),
			Carbohydrate: valid(values.Carbohydrate, maxNutrient),
			Fat:          valid(values.Fat, maxNutrient),
			Amount:       amount,
			Unit:         itemUnit,
		}, nil
	}

	return nil, ErrLabelNotFound
}

func (s *FoodScanService) FromBarcode(ctx context.Context, barcode string) (*dto.FoodDraftDTO, error) {
	product, err := s.off.Product(ctx, barcode)
	if err != nil {
		return nil, err
	}

	food := &dto.FoodDraftDTO{Name: strings.TrimSpace(product.NamePt), Barcode: &barcode}
	if food.Name == "" {
		food.Name = strings.TrimSpace(product.Name)
	}
	if brand, _, _ := strings.Cut(product.Brands, ","); strings.TrimSpace(brand) != "" {
		brand = strings.TrimSpace(brand)
		food.Brand = &brand
	}

	suffix, amount, itemUnit := "_100g", 100.0, unit.G
	if product.ServingQuantityUnit == "ml" || product.ProductQuantityUnit == "ml" || product.NutritionDataPer == "100ml" {
		itemUnit = unit.ML
	}

	serving, _ := product.ServingQuantity.Float64()
	servingUnit, ok := unit.Parse(product.ServingQuantityUnit)
	hasServing := ok && serving > 0 && serving <= maxAmount && product.Nutriment("energy-kcal_serving") != nil
	if product.Nutriment("energy-kcal_100g") == nil && hasServing {
		suffix, amount, itemUnit = "_serving", serving, servingUnit
	}

	food.Amount = &amount
	food.Unit = itemUnit
	food.Kcal = round(valid(product.Nutriment("energy-kcal"+suffix), maxNutrient))
	food.Protein = round(valid(product.Nutriment("proteins"+suffix), maxNutrient))
	food.Carbohydrate = round(valid(product.Nutriment("carbohydrates"+suffix), maxNutrient))
	food.Fat = round(valid(product.Nutriment("fat"+suffix), maxNutrient))

	return food, nil
}

func valid(value *float64, max float64) *float64 {
	if value == nil || *value < 0 || *value > max {
		return nil
	}
	return value
}

func round(value *float64) *float64 {
	if value == nil {
		return nil
	}
	rounded := math.Round(*value*10) / 10
	return &rounded
}
