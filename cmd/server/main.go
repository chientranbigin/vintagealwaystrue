package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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
	db.Exec("SET time_zone = '+07:00'")
	fmt.Println("Connected to TiDB Cloud successfully (time_zone = +07:00)!")

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

	// Auth Handler & Middleware
	auth := handlers.NewAuthHandler(db, cfg)
	r.Use(auth.AuthMiddleware())

	// Health check for Render
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "engine": "golang"})
	})

	// Login & Logout routes
	r.GET("/login", func(c *gin.Context) {
		c.File("./static/login.html")
	})
	r.POST("/api/login", auth.Login)
	r.GET("/logout", auth.Logout)

	// Serve Static assets
	r.Static("/js", "./public/js")
	r.Static("/css", "./public/css")
	r.Static("/admin-theme", "./public/admin-theme")
	r.Static("/customer-theme", "./public/customer-theme")

	// Storage proxy / file server
	r.GET("/storage/*filepath", func(c *gin.Context) {
		subPath := c.Param("filepath")
		localPath := "./storage" + subPath
		if _, err := os.Stat(localPath); err == nil {
			c.File(localPath)
			return
		}
		// If not found locally, redirect to public CDN
		mediaBase := cfg.MediaBaseURL
		if !strings.HasSuffix(mediaBase, "/") {
			mediaBase += "/"
		}
		c.Redirect(http.StatusMovedPermanently, mediaBase+strings.TrimPrefix(subPath, "/"))
	})

	// ==========================================
	// Sale V1 Handlers (Server-rendered Blade UI)
	// ==========================================
	v1 := handlers.NewSaleV1Handler(db, cfg)
	sale := r.Group("/sale")
	{
		sale.GET("", v1.Home)
		sale.GET("/", v1.Home)

		sale.GET("/product", v1.Products)
		sale.POST("/product", func(c *gin.Context) {
			values := url.Values{}
			for _, key := range []string{"search", "type", "status", "from_search_key", "from_search_value", "to_search_value", "sold_from", "sold_to"} {
				if val := strings.TrimSpace(c.PostForm(key)); val != "" {
					values.Set(key, val)
				}
			}
			redirectURL := "/sale/product"
			if qs := values.Encode(); qs != "" {
				redirectURL += "?" + qs
			}
			c.Redirect(http.StatusFound, redirectURL)
		})
		sale.GET("/available-product", v1.AvailableProducts)
		sale.POST("/available-product", func(c *gin.Context) {
			values := url.Values{}
			for _, key := range []string{"search", "type", "from_search_key", "from_search_value", "to_search_value", "from_search_key_2", "from_search_value_2", "to_search_value_2"} {
				if val := strings.TrimSpace(c.PostForm(key)); val != "" {
					values.Set(key, val)
				}
			}
			redirectURL := "/sale/available-product"
			if qs := values.Encode(); qs != "" {
				redirectURL += "?" + qs
			}
			c.Redirect(http.StatusFound, redirectURL)
		})
		sale.GET("/product/:id", v1.ProductDetail)
		sale.POST("/product/:id", v1.ProductDetail)
		sale.GET("/product/delete/:id", v1.ProductDelete)
		sale.GET("/product-create", v1.ProductCreate)
		sale.POST("/product-create", v1.ProductCreate)
		sale.GET("/product-upload", v1.ProductUpload)
		sale.POST("/product-upload", v1.ProductUpload)
		sale.GET("/product-batch-upload", v1.ProductBatchUpload)
		sale.POST("/product-batch-upload/single", v1.ProductBatchUploadSingle)

		sale.GET("/order", v1.Orders)
		sale.POST("/order", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/sale/order")
		})
		sale.GET("/order-summary", v1.SummaryOrders)
		sale.GET("/order/:id", v1.OrderDetail)
		sale.POST("/order/:id", v1.OrderDetail)
		sale.GET("/order/delete/:id", v1.OrderDelete)
		sale.GET("/order-create", v1.OrderCreate)
		sale.POST("/order-create", v1.OrderCreate)
		sale.GET("/smart-order-create", v1.SmartOrderCreate)
		sale.POST("/smart-order-create", v1.SmartOrderCreate)
		sale.GET("/complete/:id", v1.OrderComplete)
		sale.GET("/ship/:id", v1.OrderShip)
		sale.GET("/hold/:id", v1.OrderHold)
		sale.GET("/consumer", v1.Consumers)
		sale.GET("/upload-sessions", v1.UploadSessions)
		sale.GET("/upload-sessions/:sessionId/products", v1.UploadSessionProducts)
	}

	// ==========================================
	// Sale V2 API Handler (REST API for Vue SPA)
	// ==========================================
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
		api.POST("/order/:id/delete", v2.DestroyOrder)
		api.DELETE("/order/:id", v2.DestroyOrder)

		api.GET("/consumers", v2.Consumers)
		api.GET("/consumer/:id", v2.ConsumerDetail)
		api.GET("/customer-history/:phone", v2.CustomerHistory)

		api.GET("/dashboard-stats", v2.DashboardStats)
		api.GET("/orders-stats", v2.OrderStats)
		api.GET("/orders/stats", v2.OrderStats)
		api.GET("/upload-sessions", v2.UploadSessions)
		api.GET("/upload-sessions/:id/products", v2.UploadSessionProducts)
		api.GET("/products-by-upload", v2.ProductsByLatestUpload)
		api.GET("/sold-by-date", v2.SoldByDate)
		api.GET("/sold-by-date/:date/products", v2.SoldByDateProducts)
		api.GET("/upload-logs", v2.UploadLogs)
		api.POST("/upload-smart-group", v2.UploadSmartGroup)
		api.GET("/images/cleanup", v2.CleanupImages)
		api.POST("/images/cleanup", v2.CleanupImages)
	}

	// Direct endpoint for external cron jobs (e.g. cron-job.org)
	r.GET("/api/images/cleanup", v2.CleanupImages)
	r.POST("/api/images/cleanup", v2.CleanupImages)

	// Serve SaleV2 Vue SPA
	r.GET("/salev2", func(c *gin.Context) {
		c.File("./static/salev2/index.html")
	})

	// SPA fallback for all frontend routes
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/salev2/api") || strings.HasPrefix(path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API route not found"})
			return
		}
		if strings.HasPrefix(path, "/salev2") {
			c.File("./static/salev2/index.html")
			return
		}
		if path == "/sale" || strings.HasPrefix(path, "/sale/") {
			c.Redirect(http.StatusFound, "/sale")
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
