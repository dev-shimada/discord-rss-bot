package usecase

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/domain/repository"
	"github.com/mmcdole/gofeed"
)

type RssEntriesUsecase struct {
	rr         repository.RssEnrtyRepository
	rssFetcher repository.RssFetcher
}

func NewRssEntriesUsecase(rr repository.RssEnrtyRepository, rss repository.RssFetcher) RssEntriesUsecase {
	return RssEntriesUsecase{rr: rr, rssFetcher: rss}
}

func (f RssEntriesUsecase) Check(s model.Subscription) model.RssEntry {
	if s.RSSURL == "" {
		return model.RssEntry{}
	}
	items, err := f.rssFetcher.Fetch(s.RSSURL)
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to fetch RSS: %v", err))
		return model.RssEntry{}
	}
	if len(items) == 0 {
		return model.RssEntry{}
	}
	item := items[0]
	pub, _ := publishedAt(item)
	return model.RssEntry{
		RSSURL:           s.RSSURL,
		EntryTitle:       item.Title,
		EntryLink:        item.Link,
		EntryDescription: item.Description,
		PublishedAt:      pub,
	}
}

// publishedAt returns the item's publish time, falling back to the update
// time. ok is false when the feed provides no parseable timestamp at all.
func publishedAt(item *gofeed.Item) (time.Time, bool) {
	if item.PublishedParsed != nil {
		return *item.PublishedParsed, true
	}
	if item.UpdatedParsed != nil {
		return *item.UpdatedParsed, true
	}
	return time.Time{}, false
}

func (f RssEntriesUsecase) CheckNewEntries(s []model.Subscription) []model.RssEntry {
	if len(s) == 0 {
		return []model.RssEntry{}
	}
	res := make([]model.RssEntry, 0, len(s))

	// fetch each URL only once even when multiple channels subscribe to it
	type fetchResult struct {
		items []*gofeed.Item
		err   error
	}
	fetched := map[string]fetchResult{}

	for _, sub := range s {
		fr, ok := fetched[sub.RSSURL]
		if !ok {
			items, err := f.rssFetcher.Fetch(sub.RSSURL)
			fr = fetchResult{items: items, err: err}
			fetched[sub.RSSURL] = fr
		}
		if fr.err != nil {
			slog.Warn(fmt.Sprintf("failed to fetch RSS: %v", fr.err))
			continue
		}
		for _, item := range fr.items {
			pub, ok := publishedAt(item)
			if !ok {
				slog.Warn(fmt.Sprintf("skipping entry without a parseable date: %s", item.Link))
				continue
			}
			// skip if the item is older than the subscribed date
			if sub.CreatedAt.After(pub) {
				continue
			}
			res = append(res, model.RssEntry{
				RSSURL:           sub.RSSURL,
				EntryTitle:       item.Title,
				EntryLink:        item.Link,
				EntryDescription: item.Description,
				PublishedAt:      pub,
			})
		}
	}

	existingEntries := f.rr.Find(res)
	newEntries := diff(res, existingEntries)
	uniqueNewEntries := unique(newEntries)

	err := f.rr.Create(uniqueNewEntries)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to save RSS entries: %v", err))
		return nil
	}
	return uniqueNewEntries
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

func unique(s []model.RssEntry) []model.RssEntry {
	m := map[string]struct{}{}
	res := []model.RssEntry{}

	for _, v := range s {
		if _, ok := m[v.EntryLink]; !ok {
			m[v.EntryLink] = struct{}{}
			res = append(res, v)
		}
	}
	return res
}
