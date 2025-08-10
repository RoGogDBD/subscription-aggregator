package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/RoGogDBD/subscription-aggregator/internal/models"
	"github.com/RoGogDBD/subscription-aggregator/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *service.SubscriptionService
}

func NewHandler(service *service.SubscriptionService) *Handler {
	return &Handler{service: service}
}

// CreateSubscription creates a subscription
// @Summary Create subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "subscription to create"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if sub.ServiceName == "" {
		http.Error(w, "service_name is required", http.StatusBadRequest)
		return
	}
	if sub.Price <= 0 {
		http.Error(w, "price must be positive", http.StatusBadRequest)
		return
	}
	if sub.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if sub.StartDate.IsZero() {
		http.Error(w, "start_date is required and must be in MM-YYYY format", http.StatusBadRequest)
		return
	}

	if sub.ID == "" {
		sub.ID = uuid.New().String()
	}

	if err := h.service.Create(sub); err != nil {
		http.Error(w, "failed to create subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(sub)
}

// ListSubscriptions lists subscriptions
// @Summary List subscriptions
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "filter by user id"
// @Param service_name query string false "filter by service name"
// @Success 200 {array} models.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *Handler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("user_id")
	serviceName := query.Get("service_name")

	subs, err := h.service.ListFiltered(userID, serviceName)
	if err != nil {
		http.Error(w, "failed to get subscriptions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(subs)
}

// CalculateTotal returns sum of subscriptions for period
// @Summary Calculate total price
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "filter by user id"
// @Param service_name query string false "filter by service name"
// @Param start_date query string false "start month (MM-YYYY)"
// @Param end_date query string false "end month (MM-YYYY)"
// @Success 200 {object} map[string]int "total"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/calculate [get]
func (h *Handler) CalculateTotal(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	userID := query.Get("user_id")
	serviceName := query.Get("service_name")
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	if startDate != "" {
		if ok, _ := regexp.MatchString(`^(0[1-9]|1[0-2])-\d{4}$`, startDate); !ok {
			http.Error(w, "start_date must be in MM-YYYY format", http.StatusBadRequest)
			return
		}
	}
	if endDate != "" {
		if ok, _ := regexp.MatchString(`^(0[1-9]|1[0-2])-\d{4}$`, endDate); !ok {
			http.Error(w, "end_date must be in MM-YYYY format", http.StatusBadRequest)
			return
		}
	}

	total, err := h.service.TotalFiltered(userID, serviceName, startDate, endDate)
	if err != nil {
		http.Error(w, "failed to calculate total: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int{"total": total})
}

// GetSubscription returns a subscription by ID
// @Summary Get subscription by ID
// @Description Get subscription details by subscription ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing subscription id", http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "invalid subscription id format", http.StatusBadRequest)
		return
	}

	sub, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "subscription not found: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sub)
}

// UpdateSubscription updates a subscription by ID
// @Summary Update subscription by ID
// @Description Update subscription details by subscription ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body models.Subscription true "Subscription data to update"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing subscription id", http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "invalid subscription id format", http.StatusBadRequest)
		return
	}

	var upd models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		http.Error(w, "invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if upd.ID != "" && upd.ID != id {
		http.Error(w, "id in body does not match id in URL", http.StatusBadRequest)
		return
	}
	upd.ID = id

	if upd.ServiceName == "" {
		http.Error(w, "service_name is required", http.StatusBadRequest)
		return
	}
	if upd.Price <= 0 {
		http.Error(w, "price must be positive", http.StatusBadRequest)
		return
	}
	if upd.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if upd.StartDate.IsZero() {
		http.Error(w, "start_date is required and must be in MM-YYYY format", http.StatusBadRequest)
		return
	}

	if err := h.service.Update(upd); err != nil {
		http.Error(w, "failed to update subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(upd)
}

// DeleteSubscription deletes a subscription by ID
// @Summary Delete subscription by ID
// @Description Delete subscription by subscription ID
// @Tags subscriptions
// @Param id path string true "Subscription ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing subscription id", http.StatusBadRequest)
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "invalid subscription id format", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, "failed to delete subscription: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HandleHealth returns service health status
// @Summary Health check
// @Description Returns health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
