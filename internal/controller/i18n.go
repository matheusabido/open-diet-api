package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/matheusabido/open-diet-api/internal/i18n"
)

func translator(c fiber.Ctx) i18n.Translator {
	return i18n.New(c.AcceptsLanguages(i18n.Supported...))
}
