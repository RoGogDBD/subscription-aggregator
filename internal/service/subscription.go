package service

import (
	"fmt"

	"github.com/RoGogDBD/subscription-aggregator/internal/repository"
)

type SubscriptionService struct {
	storage repository.Storage
}

func NewSubscriptionService(storage repository.Storage) *SubscriptionService {
	return &SubscriptionService{storage: storage}
}

func (s *SubscriptionService) Create(sub repository.Subscription) error {
	_, err := s.storage.GetByID(sub.ID)
	if err == nil {
		return fmt.Errorf("subscription with id %s already exists", sub.ID)
	}
	return s.storage.Create(sub)
}

func (s *SubscriptionService) GetByID(id string) (repository.Subscription, error) {
	return s.storage.GetByID(id)
}

func (s *SubscriptionService) Update(sub repository.Subscription) error {
	return s.storage.Update(sub)
}

func (s *SubscriptionService) Delete(id string) error {
	return s.storage.Delete(id)
}

func (s *SubscriptionService) List() ([]repository.Subscription, error) {
	return s.storage.List()
}

func (s *SubscriptionService) Total() (float64, error) {
	return s.storage.Total()
}
