package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Subscription модель подписки
// swagger:model Subscription
type Subscription struct {
	// ID подписки
	// example: 10275647-eab5-4828-8dc1-e5335ae52abd
	ID string `json:"id,omitempty"`
	// Название сервиса
	// example: Yandex Plus
	ServiceName string `json:"service_name"`
	// Стоимость в рублях
	// example: 400
	Price int `json:"price"`
	// Пользователь (UUID)
	// example: 60601fee-2bf1-4721-ae6f-7636e79a0cba
	UserID string `json:"user_id"`
	// Дата начала (MM-YYYY)
	// example: 07-2025
	StartDate time.Time `json:"start_date"`
	// Дата окончания (MM-YYYY)
	// example: 12-2025
	EndDate *time.Time `json:"end_date,omitempty"`
}

func parseMonthYear(input string) (time.Time, error) {
	parts := strings.Split(input, "-")
	if len(parts) != 2 {
		return time.Time{}, errors.New("ожидается формат ММ-ГГГГ")
	}
	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return time.Time{}, errors.New("некорректный месяц")
	}
	yearStr := parts[1]

	if len(yearStr) == 2 {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return time.Time{}, errors.New("некорректный год")
		}
		year += 2000
		return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
	} else if len(yearStr) == 4 {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return time.Time{}, errors.New("некорректный год")
		}
		if year < 1900 || year > 9999 {
			return time.Time{}, errors.New("год вне диапазона")
		}
		return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
	} else {
		return time.Time{}, errors.New("год должен быть в формате ГГГГ или ГГ")
	}
}

func formatToMMYYYY(t time.Time) string {
	year := t.Year()
	return fmt.Sprintf("%02d-%04d", t.Month(), year)
}

func (s Subscription) MarshalJSON() ([]byte, error) {
	type Alias Subscription
	aux := struct {
		ID          string  `json:"id"`
		ServiceName string  `json:"service_name"`
		Price       int     `json:"price"`
		UserID      string  `json:"user_id"`
		StartDate   string  `json:"start_date"`
		EndDate     *string `json:"end_date,omitempty"`
	}{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   formatToMMYYYY(s.StartDate),
	}

	if s.EndDate != nil {
		endDateStr := formatToMMYYYY(*s.EndDate)
		aux.EndDate = &endDateStr
	}

	return json.Marshal(aux)
}

func (s *Subscription) UnmarshalJSON(data []byte) error {
	type Alias Subscription
	aux := &struct {
		StartDate string  `json:"start_date"`
		EndDate   *string `json:"end_date,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	startDate, err := parseMonthYear(aux.StartDate)
	if err != nil {
		return fmt.Errorf("ошибка парсинга start_date: %w", err)
	}
	s.StartDate = startDate

	if aux.EndDate != nil {
		endDate, err := parseMonthYear(*aux.EndDate)
		if err != nil {
			return fmt.Errorf("ошибка парсинга end_date: %w", err)
		}
		s.EndDate = &endDate
	} else {
		s.EndDate = nil
	}

	return nil
}
