package usecase

import (
	"fmt"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/domain/repository"
	"golang.org/x/exp/slog"
)

type SubscriptionUsecase struct {
	sr repository.SubscriptionRepository
}

func NewSubscriptionUsecase(sr repository.SubscriptionRepository) SubscriptionUsecase {
	return SubscriptionUsecase{sr: sr}
}

func (s SubscriptionUsecase) Create(sub model.Subscription) string {
	err := s.sr.Create(sub)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to subscribe: %v", err))
		return "Failed to subscribe to RSS feed."
	}

	return "Successfully subscribed to RSS feed."
}

func (s SubscriptionUsecase) List(sub model.Subscription) ([]model.Subscription, error) {
	return s.sr.FindByModel(sub)
}

func (s SubscriptionUsecase) Delete(sub model.Subscription) error {
	return s.sr.Delete(sub)
}

func (s SubscriptionUsecase) FindAll() ([]model.Subscription, error) {
	return s.sr.FindAll()
}
