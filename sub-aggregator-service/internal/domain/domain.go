package domain

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Subscription struct {
	ID          uuid.UUID
	ServiceName string
	Price       int64
	UserID      uuid.UUID
	StartDate   time.Time
	EndDate     *time.Time
}

type SubscriptionDTO struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int64     `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date"`
}

func ToSubscriptionDTO(sub *Subscription) *SubscriptionDTO {
	var endDate *string
	if sub.EndDate != nil {
		s := sub.EndDate.Format("01-2006")
		endDate = &s
	}
	startDate := sub.StartDate.Format("01-2006")

	return &SubscriptionDTO{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}
}

func ToSubscription(dto *SubscriptionDTO) (*Subscription, error) {
	dateFormat := "01-2006"
	// если месяц без ведущего ноля (1-9)
	if len(dto.StartDate) != 7 {
		dto.StartDate = "0" + dto.StartDate
	}
	startDate, err := time.Parse(dateFormat, dto.StartDate)
	if err != nil {
		return nil, fmt.Errorf("parse date from string: %w", err)
	}

	var endDate *time.Time
	if dto.EndDate != nil {
		if len(*dto.EndDate) != 7 {
			*dto.EndDate = "0" + *dto.EndDate
		}
		tmp, err := time.Parse(dateFormat, *dto.EndDate)
		if err != nil {
			return nil, fmt.Errorf("parse date from string: %w", err)
		}
		endDate = &tmp
	}

	return &Subscription{
		ID:          dto.ID,
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserID:      dto.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}

type SubscriptionCostDTO struct {
	UserID      uuid.UUID `json:"user_id"`
	ServiceName *string   `json:"service_name,omitempty"`
	From        string    `json:"from"`
	To          string    `json:"to"`
	TotalCost   int64     `json:"total_cost"`
}
