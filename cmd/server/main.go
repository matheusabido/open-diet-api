package main

import (
	"context"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/matheusabido/open-diet-api/internal/controller"
	"github.com/matheusabido/open-diet-api/internal/database"
	"github.com/matheusabido/open-diet-api/internal/groq"
	"github.com/matheusabido/open-diet-api/internal/middleware"
	"github.com/matheusabido/open-diet-api/internal/model"
	"github.com/matheusabido/open-diet-api/internal/openfoodfacts"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"github.com/matheusabido/open-diet-api/internal/service"
	"github.com/matheusabido/open-diet-api/internal/whisper"
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

	transcriber := whisper.New(os.Getenv("WHISPER_CLI_PATH"), os.Getenv("WHISPER_MODEL_PATH"))
	if !transcriber.Configured() {
		logrus.Warn("WHISPER_CLI_PATH and WHISPER_MODEL_PATH are not set, voice registration will not work")
	}

	groqClient := groq.New(os.Getenv("GROQ_KEY"))
	if !groqClient.Configured() {
		logrus.Warn("GROQ_KEY is not set, adding foods by photo will not work")
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
	foodRepository := repository.NewFoodRepository(db)
	foodItemRepository := repository.NewFoodItemRepository(db)
	foodAliasRepository := repository.NewFoodAliasRepository(db)

	tokenService := service.NewTokenService(refreshTokenRepository, jwtSecret)
	userService := service.NewUserService(userRepository, tokenService)
	foodService := service.NewFoodService(foodRepository)
	foodScanService := service.NewFoodScanService(groqClient, openfoodfacts.New())
	foodItemService := service.NewFoodItemService(foodItemRepository, foodRepository)
	foodAliasService := service.NewFoodAliasService(foodAliasRepository, foodRepository)
	voiceService := service.NewVoiceService(transcriber, foodAliasRepository, foodRepository)

	userController := controller.NewUserController(userService)
	foodController := controller.NewFoodController(foodService, foodScanService)
	foodItemController := controller.NewFoodItemController(foodItemService, voiceService)

	foodAliasController := controller.NewFoodAliasController(foodAliasService)

	auth := middleware.Auth(tokenService, userService)
	adminOnly := middleware.RequireRole(model.RoleAdmin)

	app.Post("/register", userController.Register)
	app.Post("/login", userController.Login)
	app.Post("/refresh", userController.Refresh)

	me := app.Group("/users/me", auth)
	me.Get("/", userController.Me)
	me.Patch("/", userController.UpdateMe)
	me.Patch("/password", userController.UpdatePassword)

	foods := app.Group("/foods", auth)
	foods.Get("/", foodController.List)
	foods.Post("/", foodController.Create)
	foods.Post("/photo", foodController.Photo)
	foods.Get("/barcode/:barcode", foodController.Barcode)
	foods.Get("/:id", foodController.Show)
	foods.Put("/:id", adminOnly, foodController.Update)
	foods.Delete("/:id", adminOnly, foodController.Delete)

	foodItems := app.Group("/food-items", auth)
	foodItems.Get("/", foodItemController.List)
	foodItems.Post("/", foodItemController.Create)
	foodItems.Post("/voice", foodItemController.Voice)
	foodItems.Get("/:id", foodItemController.Show)
	foodItems.Put("/:id", foodItemController.Update)
	foodItems.Delete("/:id", foodItemController.Delete)

	foodAliases := app.Group("/food-aliases", auth)
	foodAliases.Put("/:id", foodAliasController.Update)

	if err := app.Listen(":3000"); err != nil {
		logrus.Errorf("Could not start server: %v", err)
	}
}
