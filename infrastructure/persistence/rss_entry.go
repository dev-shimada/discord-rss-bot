package persistence

import (
	"fmt"
	"log/slog"

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
	links := make([]string, 0, len(entries))
	for _, e := range entries {
		links = append(links, e.EntryLink)
	}
	found := []model.RssEntry{}
	if err := r.db.Where("entry_link IN ?", links).Find(&found).Error; err != nil {
		slog.Error(fmt.Sprintf("failed to find RSS entries: %v", err))
		return []model.RssEntry{}
	}
	return found
}
