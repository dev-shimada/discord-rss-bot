package persistence_test

import (
	"os"
	"testing"
	"time"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/infrastructure/database"
	"github.com/dev-shimada/discord-rss-bot/infrastructure/persistence"
	"github.com/google/go-cmp/cmp"
)

func TestRssEntryPersistenceCreate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		args []model.RssEntry
		want []model.RssEntry
	}{
		{
			name: "empty",
			args: []model.RssEntry{},
			want: []model.RssEntry{},
		},
		{
			name: "multiple",
			args: []model.RssEntry{
				{RSSURL: "https://example.com/", EntryTitle: "title1", EntryLink: "https://example.com/entry1", PublishedAt: now},
				{RSSURL: "https://example.com/", EntryTitle: "title2", EntryLink: "https://example.com/entry2", PublishedAt: now},
			},
			want: []model.RssEntry{
				{ID: 1, RSSURL: "https://example.com/", EntryTitle: "title1", EntryLink: "https://example.com/entry1", PublishedAt: now, CreatedAt: time.Time{}},
				{ID: 2, RSSURL: "https://example.com/", EntryTitle: "title2", EntryLink: "https://example.com/entry2", PublishedAt: now, CreatedAt: time.Time{}},
			},
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
			r := persistence.NewRssEntryPersistence(db)

			// test
			err := r.Create(tt.args)
			got := []model.RssEntry{}
			db.Find(&got)

			// remove CreatedAt
			for i := range got {
				got[i].CreatedAt = time.Time{}
			}

			// assert
			if err != nil {
				t.Errorf("error: %v", err)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestRssEntryPersistenceFind(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		args []model.RssEntry
		want []model.RssEntry
	}{
		{
			name: "empty",
			args: []model.RssEntry{},
			want: []model.RssEntry{},
		},
		{
			name: "multiple",
			args: []model.RssEntry{{ID: 1}},
			want: []model.RssEntry{
				{ID: 1, RSSURL: "https://example.com/", EntryTitle: "title1", EntryLink: "https://example.com/entry1", PublishedAt: now},
				{ID: 2, RSSURL: "https://example.com/", EntryTitle: "title2", EntryLink: "https://example.com/entry2", PublishedAt: now},
			},
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
			r := persistence.NewRssEntryPersistence(db)

			// prepare
			db.Create(tt.want)

			// test
			got := r.Find(tt.args)

			// assert
			if !cmp.Equal(got, tt.want) {
				t.Errorf("diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}
