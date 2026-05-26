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
	EndDate     time.Time
}

type SubscriptionDTO struct {
	ID          uuid.UUID `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int64     `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
}

func ToSubscriptionDTO(sub *Subscription) *SubscriptionDTO {
	startYear, startMonth, _ := sub.StartDate.Date()
	endYear, endMonth, _ := sub.EndDate.Date()

	startDate := fmt.Sprintf("%d-%d", startMonth, startYear)
	endDate := fmt.Sprintf("%d-%d", endMonth, endYear)
	// если месяц без ведущего ноля
	if startMonth < 10 {
		startDate = "0" + startDate
	}
	if endMonth < 10 {
		endDate = "0" + endDate
	}

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
	if len(dto.EndDate) != 7 {
		dto.EndDate = "0" + dto.EndDate
	}
	endDate, err := time.Parse(dateFormat, dto.EndDate)
	if err != nil {
		return nil, fmt.Errorf("parse date from string: %w", err)
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
