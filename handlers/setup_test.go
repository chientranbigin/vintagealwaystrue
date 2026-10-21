package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"sale-app/config"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	testDB     *gorm.DB
	testCfg    *config.Config
	testInitMu sync.Once
)

func getTestDB(t *testing.T) (*gorm.DB, *config.Config) {
	testInitMu.Do(func() {
		gin.SetMode(gin.TestMode)
		testCfg = config.LoadConfig()

		db, err := gorm.Open(mysql.Open(testCfg.DatabaseDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			t.Fatalf("Failed to connect to test database: %v", err)
		}
		testDB = db
	})
	if testDB == nil {
		t.Fatal("testDB is nil")
	}
	return testDB, testCfg
}

// setupSaleV2TestRouter initializes a test Gin router wired with the provided SaleV2Handler
func setupSaleV2TestRouter(v2 *SaleV2Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())

	routes := func(group *gin.RouterGroup) {
		group.GET("/products", v2.Products)
		group.GET("/product/:id", v2.ProductDetail)
		group.POST("/product/store", v2.StoreProduct)
		group.POST("/product/:id/update", v2.UpdateProduct)
		group.DELETE("/product/:id", v2.DestroyProduct)

		group.GET("/orders", v2.Orders)
		group.GET("/order/:id", v2.OrderDetail)
		group.POST("/order", v2.StoreOrder)
		group.POST("/api/order", v2.StoreOrder)
		group.POST("/order/:id", v2.UpdateOrder)
		group.POST("/order/:id/update", v2.UpdateOrder)
		group.POST("/order/:id/delete", v2.DestroyOrder)
		group.DELETE("/order/:id", v2.DestroyOrder)

		group.GET("/consumers", v2.Consumers)
		group.GET("/consumer/:id", v2.ConsumerDetail)
		group.GET("/customer-history/:phone", v2.CustomerHistory)

		group.GET("/dashboard-stats", v2.DashboardStats)
		group.GET("/upload-sessions", v2.UploadSessions)
		group.GET("/upload-sessions/:id/products", v2.UploadSessionProducts)
		group.GET("/products-by-upload", v2.ProductsByLatestUpload)
		group.GET("/sold-by-date", v2.SoldByDate)
		group.GET("/sold-by-date/:date/products", v2.SoldByDateProducts)
		group.GET("/upload-logs", v2.UploadLogs)
		group.POST("/upload-smart-group", v2.UploadSmartGroup)
		group.GET("/images/cleanup", v2.CleanupImages)
		group.POST("/images/cleanup", v2.CleanupImages)
	}

	routes(r.Group("/salev2/api"))
	routes(r.Group("/api"))

	return r
}

// performRequest executes a simulated HTTP request against a gin router
func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func uintPtr(u uint) *uint {
	return &u
}
