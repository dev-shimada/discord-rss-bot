package usecase

import (
	"fmt"
	"log/slog"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/domain/repository"
	"github.com/mmcdole/gofeed"
)

type RssEntriesUsecase interface {
	CheckNewEntries(s []model.Subscription) []model.RssEntry
}

type rssEntriesUsecase struct {
	rr repository.RssEnrtyRepository
}

func NewRssEntriesUsecase(rr repository.RssEnrtyRepository) RssEntriesUsecase {
	return &rssEntriesUsecase{rr: rr}
}

func (f rssEntriesUsecase) CheckNewEntries(s []model.Subscription) []model.RssEntry {
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
			res = append(res, model.RssEntry{
				RSSURL:      sub.RSSURL,
				EntryTitle:  item.Title,
				EntryLink:   item.Link,
				PublishedAt: *item.PublishedParsed,
			})
		}
	}
	cpRes := make([]model.RssEntry, 0, len(s))
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

func diff[T comparable](s1, s2 []T) []T {
	diffSlice := []T{}
	cmpMap := map[T]int{}

	// slice2が各要素が何個あるのかmapに格納
	for _, v := range s2 {
		cmpMap[v] += 1
	}

	// slice2にある要素をslice1にあるか確認して、なければdiffSliceに格納
	for _, v := range s1 {
		t, ok := cmpMap[v]
		// slice1の要素がslice2になければ、配列に格納
		if !ok {
			diffSlice = append(diffSlice, v)
			continue
		}
		if t == 1 {
			delete(cmpMap, v)
		} else {
			cmpMap[v] -= 1
		}
	}
	return diffSlice
}
