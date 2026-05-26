package service

import (
	"context"
	"github.com/AndreSS-ntp/subscription-aggregator/sub-aggregator-service/internal/domain"
	"github.com/google/uuid"
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
