package usecase

import (
	"fmt"
	"log/slog"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/domain/repository"
	"github.com/mmcdole/gofeed"
)

type RssEntriesUsecase struct {
	rr repository.RssEnrtyRepository
}

func NewRssEntriesUsecase(rr repository.RssEnrtyRepository) RssEntriesUsecase {
	return RssEntriesUsecase{rr: rr}
}

func (f RssEntriesUsecase) CheckNewEntries(s []model.Subscription) []model.RssEntry {
	if len(s) == 0 {
		return []model.RssEntry{}
	}
	res := make([]model.RssEntry, 0, len(s))

	for _, sub := range s {
		items, err := fetchRSS(sub.RSSURL)
		if err != nil {
			slog.Warn(fmt.Sprintf("failed to fetch RSS: %v", err))
			continue
		}
		for _, item := range items {
			// skip if the item is older than the subscribed date
			if sub.CreatedAt.After(*item.PublishedParsed) {
				continue
			}
			res = append(res, model.RssEntry{
				RSSURL:      sub.RSSURL,
				EntryTitle:  item.Title,
				EntryLink:   item.Link,
				PublishedAt: *item.PublishedParsed,
			})
		}
	}
	cpRes := make([]model.RssEntry, len(s))
	copy(cpRes, res)

	existingEntries := f.rr.Find(cpRes)
	newEntries := diff(res, existingEntries)
	err := f.rr.Create(newEntries)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to save RSS entries: %v", err))
		return nil
	}
	return newEntries
}

func fetchRSS(rssURL string) ([]*gofeed.Item, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		return nil, err
	}
	return feed.Items, nil
}

func diff(s1, s2 []model.RssEntry) []model.RssEntry {
	diffSlice := []model.RssEntry{}
	cmpMap := map[string]int{}

	for _, v := range s2 {
		cmpMap[v.EntryLink] += 1
	}

	for _, v := range s1 {
		t, ok := cmpMap[v.EntryLink]
		if !ok {
			diffSlice = append(diffSlice, v)
			continue
		}
		if t == 1 {
			delete(cmpMap, v.EntryLink)
		} else {
			cmpMap[v.EntryLink] -= 1
		}
	}
	return diffSlice
}
