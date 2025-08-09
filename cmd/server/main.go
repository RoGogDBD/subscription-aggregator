package main

import (
	"log"
	"net/http"

	"github.com/RoGogDBD/subscription-aggregator/internal/handlers"
	"github.com/RoGogDBD/subscription-aggregator/internal/logger"
	"github.com/RoGogDBD/subscription-aggregator/internal/repository"
	"github.com/RoGogDBD/subscription-aggregator/internal/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	parseFlags()

	if err := logger.Initialize(flagLogLevel); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	r, err := run()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	logger.Log.Info("Running server", zap.String("address", flagRunAddr))
	log.Fatal(http.ListenAndServe(flagRunAddr, r))
}

func run() (http.Handler, error) {
	storage := repository.NewMemStorage()
	subscriptionService := service.NewSubscriptionService(storage)
	h := handlers.NewHandler(subscriptionService)

	r := chi.NewRouter()

	r.Use(logger.RequestLogger)

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.CreateSubscription)
		r.Get("/", h.ViewSubscription)
		r.Get("/calculate", h.CalculateTotal)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetSubscription)
			r.Put("/", h.UpdateSubscription)
			r.Delete("/", h.DeleteSubscription)
		})
	})

	r.Get("/health", h.HandleHealth)

	return r, nil
}
