package persistence

import (
	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/domain/repository"
	"gorm.io/gorm"
)

type RssEntryPersistence struct {
	db *gorm.DB
}

func NewRssEntryPersistence(db *gorm.DB) repository.RssEnrtyRepository {
	return &RssEntryPersistence{db: db}
}

// func (r RssEntryPersistence) saveRSSEntries(rssURL string, entries []*gofeed.Item) {
func (r RssEntryPersistence) Create(entries []model.RssEntry) error {
	if len(entries) == 0 {
		return nil
	}
	res := r.db.Create(&entries)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r RssEntryPersistence) Find(entries []model.RssEntry) []model.RssEntry {
	if len(entries) == 0 {
		return []model.RssEntry{}
	}
	r.db.Find(&entries)
	return entries
}

func (r RssEntryPersistence) FindByModels(entries []model.RssEntry) []model.RssEntry {
	res := []model.RssEntry{}
	r.db.Where(&entries).Find(&res)
	return res
}
