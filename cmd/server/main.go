// @title Subscription Aggregator API
// @version 1.0
// @description REST сервис для агрегации онлайн-подписок.

// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"net/http"

	_ "github.com/RoGogDBD/subscription-aggregator/api/docs"
	"github.com/RoGogDBD/subscription-aggregator/internal/config"
	"github.com/RoGogDBD/subscription-aggregator/internal/handlers"
	"github.com/RoGogDBD/subscription-aggregator/internal/logger"
	"github.com/RoGogDBD/subscription-aggregator/internal/repository"
	"github.com/RoGogDBD/subscription-aggregator/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or invalid: %v", err)
	}
	config.ParseFlags()

	if err := logger.Initialize(config.FlagLogLevel); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	pool := service.Connect()
	defer pool.Close()

	storage := repository.NewRepository(pool)
	subscriptionService := service.NewSubscriptionService(storage)
	h := handlers.NewHandler(subscriptionService)

	r, err := run(h)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	logger.Log.Info("Running server", zap.String("address", config.FlagRunAddr))
	log.Fatal(http.ListenAndServe(config.FlagRunAddr, r))
}

func run(h *handlers.Handler) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(logger.RequestLogger)

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.CreateSubscription)
		r.Get("/", h.ListSubscriptions)
		r.Get("/calculate", h.CalculateTotal)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetSubscription)
			r.Put("/", h.UpdateSubscription)
			r.Delete("/", h.DeleteSubscription)
		})
	})

	r.Get("/health", h.HandleHealth)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r, nil
}
