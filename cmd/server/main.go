package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mucunga90/ecommerce/config"
	"github.com/mucunga90/ecommerce/internal"
	"github.com/mucunga90/ecommerce/internal/api"
	"github.com/mucunga90/ecommerce/internal/database"
	"github.com/mucunga90/ecommerce/internal/events"
	"github.com/mucunga90/ecommerce/internal/manager"
	"github.com/mucunga90/ecommerce/internal/service"
	"github.com/mucunga90/ecommerce/internal/storage"
	"github.com/redis/go-redis/v9"

	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()

	db, err := database.New(cfg.DB.ConnectionStr) // connect to PostgreSQL
	if err != nil {
		log.Fatal(err)
	}
	defer closeDB(db) // close the database connection when the program exits

	if err := runMigrations(db); err != nil {
		log.Fatal(err)
	}

	// Ensure consumer group exists
	streamName := "order_notifications"
	groupName := "notifier_group"
	consumerName := "consumer-1"

	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr, Password: cfg.Redis.Password})

	publisher := events.NewPublisher(redisClient, streamName) // initialize Redis publisher

	storage := storage.New(db)

	manager := manager.New(cfg.Email.AdminEmail, storage, publisher)

	notifier := service.NewNotifier(cfg)

	if err := redisClient.XGroupCreateMkStream(ctx, streamName, groupName, "0").Err(); err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			log.Fatalf("❌ Failed to create Redis consumer group: %v", err)
		}
	}

	go events.NewNotifierConsumer(redisClient, streamName, groupName, consumerName, notifier).Start(ctx)

	// OIDC verifier
	verifier, err := service.NewVerifier(context.Background(), cfg.OIDC.Issuer, cfg.OIDC.JWKSURL, cfg.OIDC.Audience)
	if err != nil {
		log.Fatalf("OIDC verifier init failed: %v", err)
	}

	// Protected routes
	middlewares := []echo.MiddlewareFunc{
		verifier.EchoJWTMiddleware(),
		verifier.RequireValidClaims,
	}

	e := echo.New()
	handler := e.Group("/product", middlewares...)

	// Health check
	e.GET("/healthcheck", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	// Order routes
	handler.POST("/product/order", api.NewOrderHandler(manager).MakeOrder, service.RequireScopes("orders:write"))

	// Product routes
	handler.POST("/product/upload", api.NewUploadHandler(manager).ProductUpload, service.RequireScopes("product:write"))

	// Product pricing routes
	handler.GET("/product/prices", api.NewPriceHandler(manager).ProductPrices, service.RequireScopes("product:read"))

	port := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port))) // start the server on the specified port
}

func closeDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get DB instance:", err)
	}
	if err := sqlDB.Close(); err != nil {
		log.Fatal("Failed to close DB connection:", err)
	}
}

// Migrations runs the database migrations
func runMigrations(db *gorm.DB) error {
	// Run migrations here if needed
	// For example, you can use GORM's AutoMigrate or a migration library
	return db.AutoMigrate(&internal.Product{}, &internal.Order{}, &internal.Customer{})
}
