package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/banking-superapp/wealth-service/config"
	"github.com/banking-superapp/wealth-service/handler"
	"github.com/banking-superapp/wealth-service/repository"
	"github.com/banking-superapp/wealth-service/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	cfg := config.Load()

	mongoClient, err := repository.NewMongoClient(cfg.MongoAtlasURI)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database("banking_wealth")
	if err := repository.CreateIndexes(db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	mfRepo := repository.NewMFSchemeRepo(db)
	sipRepo := repository.NewSIPRepo(db)
	portRepo := repository.NewPortfolioRepo(db)
	riskRepo := repository.NewRiskProfileRepo(db)

	wealthSvc := service.NewWealthService(mfRepo, sipRepo, portRepo, riskRepo)
	wealthHandler := handler.NewWealthHandler(wealthSvc)

	app := fiber.New(fiber.Config{
		AppName:      cfg.ServiceName,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": cfg.ServiceName})
	})

	v1 := app.Group("/v1")
	wealth := v1.Group("/wealth")
	wealth.Get("/mf/catalogue", wealthHandler.GetCatalogue)
	wealth.Post("/mf/sip/create", wealthHandler.CreateSIP)
	wealth.Get("/portfolio", wealthHandler.GetPortfolio)
	wealth.Get("/portfolio/analytics", wealthHandler.GetPortfolioAnalytics)
	wealth.Post("/risk-profile", wealthHandler.AssessRiskProfile)
	wealth.Get("/risk-profile", wealthHandler.GetRiskProfile)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting %s on port %s", cfg.ServiceName, cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = app.ShutdownWithContext(ctx)
}
