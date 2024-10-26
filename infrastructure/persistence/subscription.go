package persistence

import (
	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/domain/repository"
	"gorm.io/gorm"
)

type subscriptionPersistence struct {
	db *gorm.DB
}

func NewSubscriptionPersistence(db *gorm.DB) repository.SubscriptionRepository {
	return &subscriptionPersistence{db: db}
}

func (s subscriptionPersistence) Create(sub model.Subscription) error {
	return s.db.Create(&sub).Error
}

func (s subscriptionPersistence) Find(m []model.Subscription) ([]model.Subscription, error) {
	res := s.db.Find(&m)
	if res.Error != nil {
		return []model.Subscription{}, res.Error
	}
	return m, nil
}

func (s subscriptionPersistence) FindByModel(m model.Subscription) ([]model.Subscription, error) {
	var subs []model.Subscription
	res := s.db.Where(m).Find(&subs)
	if res.Error != nil {
		return []model.Subscription{}, res.Error
	}
	return subs, nil
}

func (s subscriptionPersistence) FindAll() ([]model.Subscription, error) {
	var subs []model.Subscription
	res := s.db.Find(&subs)
	if res.Error != nil {
		return []model.Subscription{}, res.Error
	}
	return subs, nil
}
