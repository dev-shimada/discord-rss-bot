package fetch

import (
	"github.com/mmcdole/gofeed"
)

type Rss struct {
	*gofeed.Parser
}

func NewRss() Rss {
	return Rss{gofeed.NewParser()}
}

func (r Rss) Fetch(rssURL string) ([]*gofeed.Item, error) {
	feed, err := r.ParseURL(rssURL)
	if err != nil {
		return nil, err
	}
	return feed.Items, nil
}
