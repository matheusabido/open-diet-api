package voice

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/matheusabido/open-diet-api/internal/unit"
	"golang.org/x/text/unicode/norm"
)

var (
	amountPattern    = regexp.MustCompile(`^(\d+(?:[.,]\d+)*)(\p{L}*)$`)
	thousandsPattern = regexp.MustCompile(`^\d{1,3}(\.\d{3})+$`)
)

var units = map[string]unit.Unit{
	"g": unit.G, "gr": unit.G, "grama": unit.G, "gramas": unit.G,
	"kg": unit.KG, "quilo": unit.KG, "quilos": unit.KG, "kilo": unit.KG, "kilos": unit.KG,
	"quilograma": unit.KG, "quilogramas": unit.KG,
	"mg": unit.MG, "miligrama": unit.MG, "miligramas": unit.MG,
	"oz": unit.OZ, "onca": unit.OZ, "oncas": unit.OZ,
	"lb": unit.LB, "libra": unit.LB, "libras": unit.LB,
	"ml": unit.ML, "mililitro": unit.ML, "mililitros": unit.ML,
	"l": unit.L, "litro": unit.L, "litros": unit.L,
	"dl": unit.DL, "decilitro": unit.DL, "decilitros": unit.DL,
	"un": unit.UN, "unidade": unit.UN, "unidades": unit.UN,
}

var numbers = map[string]float64{
	"meio": 0.5, "meia": 0.5, "½": 0.5, "1/2": 0.5,
	"um": 1, "uma": 1, "dois": 2, "duas": 2, "tres": 3, "quatro": 4, "cinco": 5,
	"seis": 6, "sete": 7, "oito": 8, "nove": 9, "dez": 10, "onze": 11, "doze": 12, "treze": 13,
	"catorze": 14, "quatorze": 14, "quinze": 15, "dezesseis": 16, "dezessete": 17, "dezoito": 18,
	"dezenove": 19, "vinte": 20, "trinta": 30, "quarenta": 40, "cinquenta": 50, "sessenta": 60,
	"setenta": 70, "oitenta": 80, "noventa": 90, "cem": 100, "cento": 100,
	"duzentos": 200, "duzentas": 200, "trezentos": 300, "trezentas": 300,
	"quatrocentos": 400, "quatrocentas": 400, "quinhentos": 500, "quinhentas": 500,
	"seiscentos": 600, "seiscentas": 600, "setecentos": 700, "setecentas": 700,
	"oitocentos": 800, "oitocentas": 800, "novecentos": 900, "novecentas": 900,
}

var prepositions = map[string]bool{"de": true, "do": true, "da": true, "dos": true, "das": true}

func Parse(text string) (amount float64, u unit.Unit, name string, ok bool) {
	words := tokenize(text)

	amount, suffix, words := parseAmount(words)
	if amount <= 0 {
		return 0, "", "", false
	}

	u = unit.UN
	if suffix != "" {
		if u, ok = units[suffix]; !ok {
			return 0, "", "", false
		}
	} else if len(words) > 0 {
		if spoken, ok := units[words[0]]; ok {
			u = spoken
			words = words[1:]
		}
	}

	if len(words) > 0 && prepositions[words[0]] {
		words = words[1:]
	}
	if len(words) == 0 {
		return 0, "", "", false
	}

	return amount, u, strings.Join(words, " "), true
}

func Normalize(text string) string {
	var b strings.Builder
	for _, r := range norm.NFD.String(strings.ToLower(text)) {
		if !unicode.Is(unicode.Mn, r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func tokenize(text string) []string {
	var words []string
	for _, word := range strings.Fields(Normalize(text)) {
		if word = strings.Trim(word, ".,!?;:\"'()[]"); word != "" {
			words = append(words, word)
		}
	}
	return words
}

func parseAmount(words []string) (amount float64, suffix string, rest []string) {
	if len(words) == 0 {
		return 0, "", nil
	}

	if match := amountPattern.FindStringSubmatch(words[0]); match != nil {
		digits := match[1]
		if thousandsPattern.MatchString(digits) {
			digits = strings.ReplaceAll(digits, ".", "")
		}
		amount, err := strconv.ParseFloat(strings.ReplaceAll(digits, ",", "."), 64)
		if err != nil {
			return 0, "", nil
		}
		return amount, match[2], words[1:]
	}

	var thousands, current float64
	i := 0
	for i < len(words) {
		word := words[i]
		if word == "e" && i > 0 && i+1 < len(words) && isNumber(words[i+1]) {
			i++
			continue
		}
		if word == "mil" {
			thousands = max(current, 1) * 1000
			current = 0
		} else if value, ok := numbers[word]; ok {
			current += value
		} else {
			break
		}
		i++
	}

	return thousands + current, "", words[i:]
}

func isNumber(word string) bool {
	_, ok := numbers[word]
	return ok || word == "mil"
}
