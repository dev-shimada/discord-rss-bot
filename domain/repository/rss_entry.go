package repository

import "github.com/dev-shimada/discord-rss-bot/domain/model"

type RssEnrtyRepository interface {
	Create(entries []model.RssEntry) error
	Find(entries []model.RssEntry) []model.RssEntry
	FindByModels(entries []model.RssEntry) []model.RssEntry
}
