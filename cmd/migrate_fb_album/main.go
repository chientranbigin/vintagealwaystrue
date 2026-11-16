package main

import (
	"fmt"
	"log"

	"sale-app/config"
	"sale-app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println("Connecting to TiDB Cloud MySQL database...")
	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 1. AutoMigrate Category table
	fmt.Println("Migrating Category table...")
	if err := db.AutoMigrate(&models.Category{}); err != nil {
		log.Fatalf("Failed to migrate categories table: %v", err)
	}

	// 2. Add fb_album_photo_id column to products if not exists
	fmt.Println("Checking fb_album_photo_id in products table...")
	db.Exec("ALTER TABLE products ADD COLUMN IF NOT EXISTS fb_album_photo_id VARCHAR(100) NULL")

	// 3. Seed default categories if not exists
	defaultCategories := []struct {
		Code string
		Name string
	}{
		{"TROUSERS", "Quần Âu (Trousers)"},
		{"JACKET", "Áo Khoác (Jacket)"},
		{"HARRINGTON_JACKET", "Áo Khoác Harrington"},
		{"SHIRT", "Áo Sơ Mi (Shirt)"},
		{"BLAZER", "Áo Blazer"},
		{"SUIT", "Bộ Suit"},
		{"GILE", "Áo Gile"},
		{"POLO SHIRT", "Áo Polo"},
		{"TIE", "Cà Vạt (Tie)"},
		{"BELT", "Thắt Lưng (Belt)"},
		{"HAT", "Mũ (Hat)"},
	}

	for _, item := range defaultCategories {
		var cat models.Category
		res := db.Where("code = ?", item.Code).First(&cat)
		if res.Error != nil {
			newCat := models.Category{
				Code:     item.Code,
				Name:     item.Name,
				IsActive: true,
			}
			db.Create(&newCat)
			fmt.Printf("Created category: %s (%s)\n", item.Code, item.Name)
		} else {
			fmt.Printf("Category already exists: %s (%s), AlbumID: %v\n", cat.Code, cat.Name, cat.FBAlbumID)
		}
	}

	fmt.Println("Migration and seeding completed successfully!")
}
