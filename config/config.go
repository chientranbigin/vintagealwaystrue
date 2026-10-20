package config

import (
	"os"

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
		dsn = "24jHtiebmCYdexH.root:SaleDbPassword@2026@tcp(gateway01.ap-southeast-1.prod.aws.tidbcloud.com:4000)/sale_db?tls=true&charset=utf8mb4&parseTime=True&loc=Local"
	}

	mediaBase := os.Getenv("MEDIA_BASE_URL")
	if mediaBase == "" {
		// Public CDN endpoint from existing server with SSL
		mediaBase = "https://www.tadimedia.vn/storage/"
	}

	return &Config{
		Port:         port,
		DatabaseDSN:  dsn,
		MediaBaseURL: mediaBase,
	}
}
