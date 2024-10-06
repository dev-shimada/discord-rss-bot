package usecase

import (
	"fmt"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/domain/repository"
	"golang.org/x/exp/slog"
)

type SubscriptionUsecase interface {
	FindAll() ([]model.Subscription, error)
	Create(sub model.Subscription) string
}

type subscriptionUsecase struct {
	sr repository.SubscriptionRepository
}

func NewSubscriptionUsecase(sr repository.SubscriptionRepository) SubscriptionUsecase {
	return &subscriptionUsecase{sr: sr}
}

func (s subscriptionUsecase) Create(sub model.Subscription) string {
	err := s.sr.Create(sub)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to subscribe: %v", err))
		return "Failed to subscribe to RSS feed."
	}

	return "Successfully subscribed to RSS feed."
}

func (s subscriptionUsecase) FindAll() ([]model.Subscription, error) {
	return s.sr.FindAll()
}
