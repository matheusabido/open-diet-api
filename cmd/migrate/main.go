package main

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/matheusabido/open-diet-api/internal/database"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logrus.Warnf("Could not load .env file: %v", err)
	}

	databaseConfig, err := database.ConfigFromEnv()
	if err != nil {
		logrus.Fatal(err)
	}

	ctx := context.Background()

	db, err := database.Connect(ctx, databaseConfig)
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
