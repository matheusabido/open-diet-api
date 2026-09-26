package i18n

import "fmt"

type Key string

const (
	Required             Key = "required"
	MinLength            Key = "min_length"
	MaxLength            Key = "max_length"
	Between              Key = "between"
	UsernameCharset      Key = "username_charset"
	PasswordCharset      Key = "password_charset"
	PasswordConfirmation Key = "password_confirmation"
	InvalidBody          Key = "invalid_body"
	UsernameTaken        Key = "username_taken"
	InvalidCredentials   Key = "invalid_credentials"
	InvalidRefreshToken  Key = "invalid_refresh_token"
	Unauthorized         Key = "unauthorized"
	UserNotFound         Key = "user_not_found"
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
		Between:              "deve estar entre %d e %d",
		UsernameCharset:      "deve conter apenas letras, números, '_' e '.'",
		PasswordCharset:      "deve conter apenas letras, números e os caracteres especiais %s",
		PasswordConfirmation: "deve ser igual à senha",
		InvalidBody:          "corpo da requisição inválido",
		UsernameTaken:        "nome de usuário já está em uso",
		InvalidCredentials:   "usuário ou senha inválidos",
		InvalidRefreshToken:  "refresh token inválido ou expirado",
		Unauthorized:         "não autenticado",
		UserNotFound:         "usuário não encontrado",
		InternalError:        "erro interno do servidor",
	},
	En: {
		Required:             "is required",
		MinLength:            "must be at least %d characters",
		MaxLength:            "must be at most %d characters",
		Between:              "must be between %d and %d",
		UsernameCharset:      "may only contain letters, numbers, '_' and '.'",
		PasswordCharset:      "may only contain letters, numbers and the special characters %s",
		PasswordConfirmation: "must match password",
		InvalidBody:          "invalid request body",
		UsernameTaken:        "username already taken",
		InvalidCredentials:   "invalid username or password",
		InvalidRefreshToken:  "invalid or expired refresh token",
		Unauthorized:         "unauthenticated",
		UserNotFound:         "user not found",
		InternalError:        "internal server error",
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
