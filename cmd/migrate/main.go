package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/matheusabido/open-diet-api/internal/database"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warnf("Could not load .env file: %v", err)
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

	applied, err := database.Migrate(ctx, db)
	for _, name := range applied {
		logrus.Infof("Applied migration %s", name)
	}
	if err != nil {
		logrus.Fatalf("Could not run migrations: %v", err)
	}

	if len(applied) == 0 {
		logrus.Info("No pending migrations")
	}
}
