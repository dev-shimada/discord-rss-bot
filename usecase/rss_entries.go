package usecase

import (
	"fmt"
	"log/slog"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/domain/repository"
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
	return model.RssEntry{
		RSSURL:      s.RSSURL,
		EntryTitle:  item.Title,
		EntryLink:   item.Link,
		PublishedAt: *item.PublishedParsed,
	}
}

func (f RssEntriesUsecase) CheckNewEntries(s []model.Subscription) []model.RssEntry {
	if len(s) == 0 {
		return []model.RssEntry{}
	}
	res := make([]model.RssEntry, 0, len(s))

	for _, sub := range s {
		items, err := f.rssFetcher.Fetch(sub.RSSURL)
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

	existingEntries := f.rr.FindByModels(cpRes)
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
