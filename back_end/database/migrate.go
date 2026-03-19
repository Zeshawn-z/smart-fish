package database

import (
	"log"

	"smart-fish/back_end/models"
)

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Region{},
		&models.FishingSpot{},
		&models.Device{},
		&models.Gateway{},
		&models.HistoricalData{},
		&models.EnvironmentData{},
		&models.WaterQualityData{},
		&models.Reminder{},
		&models.Notice{},
		&models.FishingSuggestion{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")
}
