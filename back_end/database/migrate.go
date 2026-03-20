package database

import (
	"log"

	"smart-fish/back_end/models"
)

func Migrate() {
	err := DB.AutoMigrate(
		// 原有模型
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
		// Flask SFR 兼容模型
		&models.Post{},
		&models.Comment{},
		&models.CommentOnComments{},
		&models.Image{},
		&models.LikeOnPosts{},
		&models.LikeOnComments{},
		&models.LikeOnCOCS{},
		&models.FishingRecord{},
		&models.FishCaught{},
		&models.IoTDevice{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migration completed")
}
