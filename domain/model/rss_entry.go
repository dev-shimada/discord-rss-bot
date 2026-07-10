package model

import (
	"time"
)

type RssEntry struct {
	ID               uint `gorm:"primaryKey"`
	RSSURL           string
	EntryTitle       string
	EntryLink        string
	EntryDescription string
	PublishedAt      time.Time
	CreatedAt        time.Time
}
