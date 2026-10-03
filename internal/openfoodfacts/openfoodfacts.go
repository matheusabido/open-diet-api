package openfoodfacts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	endpoint  = "https://world.openfoodfacts.org/api/v2/product/%s.json?fields=" + fields
	fields    = "product_name,product_name_pt,brands,nutriments,serving_quantity,serving_quantity_unit,product_quantity_unit,nutrition_data_per"
	userAgent = "open-diet-api/1.0"
	timeout   = 10 * time.Second
)

var ErrProductNotFound = errors.New("product not found")

type Product struct {
	Name                string         `json:"product_name"`
	NamePt              string         `json:"product_name_pt"`
	Brands              string         `json:"brands"`
	Nutriments          map[string]any `json:"nutriments"`
	ServingQuantity     json.Number    `json:"serving_quantity"`
	ServingQuantityUnit string         `json:"serving_quantity_unit"`
	ProductQuantityUnit string         `json:"product_quantity_unit"`
	NutritionDataPer    string         `json:"nutrition_data_per"`
}

// Nutriment returns a value such as "proteins_100g" or "energy-kcal_serving".
// Open Food Facts sends the numbers either as numbers or as strings.
func (p *Product) Nutriment(key string) *float64 {
	var value float64
	switch raw := p.Nutriments[key].(type) {
	case float64:
		value = raw
	case string:
		if _, err := fmt.Sscan(raw, &value); err != nil {
			return nil
		}
	default:
		return nil
	}
	return &value
}

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{http: &http.Client{Timeout: timeout}}
}

func (c *Client) Product(ctx context.Context, barcode string) (*Product, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(endpoint, barcode), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("open food facts: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrProductNotFound
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("open food facts: status %d", res.StatusCode)
	}

	var body struct {
		Status  int      `json:"status"`
		Product *Product `json:"product"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("open food facts: %w", err)
	}
	if body.Status != 1 || body.Product == nil {
		return nil, ErrProductNotFound
	}

	return body.Product, nil
}
