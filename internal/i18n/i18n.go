package i18n

import "fmt"

type Key string

const (
	Required             Key = "required"
	MinLength            Key = "min_length"
	MaxLength            Key = "max_length"
	Between              Key = "between"
	GreaterThan          Key = "greater_than"
	MaxValue             Key = "max_value"
	InvalidUnit          Key = "invalid_unit"
	IncompatibleUnit     Key = "incompatible_unit"
	InvalidISODate       Key = "invalid_iso_date"
	InvalidMeal          Key = "invalid_meal"
	UsernameCharset      Key = "username_charset"
	PasswordCharset      Key = "password_charset"
	PasswordConfirmation Key = "password_confirmation"
	InvalidBody          Key = "invalid_body"
	InvalidQuery         Key = "invalid_query"
	UsernameTaken        Key = "username_taken"
	InvalidCredentials   Key = "invalid_credentials"
	InvalidRefreshToken  Key = "invalid_refresh_token"
	Unauthorized         Key = "unauthorized"
	Forbidden            Key = "forbidden"
	FoodNotFound         Key = "food_not_found"
	FoodItemNotFound     Key = "food_item_not_found"
	InternalError        Key = "internal_error"
)

const (
	PtBR = "pt-BR"
	En   = "en"

	Default = PtBR
)

var Supported = []string{PtBR, En}

var messages = map[string]map[Key]string{
	PtBR: {
		Required:             "é obrigatório",
		MinLength:            "deve ter pelo menos %d caracteres",
		MaxLength:            "deve ter no máximo %d caracteres",
		Between:              "deve estar entre %v e %v",
		GreaterThan:          "deve ser maior que %v",
		MaxValue:             "deve ser no máximo %v",
		InvalidUnit:          "deve ser uma das unidades: %s",
		IncompatibleUnit:     "deve ser compatível com a unidade da comida: %s",
		InvalidISODate:       "deve ser uma data ISO 8601 (ex: 2026-01-31)",
		InvalidMeal:          "deve ser uma das refeições: %s",
		UsernameCharset:      "deve conter apenas letras, números, '_' e '.'",
		PasswordCharset:      "deve conter apenas letras, números e os caracteres especiais %s",
		PasswordConfirmation: "deve ser igual à senha",
		InvalidBody:          "corpo da requisição inválido",
		InvalidQuery:         "parâmetros de consulta inválidos",
		UsernameTaken:        "nome de usuário já está em uso",
		InvalidCredentials:   "usuário ou senha inválidos",
		InvalidRefreshToken:  "refresh token inválido ou expirado",
		Unauthorized:         "não autenticado",
		Forbidden:            "acesso negado",
		FoodNotFound:         "comida não encontrada",
		FoodItemNotFound:     "item não encontrado",
		InternalError:        "Ocorreu um erro interno no servidor. Por favor, tente novamente.",
	},
	En: {
		Required:             "is required",
		MinLength:            "must be at least %d characters",
		MaxLength:            "must be at most %d characters",
		Between:              "must be between %v and %v",
		GreaterThan:          "must be greater than %v",
		MaxValue:             "must be at most %v",
		InvalidUnit:          "must be one of the units: %s",
		IncompatibleUnit:     "must be compatible with the food unit: %s",
		InvalidISODate:       "must be an ISO 8601 date (e.g. 2026-01-31)",
		InvalidMeal:          "must be one of the meals: %s",
		UsernameCharset:      "may only contain letters, numbers, '_' and '.'",
		PasswordCharset:      "may only contain letters, numbers and the special characters %s",
		PasswordConfirmation: "must match password",
		InvalidBody:          "invalid request body",
		InvalidQuery:         "invalid query parameters",
		UsernameTaken:        "username already taken",
		InvalidCredentials:   "invalid username or password",
		InvalidRefreshToken:  "invalid or expired refresh token",
		Unauthorized:         "unauthenticated",
		Forbidden:            "access denied",
		FoodNotFound:         "food not found",
		FoodItemNotFound:     "item not found",
		InternalError:        "An internal server error occurred. Please try again.",
	},
}

type Translator func(key Key, args ...any) string

func New(lang string) Translator {
	catalog, ok := messages[lang]
	if !ok {
		catalog = messages[Default]
	}

	return func(key Key, args ...any) string {
		message, ok := catalog[key]
		if !ok {
			message, ok = messages[Default][key]
		}
		if !ok {
			return string(key)
		}
		if len(args) == 0 {
			return message
		}
		return fmt.Sprintf(message, args...)
	}
}
