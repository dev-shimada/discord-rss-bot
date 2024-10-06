package model

import (
	"time"
)

type Subscription struct {
	ID        uint `gorm:"primaryKey"`
	ChannelID string
	RSSURL    string
	CreatedAt time.Time
}
