package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/faisal/crypto/backend/internal/config"
	"github.com/faisal/crypto/backend/internal/handlers"
	"github.com/faisal/crypto/backend/internal/services/market"
	"github.com/faisal/crypto/backend/internal/services/portfolio"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	router := gin.Default()
	router.Use(config.CORSMiddleware(cfg.AllowedOrigins))

	api := router.Group("/api")

	marketService := market.NewService(cfg)
	marketHandler := handlers.NewMarketHandler(marketService)
	marketHandler.Register(api)

	// Using in-memory storage for development
	// TODO: Switch to MongoDB when connection is ready
	portfolioService := portfolio.NewService(cfg)
	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)
	portfolioHandler.Register(api)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Go backend listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
