package handlers

import (
	"net/http"

	"github.com/RoGogDBD/subscription-aggregator/internal/service"
)

type Handler struct {
	service *service.SubscriptionService
}

func NewHandler(service *service.SubscriptionService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	// TODO: логика создания подписки
}

func (h *Handler) ViewSubscription(w http.ResponseWriter, r *http.Request) {
	// TODO: логика для ListenAndServe (возможно, это не очень подходящее название метода для обработчика)
}

func (h *Handler) CalculateTotal(w http.ResponseWriter, r *http.Request) {
	// TODO: логика подсчёта чего-то (например, общей суммы)
}

func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	// TODO: логика получения подписки по id
}

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	// TODO: логика обновления подписки
}

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	// TODO: логика удаления подписки
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	// TODO: логика health check
}
