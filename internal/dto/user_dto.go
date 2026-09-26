package dto

import (
	"regexp"

	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/model"
	"github.com/matheusabido/open-diet-api/internal/validation"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+=.,?-]+$`)
)

const passwordSpecialChars = "!@#$%^&*()_+-=.,?"

type UserRegisterDTO struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (d UserRegisterDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	v.Check(validation.NotBlank(d.Username), "username", t(i18n.Required))
	v.Check(validation.MinLength(d.Username, 3), "username", t(i18n.MinLength, 3))
	v.Check(validation.MaxLength(d.Username, 32), "username", t(i18n.MaxLength, 32))
	v.Check(validation.Matches(d.Username, usernameRegex), "username", t(i18n.UsernameCharset))

	v.Check(validation.NotBlank(d.Password), "password", t(i18n.Required))
	v.Check(validation.MinLength(d.Password, 8), "password", t(i18n.MinLength, 8))
	v.Check(validation.MaxLength(d.Password, 72), "password", t(i18n.MaxLength, 72))
	v.Check(validation.Matches(d.Password, passwordRegex), "password", t(i18n.PasswordCharset, passwordSpecialChars))

	v.Check(validation.Equal(d.Password, d.PasswordConfirmation), "password_confirmation", t(i18n.PasswordConfirmation))

	return v.Errors()
}

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (d UserLoginDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	v.Check(validation.NotBlank(d.Username), "username", t(i18n.Required))
	v.Check(validation.NotBlank(d.Password), "password", t(i18n.Required))

	return v.Errors()
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"`
}

func (d RefreshTokenDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	v.Check(validation.NotBlank(d.RefreshToken), "refresh_token", t(i18n.Required))

	return v.Errors()
}

type AuthResponseDTO struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	User         *model.User `json:"user"`
}
