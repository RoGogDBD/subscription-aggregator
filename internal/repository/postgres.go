package repository

import (
	"context"

	"github.com/RoGogDBD/subscription-aggregator/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(subscription models.Subscription) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO subscriptions 
			(id, service_name, price, user_id, start_date, end_date)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		subscription.ID,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	)
	return err
}

func (r *Repository) GetByID(id string) (models.Subscription, error) {
	var sub models.Subscription
	err := r.db.QueryRow(context.Background(),
		`SELECT id, service_name, price, user_id, start_date, end_date 
		 FROM subscriptions WHERE id=$1`, id).
		Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	return sub, err
}

func (r *Repository) Update(subscription models.Subscription) error {
	_, err := r.db.Exec(context.Background(),
		`UPDATE subscriptions SET 
			service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5
		 WHERE id=$6`,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
		subscription.ID,
	)
	return err
}

func (r *Repository) Delete(id string) error {
	_, err := r.db.Exec(context.Background(),
		"DELETE FROM subscriptions WHERE id=$1", id)
	return err
}

func (r *Repository) List() ([]models.Subscription, error) {
	rows, err := r.db.Query(context.Background(),
		"SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (r *Repository) Total() (int, error) {
	var total int
	err := r.db.QueryRow(context.Background(),
		"SELECT COALESCE(SUM(price), 0) FROM subscriptions").
		Scan(&total)
	return total, err
}

func (r *Repository) ListFiltered(date string) ([]models.Subscription, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE start_date <= $1 AND end_date >= $1`, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var sub models.Subscription
		if err := rows.Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate); err != nil {
			return nil, err
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

func (r *Repository) TotalFiltered(date string) (int, error) {
	var total int
	err := r.db.QueryRow(context.Background(),
		`SELECT COALESCE(SUM(price), 0) FROM subscriptions
		 WHERE start_date <= $1 AND end_date >= $1`, date).
		Scan(&total)
	return total, err
}
