package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseDSN  string
	MediaBaseURL string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		// Default to TiDB Cloud Serverless credentials
		dsn = "24jHtiebmCYdexH.root:SaleDbPassword@2026@tcp(gateway01.ap-southeast-1.prod.aws.tidbcloud.com:4000)/sale_db?tls=true&charset=utf8mb4&parseTime=True&loc=Asia%2FHo_Chi_Minh&time_zone=%27%2B07%3A00%27"
	} else {
		if !strings.Contains(dsn, "time_zone=") {
			dsn += "&time_zone=%27%2B07%3A00%27"
		}
		if strings.Contains(dsn, "loc=Local") {
			dsn = strings.Replace(dsn, "loc=Local", "loc=Asia%2FHo_Chi_Minh", 1)
		} else if !strings.Contains(dsn, "loc=") {
			dsn += "&loc=Asia%2FHo_Chi_Minh"
		}
	}

	mediaBase := os.Getenv("MEDIA_BASE_URL")
	if mediaBase == "" {
		// Cloudinary CDN endpoint for sale storage
		mediaBase = "https://res.cloudinary.com/leq5n17r/image/upload/"
	}

	return &Config{
		Port:         port,
		DatabaseDSN:  dsn,
		MediaBaseURL: mediaBase,
	}
}
