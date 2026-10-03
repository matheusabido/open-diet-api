package dto

import (
	"math"
	"regexp"
	"strings"

	"github.com/matheusabido/open-diet-api/internal/i18n"
	"github.com/matheusabido/open-diet-api/internal/model"
	"github.com/matheusabido/open-diet-api/internal/validation"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+=.,?-]+$`)
)

const passwordSpecialChars = "!@#$%^&*()_+-=.,?"

const (
	minAge = 1
	maxAge = 120

	minWeight = 20.0
	maxWeight = 500.0

	minHeight = 50
	maxHeight = 300
)

type UserRegisterDTO struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (d UserRegisterDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	v.Check(validation.NotBlank(d.Username), "username", t(i18n.Required))
	v.Check(validation.MinLength(d.Username, 5), "username", t(i18n.MinLength, 5))
	v.Check(validation.MaxLength(d.Username, 32), "username", t(i18n.MaxLength, 32))
	v.Check(validation.Matches(d.Username, usernameRegex), "username", t(i18n.UsernameCharset))

	checkPassword(v, t, d.Password, d.PasswordConfirmation)

	return v.Errors()
}

type UserPasswordUpdateDTO struct {
	CurrentPassword      string `json:"current_password"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (d UserPasswordUpdateDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	v.Check(validation.NotBlank(d.CurrentPassword), "current_password", t(i18n.Required))
	checkPassword(v, t, d.Password, d.PasswordConfirmation)

	return v.Errors()
}

func checkPassword(v *validation.Validator, t i18n.Translator, password, confirmation string) {
	v.Check(validation.NotBlank(password), "password", t(i18n.Required))
	v.Check(validation.MinLength(password, 6), "password", t(i18n.MinLength, 6))
	v.Check(validation.MaxLength(password, 72), "password", t(i18n.MaxLength, 72))
	v.Check(validation.Matches(password, passwordRegex), "password", t(i18n.PasswordCharset, passwordSpecialChars))

	v.Check(validation.Equal(password, confirmation), "password_confirmation", t(i18n.PasswordConfirmation))
}

type UserUpdateDTO struct {
	KcalGoal         *int `json:"kcal_goal"`
	ProteinGoal      *int `json:"protein_goal"`
	CarbohydrateGoal *int `json:"carbohydrate_goal"`
	FatGoal          *int `json:"fat_goal"`

	Sex           *string  `json:"sex"`
	Age           *int     `json:"age"`
	Weight        *float64 `json:"weight"`
	Height        *int     `json:"height"`
	ActivityLevel *string  `json:"activity_level"`
	Objective     *string  `json:"objective"`
}

func (d UserUpdateDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	checkGoal(v, t, "kcal_goal", d.KcalGoal)
	checkGoal(v, t, "protein_goal", d.ProteinGoal)
	checkGoal(v, t, "carbohydrate_goal", d.CarbohydrateGoal)
	checkGoal(v, t, "fat_goal", d.FatGoal)

	if d.Sex != nil {
		v.Check(model.Sex(*d.Sex).Valid(), "sex", t(i18n.InvalidOption, optionNames(model.Sexes)))
	}
	if d.Age != nil {
		v.Check(validation.Between(*d.Age, minAge, maxAge), "age", t(i18n.Between, minAge, maxAge))
	}
	if d.Weight != nil {
		v.Check(validation.Between(*d.Weight, minWeight, maxWeight), "weight", t(i18n.Between, minWeight, maxWeight))
	}
	if d.Height != nil {
		v.Check(validation.Between(*d.Height, minHeight, maxHeight), "height", t(i18n.Between, minHeight, maxHeight))
	}
	if d.ActivityLevel != nil {
		v.Check(model.ActivityLevel(*d.ActivityLevel).Valid(), "activity_level", t(i18n.InvalidOption, optionNames(model.ActivityLevels)))
	}
	if d.Objective != nil {
		v.Check(model.Objective(*d.Objective).Valid(), "objective", t(i18n.InvalidOption, optionNames(model.Objectives)))
	}

	return v.Errors()
}

func optionNames[T ~string](options []T) string {
	names := make([]string, len(options))
	for i, option := range options {
		names[i] = string(option)
	}
	return strings.Join(names, ", ")
}

func checkGoal(v *validation.Validator, t i18n.Translator, field string, value *int) {
	if value == nil {
		return
	}
	v.Check(validation.Between(*value, 0, math.MaxUint16), field, t(i18n.Between, 0, math.MaxUint16))
}

type UserLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (d UserLoginDTO) Validate(t i18n.Translator) validation.Errors {
	v := validation.New()

	v.Check(validation.NotBlank(d.Username), "username", t(i18n.Required))
	v.Check(validation.NotBlank(d.Password), "password", t(i18n.Required))

	v.Check(validation.MinLength(d.Username, 5), "username", t(i18n.MinLength, 5))
	v.Check(validation.MinLength(d.Password, 6), "password", t(i18n.MinLength, 6))

	v.Check(validation.MaxLength(d.Username, 32), "username", t(i18n.MaxLength, 32))
	v.Check(validation.MaxLength(d.Password, 72), "password", t(i18n.MaxLength, 72))

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
