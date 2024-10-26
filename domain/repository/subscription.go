package repository

import "github.com/dev-shimada/discord-rss-bot/domain/model"

type SubscriptionRepository interface {
	Create(sub model.Subscription) error
	Find(m []model.Subscription) ([]model.Subscription, error)
	FindByModel(m model.Subscription) ([]model.Subscription, error)
	FindAll() ([]model.Subscription, error)
	Delete(m model.Subscription) error
}
