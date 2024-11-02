package repository

import (
	"github.com/mmcdole/gofeed"
)

type RssFetcher interface {
	Fetch(rssURL string) ([]*gofeed.Item, error)
}
