package validation

import (
	"cmp"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

type Errors map[string][]string

type Validator struct {
	errors Errors
}

func New() *Validator {
	return &Validator{errors: Errors{}}
}

func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.errors[field] = append(v.errors[field], message)
	}
}

func (v *Validator) Valid() bool {
	return len(v.errors) == 0
}

func (v *Validator) Errors() Errors {
	return v.errors
}

func NotBlank(s string) bool {
	return strings.TrimSpace(s) != ""
}

func MinLength(s string, n int) bool {
	return utf8.RuneCountInString(s) >= n
}

func MaxLength(s string, n int) bool {
	return utf8.RuneCountInString(s) <= n
}

func MaxBytes(s string, n int) bool {
	return len(s) <= n
}

func Matches(s string, re *regexp.Regexp) bool {
	return re.MatchString(s)
}

func Equal[T comparable](a, b T) bool {
	return a == b
}

func Between[T cmp.Ordered](value, min, max T) bool {
	return value >= min && value <= max
}

func Min[T cmp.Ordered](value, min T) bool {
	return value >= min
}

func Max[T cmp.Ordered](value, max T) bool {
	return value <= max
}

func GreaterThan[T cmp.Ordered](value, min T) bool {
	return value > min
}

func In[T comparable](value T, allowed ...T) bool {
	return slices.Contains(allowed, value)
}

func ISODateTime(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}

func ISODate(s string) bool {
	_, err := time.Parse(time.DateOnly, s)
	return err == nil
}
