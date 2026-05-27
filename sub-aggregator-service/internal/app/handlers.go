package app

import (
	"context"
	"encoding/json"
	"errors"
	alogger "github.com/AndreSS-ntp/logger"
	"github.com/AndreSS-ntp/subscription-aggregator/sub-aggregator-service/internal/domain"
	"github.com/google/uuid"
	"net/http"
	"strconv"
)

type Command struct {
	Description string
	Handler     func(http.ResponseWriter, *http.Request)
}

type App struct {
	Commands map[string]Command
	Service  Service
}

type Service interface {
	CreateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	UpdateSubscription(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error)
	DeleteSubscription(ctx context.Context, id uuid.UUID) error
	GetSubscriptionById(ctx context.Context, id uuid.UUID) (*domain.Subscription, error)
	ListSubscriptions(ctx context.Context, limit, offset int) ([]*domain.Subscription, error)
}

type ErrorResponse struct {
	ErrMsg string `json:"error"`
}

func NewApp(s Service) *App {
	a := App{}
	var commands = map[string]Command{
		"POST /v1/subscription":        Command{"Создать новую запись о подписке.", a.CreateSubscription},
		"DELETE /v1/subscription/{id}": Command{"Удалить запись о подписке по ID.", a.DeleteSubscription},
		"GET /v1/subscription/{id}":    Command{"Получить запись о подписке по ID.", a.GetSubscription},
		"PUT /v1/subscription/{id}":    Command{"Обновить запись о подписке по ID", a.UpdateSubscription},
		"GET /v1/subscriptions":        Command{"Получить список всех записей подписок (параметры пагинации: limit, offset)", a.ListSubscriptions},
	}
	a.Commands = commands
	a.Service = s
	return &a
}

func (a *App) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var subDTO *domain.SubscriptionDTO
	err := json.NewDecoder(r.Body).Decode(&subDTO)
	if err != nil {
		sendError(ctx, w, "invalid json", 400)
		return
	}

	sub, err := domain.ToSubscription(subDTO)
	if err != nil {
		sendError(ctx, w, "internal server error", 500)
		return
	}

	createdSub, err := a.Service.CreateSubscription(ctx, sub)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			sendError(ctx, w, domain.ErrAlreadyExists.Error(), 409)
			return
		}
		sendError(ctx, w, "internal server error", 500)
		return
	}

	data, err := json.Marshal(domain.ToSubscriptionDTO(createdSub))
	if err != nil {
		sendError(ctx, w, "internal server error", 500)
		return
	}

	sendOk(ctx, w, data, 201)
}

func (a *App) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		sendError(ctx, w, "invalid id", 400)
		return
	}

	err = a.Service.DeleteSubscription(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(ctx, w, domain.ErrNotFound.Error(), 404)
			return
		}
		sendError(ctx, w, "internal server error", 500)
		return
	}
	w.WriteHeader(204)
}

func (a *App) GetSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		sendError(ctx, w, "invalid id", 400)
		return
	}

	sub, err := a.Service.GetSubscriptionById(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(ctx, w, domain.ErrNotFound.Error(), 404)
			return
		}
		sendError(ctx, w, "internal server error", 500)
		return
	}

	data, err := json.Marshal(domain.ToSubscriptionDTO(sub))
	if err != nil {
		sendError(ctx, w, "internal server error", 500)
		return
	}

	sendOk(ctx, w, data, 200)
}

func (a *App) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		sendError(ctx, w, "invalid id", 400)
		return
	}

	var subDTO *domain.SubscriptionDTO
	err = json.NewDecoder(r.Body).Decode(&subDTO)
	if err != nil {
		sendError(ctx, w, "invalid json", 400)
		return
	}

	sub, err := domain.ToSubscription(subDTO)
	if err != nil {
		sendError(ctx, w, "internal server error", 500)
		return
	}
	sub.ID = id

	updatedSub, err := a.Service.UpdateSubscription(ctx, sub)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			sendError(ctx, w, domain.ErrNotFound.Error(), 404)
			return
		}
		sendError(ctx, w, "internal server error", 500)
		return
	}

	data, err := json.Marshal(domain.ToSubscriptionDTO(updatedSub))
	if err != nil {
		sendError(ctx, w, "internal server error", 500)
		return
	}

	sendOk(ctx, w, data, 200)
}

func (a *App) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 20
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	subs, err := a.Service.ListSubscriptions(ctx, limit, offset)
	if err != nil {
		sendError(ctx, w, "internal server error", 500)
	}

	subsDTO := make([]*domain.SubscriptionDTO, 0, len(subs))
	for _, v := range subs {
		subsDTO = append(subsDTO, domain.ToSubscriptionDTO(v))
	}

	data, err := json.Marshal(subsDTO)
	if err != nil {
		sendError(ctx, w, "internal server error", 500)
		return
	}

	sendOk(ctx, w, data, 200)
}

func sendError(ctx context.Context, w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w_err := json.NewEncoder(w).Encode(ErrorResponse{msg})
	if w_err != nil {
		alogger.FromContext(ctx).Error(ctx, "cant write a response: "+w_err.Error())
	}
}

func sendOk(ctx context.Context, w http.ResponseWriter, data []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, w_err := w.Write(data)
	if w_err != nil {
		alogger.FromContext(ctx).Error(ctx, "cant write a response: "+w_err.Error())
	}
}
