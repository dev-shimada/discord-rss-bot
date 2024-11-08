package usecase_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/domain/repository"
	"github.com/dev-shimada/discord-rss-bot/infrastructure/database"
	"github.com/dev-shimada/discord-rss-bot/infrastructure/persistence"
	"github.com/dev-shimada/discord-rss-bot/usecase"
	"github.com/google/go-cmp/cmp"
	"github.com/mmcdole/gofeed"
)

// mockRss is a mock of RssFetcher interface
type mockRss struct {
	mockFetch func() ([]*gofeed.Item, error)
}

func (m mockRss) Fetch(rssURL string) ([]*gofeed.Item, error) {
	return m.mockFetch()
}

// mockRssEnrtyRepository is a mock of RssEnrtyRepository interface
type mockRssEnrtyRepository struct {
	repository.RssEnrtyRepository
}

func (r mockRssEnrtyRepository) Create(_ []model.RssEntry) error          { return nil }
func (r mockRssEnrtyRepository) Find(_ []model.RssEntry) []model.RssEntry { return nil }

func TestCheck(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		args  model.Subscription
		fetch func() ([]*gofeed.Item, error)
		want  model.RssEntry
	}{
		{
			name: "empty",
			args: model.Subscription{},
			fetch: func() ([]*gofeed.Item, error) {
				return []*gofeed.Item{}, nil
			},
			want: model.RssEntry{},
		},
		{
			name: "multiple entries",
			args: model.Subscription{RSSURL: "https://example.com"},
			fetch: func() ([]*gofeed.Item, error) {
				return []*gofeed.Item{
					{Link: "https://example.com/entry1", Title: "title1", PublishedParsed: &now},
					{Link: "https://example.com/entry2", Title: "title2", PublishedParsed: &now},
				}, nil
			},
			want: model.RssEntry{RSSURL: "https://example.com", EntryTitle: "title1", EntryLink: "https://example.com/entry1", PublishedAt: now},
		},
		{
			name: "fetch error",
			args: model.Subscription{RSSURL: "https://example.com"},
			fetch: func() ([]*gofeed.Item, error) {
				return []*gofeed.Item{}, errors.New("error")
			},
			want: model.RssEntry{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rr := mockRssEnrtyRepository{}
			m := mockRss{tt.fetch}
			f := usecase.NewRssEntriesUsecase(rr, m)

			// test
			got := f.Check(tt.args)

			// assert
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestCheckNewEntries(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		args  []model.Subscription
		fetch func() ([]*gofeed.Item, error)
		want  []model.RssEntry
	}{
		{
			name: "empty",
			args: []model.Subscription{},
			fetch: func() ([]*gofeed.Item, error) {
				return []*gofeed.Item{}, nil
			},
			want: []model.RssEntry{},
		},
		{
			name: "new entries",
			args: []model.Subscription{{ID: 1, ChannelID: "123", RSSURL: "https://example.com", CreatedAt: now}},
			fetch: func() ([]*gofeed.Item, error) {
				return []*gofeed.Item{
					{Link: "https://example.com/entry1", Title: "title1", PublishedParsed: &now},
					{Link: "https://example.com/entry2", Title: "title2", PublishedParsed: &now},
				}, nil
			},
			want: []model.RssEntry{
				{ID: 1, RSSURL: "https://example.com", EntryTitle: "title1", EntryLink: "https://example.com/entry1", PublishedAt: now},
				{ID: 2, RSSURL: "https://example.com", EntryTitle: "title2", EntryLink: "https://example.com/entry2", PublishedAt: now},
			},
		},
		{
			name: "drop old entries",
			args: []model.Subscription{{ID: 1, ChannelID: "123", RSSURL: "https://example.com", CreatedAt: now}},
			fetch: func() ([]*gofeed.Item, error) {
				old := now.Add(-time.Microsecond)
				return []*gofeed.Item{
					{Link: "https://example.com/entry1", Title: "title", PublishedParsed: &old},
					{Link: "https://example.com/entry2", Title: "title", PublishedParsed: &now},
				}, nil
			},
			want: []model.RssEntry{{ID: 1, RSSURL: "https://example.com", EntryTitle: "title", EntryLink: "https://example.com/entry2", PublishedAt: now}},
		},
		{
			name: "fetch error",
			args: []model.Subscription{{ID: 1, ChannelID: "123", RSSURL: "https://example.com", CreatedAt: now}},
			fetch: func() ([]*gofeed.Item, error) {
				return []*gofeed.Item{}, errors.New("error")
			},
			want: []model.RssEntry{},
		},
	}

	bfDbPath := os.Getenv("DB_PATH")
	os.Setenv("DB_PATH", "testdata/test.db")
	defer os.Setenv("DB_PATH", bfDbPath)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			os.Remove("testdata/test.db")
			db := database.NewDB()
			defer database.CloseDB(db)
			rr := persistence.NewRssEntryPersistence(db)
			m := mockRss{tt.fetch}
			f := usecase.NewRssEntriesUsecase(rr, m)

			// test
			got := f.CheckNewEntries(tt.args)

			// remove CreatedAt field
			for i := range got {
				got[i].CreatedAt = time.Time{}
			}

			// assert
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestDiff(t *testing.T) {
	type args struct {
		oldEntries []model.RssEntry
		newEntries []model.RssEntry
	}
	tests := []struct {
		name string
		args args
		want []model.RssEntry
	}{
		{
			name: "new",
			args: args{
				oldEntries: []model.RssEntry{{EntryLink: "https://old.example.com"}},
				newEntries: []model.RssEntry{{EntryLink: "https://new.example.com"}, {EntryLink: "https://old.example.com"}},
			},
			want: []model.RssEntry{{EntryLink: "https://new.example.com"}},
		},
		{
			name: "same",
			args: args{
				oldEntries: []model.RssEntry{{EntryLink: "https://old.example.com"}},
				newEntries: []model.RssEntry{{EntryLink: "https://old.example.com"}},
			},
			want: []model.RssEntry{},
		},
		{
			name: "old",
			args: args{
				oldEntries: []model.RssEntry{{EntryLink: "https://old.example.com"}, {EntryLink: "https://old2.example.com"}},
				newEntries: []model.RssEntry{{EntryLink: "https://old.example.com"}},
			},
			want: []model.RssEntry{},
		},
		{
			name: "empty",
			args: args{
				oldEntries: []model.RssEntry{{EntryLink: "https://old.example.com"}},
				newEntries: []model.RssEntry{},
			},
			want: []model.RssEntry{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := usecase.Diff(tt.args.newEntries, tt.args.oldEntries)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type args struct {
		entries []model.RssEntry
	}
	tests := []struct {
		name string
		args args
		want []model.RssEntry
	}{
		{
			name: "unique",
			args: args{
				entries: []model.RssEntry{{EntryLink: "https://example.com"}, {EntryLink: "https://example.com"}},
			},
			want: []model.RssEntry{{EntryLink: "https://example.com"}},
		},
		{
			name: "empty",
			args: args{
				entries: []model.RssEntry{},
			},
			want: []model.RssEntry{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := usecase.Unique(tt.args.entries)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}
