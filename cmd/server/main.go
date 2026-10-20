package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"sale-app/config"
	"sale-app/handlers"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.LoadConfig()

	fmt.Println("Connecting to TiDB Cloud MySQL database...")
	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
	}
	fmt.Println("Connected to TiDB Cloud successfully!")

	r := gin.Default()

	// Native CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Health check for Render
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "engine": "golang"})
	})

	// Serve Static assets
	r.Static("/js", "./static/js")
	r.Static("/css", "./static/css")

	// SaleV2 API Handler
	v2 := handlers.NewSaleV2Handler(db, cfg)

	api := r.Group("/salev2/api")
	{
		api.GET("/products", v2.Products)
		api.GET("/product/:id", v2.ProductDetail)
		api.POST("/product/store", v2.StoreProduct)
		api.POST("/product/:id/update", v2.UpdateProduct)
		api.DELETE("/product/:id", v2.DestroyProduct)

		api.GET("/orders", v2.Orders)
		api.GET("/order/:id", v2.OrderDetail)
		api.POST("/api/order", v2.StoreOrder)
		api.POST("/order", v2.StoreOrder)
		api.POST("/order/:id", v2.UpdateOrder)
		api.POST("/order/:id/update", v2.UpdateOrder)

		api.GET("/consumers", v2.Consumers)
		api.GET("/consumer/:id", v2.ConsumerDetail)
		api.GET("/customer-history/:phone", v2.CustomerHistory)

		api.GET("/dashboard-stats", v2.DashboardStats)
		api.GET("/upload-sessions", v2.UploadSessions)
		api.GET("/upload-logs", v2.UploadLogs)
	}

	// Serve SaleV2 Vue SPA
	r.GET("/salev2", func(c *gin.Context) {
		c.File("./static/salev2/index.html")
	})

	// SPA fallback for all frontend routes
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/salev2/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
			return
		}
		c.File("./static/salev2/index.html")
	})

	// Redirect root to /salev2
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/salev2")
	})

	fmt.Printf("Sale Golang server running on port :%s...\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
