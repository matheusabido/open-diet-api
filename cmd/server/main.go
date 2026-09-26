package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/matheusabido/open-diet-api/internal/controller"
	"github.com/sirupsen/logrus"
)

func main() {
	app := fiber.New()

	userController := controller.UserController{}

	app.Post("/register", userController.Register)
	app.Post("/login", userController.Login)

	if err := app.Listen(":3000"); err != nil {
		logrus.Errorf("Could not start server: %v", err)
	}
}
