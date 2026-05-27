package service

import (
	"context"
	"fmt"
	"github.com/AndreSS-ntp/subscription-aggregator/sub-aggregator-service/internal/domain"
	"github.com/google/uuid"
	"time"
)

type Service struct {
	Repository Repository
}

type Repository interface {
	CreateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetSubscriptionById(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	ListSubscriptions(ctx context.Context, limit, offset int) ([]*domain.Subscription, error)
	TotalCost(ctx context.Context, userID uuid.UUID, serviceName *string, from, to time.Time) (int64, error)
}

func NewService(repository Repository) *Service {
	return &Service{repository}
}

func (s *Service) CreateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	return s.Repository.CreateSubscription(ctx, sub)
}

func (s *Service) UpdateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	return s.Repository.UpdateSubscription(ctx, sub)
}

func (s *Service) DeleteSubscription(ctx context.Context, id uuid.UUID) error {
	return s.Repository.DeleteSubscription(ctx, id)
}

func (s *Service) GetSubscriptionById(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	return s.Repository.GetSubscriptionById(ctx, id)
}

func (s *Service) ListSubscriptions(ctx context.Context, limit, offset int) ([]*domain.Subscription, error) {
	return s.Repository.ListSubscriptions(ctx, limit, offset)
}

func (s *Service) TotalCost(ctx context.Context, userID uuid.UUID, serviceName *string, from, to string) (*domain.SubscriptionCostDTO, error) {
	fromDate, err := time.Parse("01-2006", from)
	if err != nil {
		return nil, fmt.Errorf("parse from: %w", err)
	}

	toDate, err := time.Parse("01-2006", to)
	if err != nil {
		return nil, fmt.Errorf("parse to: %w", err)
	}

	if fromDate.After(toDate) {
		return nil, fmt.Errorf("from must be before or equal to to")
	}

	total, err := s.Repository.TotalCost(ctx, userID, serviceName, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	return &domain.SubscriptionCostDTO{
		UserID:      userID,
		ServiceName: serviceName,
		From:        from,
		To:          to,
		TotalCost:   total,
	}, nil
}
