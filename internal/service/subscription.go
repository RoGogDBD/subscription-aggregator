package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RoGogDBD/subscription-aggregator/internal/models"
	"github.com/RoGogDBD/subscription-aggregator/internal/repository"
)

type SubscriptionService struct {
	storage *repository.Repository
}

func NewSubscriptionService(storage *repository.Repository) *SubscriptionService {
	return &SubscriptionService{storage: storage}
}

func (s *SubscriptionService) Create(sub models.Subscription) error {
	_, err := s.storage.GetByID(sub.ID)
	if err == nil {
		return fmt.Errorf("subscription with id %s already exists", sub.ID)
	}
	return s.storage.Create(sub)
}

func (s *SubscriptionService) GetByID(id string) (models.Subscription, error) {
	return s.storage.GetByID(id)
}

func (s *SubscriptionService) Update(sub models.Subscription) error {
	return s.storage.Update(sub)
}

func (s *SubscriptionService) Delete(id string) error {
	return s.storage.Delete(id)
}

func (s *SubscriptionService) List() ([]models.Subscription, error) {
	return s.storage.List()
}

func (s *SubscriptionService) ListFiltered(userID, serviceName string) ([]models.Subscription, error) {
	all, err := s.storage.List()
	if err != nil {
		return nil, err
	}

	var out []models.Subscription
	for _, sub := range all {
		if userID != "" && sub.UserID != userID {
			continue
		}
		if serviceName != "" {
			if !strings.Contains(strings.ToLower(sub.ServiceName), strings.ToLower(serviceName)) {
				continue
			}
		}
		out = append(out, sub)
	}
	return out, nil
}

func (s *SubscriptionService) TotalFiltered(userID, serviceName, startDate, endDate string) (int, error) {
	subs, err := s.ListFiltered(userID, serviceName)
	if err != nil {
		return 0, err
	}
	if startDate == "" && endDate == "" {
		total := 0
		for _, sub := range subs {
			total += sub.Price
		}
		return total, nil
	}

	var periodStart *time.Time
	var periodEnd *time.Time

	if startDate != "" {
		st, _, err := parseMonthYearRange(startDate)
		if err != nil {
			return 0, fmt.Errorf("invalid start_date: %w", err)
		}
		periodStart = &st
	}
	if endDate != "" {
		_, ed, err := parseMonthYearRange(endDate)
		if err != nil {
			return 0, fmt.Errorf("invalid end_date: %w", err)
		}
		periodEnd = &ed
	}

	total := 0
	for _, sub := range subs {
		if subscriptionOverlapsPeriod(sub, periodStart, periodEnd) {
			total += sub.Price
		}
	}
	return total, nil
}

func parseMonthYearRange(input string) (time.Time, time.Time, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("expected MM-YYYY")
	}
	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month")
	}
	yearStr := parts[1]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid year")
	}
	if len(yearStr) == 2 {
		year += 2000
	}
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	next := start.AddDate(0, 1, 0)
	end := next.Add(-time.Nanosecond)
	return start, end, nil
}

func subscriptionOverlapsPeriod(sub models.Subscription, periodStart, periodEnd *time.Time) bool {
	subStart := sub.StartDate
	var subEnd *time.Time
	if sub.EndDate != nil {
		t := *sub.EndDate
		subEnd = &t
	}

	if periodStart == nil && periodEnd == nil {
		return true
	}


	if periodStart != nil && subEnd != nil && subEnd.Before(*periodStart) {
		return false
	}
	
	if periodEnd != nil && subStart.After(*periodEnd) {
		return false
	}

	return true
}
