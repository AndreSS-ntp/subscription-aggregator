package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/AndreSS-ntp/subscription-aggregator/sub-aggregator-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service_name, price, user_id, start_date, end_date
	`

	row := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)

	createdSub := &domain.Subscription{}

	err := row.Scan(&createdSub.ID,
		&createdSub.ServiceName,
		&createdSub.Price,
		&createdSub.UserID,
		&createdSub.StartDate,
		&createdSub.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAlreadyExists
		}
		return nil, fmt.Errorf("scan created subscription: %w", err)
	}

	return createdSub, nil
}

func (r *Repository) UpdateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6
		RETURNING id, service_name, price, user_id, start_date, end_date
	`

	row := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)

	updatedSub := &domain.Subscription{}

	err := row.Scan(&updatedSub.ID,
		&updatedSub.ServiceName,
		&updatedSub.Price,
		&updatedSub.UserID,
		&updatedSub.StartDate,
		&updatedSub.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("scan updated subscription: %w", err)
	}

	return updatedSub, nil
}

func (r *Repository) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *Repository) GetSubscriptionById(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`

	sub := &domain.Subscription{}

	row := r.pool.QueryRow(ctx, query, id)
	err := row.Scan(&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("scan row: %w", err)
	}
	return sub, nil
}

func (r *Repository) ListSubscriptions(ctx context.Context, limit, offset int) ([]*domain.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		ORDER BY user_id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query items: %w", err)
	}
	defer rows.Close()

	var subs []*domain.Subscription

	for rows.Next() {
		sub := &domain.Subscription{}
		err = rows.Scan(&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("empty table")
			}
			return nil, fmt.Errorf("scan row: %w", err)
		}
		subs = append(subs, sub)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return subs, nil
}

func (r *Repository) TotalCost(ctx context.Context, userID uuid.UUID, serviceName *string, from, to time.Time) (int64, error) {
	query := `
        SELECT COALESCE(SUM(
            price * (
                (EXTRACT(YEAR FROM LEAST(COALESCE(end_date, $3), $3)) * 12 +
                 EXTRACT(MONTH FROM LEAST(COALESCE(end_date, $3), $3)))
                -
                (EXTRACT(YEAR FROM GREATEST(start_date, $2)) * 12 +
                 EXTRACT(MONTH FROM GREATEST(start_date, $2)))
                + 1
            )
        ), 0)
        FROM subscriptions
        WHERE user_id = $1
          AND start_date <= $3
          AND (end_date IS NULL OR end_date >= $2)
          AND ($4::text IS NULL OR service_name = $4)
    `

	var total int64
	err := r.pool.QueryRow(ctx, query, userID, from, to, serviceName).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("total cost query: %w", err)
	}

	return total, nil
}
