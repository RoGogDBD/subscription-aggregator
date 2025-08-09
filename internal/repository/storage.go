package repository

import "fmt"

type Subscription struct {
	ID     string
	Amount float64
}

type MemStorage struct {
	data map[string]Subscription
}

func NewMemStorage() *MemStorage {
	return &MemStorage{data: make(map[string]Subscription)}
}

type Storage interface {
	Create(subscription Subscription) error
	GetByID(id string) (Subscription, error)
	Update(subscription Subscription) error
	Delete(id string) error
	List() ([]Subscription, error)
	Total() (float64, error)
}

func (m *MemStorage) Create(subscription Subscription) error {
	m.data[subscription.ID] = subscription
	return nil
}

func (m *MemStorage) GetByID(id string) (Subscription, error) {
	sub, ok := m.data[id]
	if !ok {
		return Subscription{}, fmt.Errorf("subscription with id %s not found", id)
	}
	return sub, nil
}

func (m *MemStorage) Update(subscription Subscription) error {
	_, ok := m.data[subscription.ID]
	if !ok {
		return fmt.Errorf("subscription with id %s not found", subscription.ID)
	}
	m.data[subscription.ID] = subscription
	return nil
}

func (m *MemStorage) Delete(id string) error {
	_, ok := m.data[id]
	if !ok {
		return fmt.Errorf("subscription with id %s not found", id)
	}
	delete(m.data, id)
	return nil
}

func (m *MemStorage) List() ([]Subscription, error) {
	var list []Subscription
	for _, sub := range m.data {
		list = append(list, sub)
	}
	return list, nil
}

func (m *MemStorage) Total() (float64, error) {
	var sum float64
	for _, sub := range m.data {
		sum += sub.Amount
	}
	return sum, nil
}
