package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/matheusabido/open-diet-api/internal/controller"
	"github.com/matheusabido/open-diet-api/internal/database"
	"github.com/matheusabido/open-diet-api/internal/middleware"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"github.com/matheusabido/open-diet-api/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warnf("Could not load .env file: %v", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logrus.Fatal("JWT_SECRET environment variable is required")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logrus.Fatal("DATABASE_URL environment variable is required")
	}

	ctx := context.Background()

	db, err := database.Connect(ctx, databaseURL)
	if err != nil {
		logrus.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	app := fiber.New()

	userRepository := repository.NewUserRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)

	tokenService := service.NewTokenService(refreshTokenRepository, jwtSecret)
	userService := service.NewUserService(userRepository, tokenService)

	userController := controller.NewUserController(userService)

	app.Post("/register", userController.Register)
	app.Post("/login", userController.Login)
	app.Post("/refresh", userController.Refresh)

	me := app.Group("/users/me", middleware.Auth(tokenService))
	me.Get("/", userController.Me)
	me.Patch("/", userController.UpdateMe)

	if err := app.Listen(":3000"); err != nil {
		logrus.Errorf("Could not start server: %v", err)
	}
}
