package persistence_test

import (
	"os"
	"testing"
	"time"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/infrastructure/database"
	"github.com/dev-shimada/discord-rss-bot/infrastructure/persistence"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
)

func TestSubscriptionPersistenceCreate(t *testing.T) {
	test := []struct {
		name string
		args model.Subscription
		want model.Subscription
	}{
		{
			name: "success",
			args: model.Subscription{ChannelID: "1234567890", RSSURL: "https://example.com"},
			want: model.Subscription{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: time.Time{}},
		},
	}

	bfDbPath := os.Getenv("DB_PATH")
	os.Setenv("DB_PATH", "testdata/test.db")
	defer os.Setenv("DB_PATH", bfDbPath)

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			os.Remove("testdata/test.db")
			db := database.NewDB()
			defer database.CloseDB(db)
			sr := persistence.NewSubscriptionPersistence(db)

			// test
			err := sr.Create(tt.args)

			got := []model.Subscription{}
			db.Find(&got)

			for i := range got {
				got[i].CreatedAt = time.Time{}
			}

			// assert
			if err != nil {
				t.Errorf("want: nil, got: %v", err)
			}
			if len(got) != 1 {
				t.Errorf("want: 1, got: %d", len(got))
			}
			if !cmp.Equal(got[0], tt.want) {
				t.Errorf("Diff: %v", cmp.Diff(got[0], tt.want))
			}
		})
	}
}

func TestSubscriptionPersistenceFindByModel(t *testing.T) {
	now := time.Now()
	test := []struct {
		name   string
		args   model.Subscription
		create func(*gorm.DB)
		want   []model.Subscription
	}{
		{
			name:   "empty",
			args:   model.Subscription{},
			create: func(db *gorm.DB) {},
			want:   []model.Subscription{},
		},
		{
			name: "select by id",
			args: model.Subscription{ID: 1},
			create: func(db *gorm.DB) {
				db.Create(&model.Subscription{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now})
				db.Create(&model.Subscription{ID: 2, ChannelID: "0987654321", RSSURL: "https://example.com", CreatedAt: now})
			},
			want: []model.Subscription{
				{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now},
			},
		},
		{
			name: "select by ChannelID",
			args: model.Subscription{ChannelID: "1234567890"},
			create: func(db *gorm.DB) {
				db.Create(&model.Subscription{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now})
				db.Create(&model.Subscription{ID: 2, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now})
				db.Create(&model.Subscription{ID: 3, ChannelID: "0987654321", RSSURL: "https://example.com", CreatedAt: now})
			},
			want: []model.Subscription{
				{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now},
				{ID: 2, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now},
			},
		},
		{
			name: "select by ChannelID and RSSURL",
			args: model.Subscription{ChannelID: "1234567890", RSSURL: "https://example.com"},
			create: func(db *gorm.DB) {
				db.Create(&model.Subscription{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now})
				db.Create(&model.Subscription{ID: 2, ChannelID: "0987654321", RSSURL: "https://example.com", CreatedAt: now})
				db.Create(&model.Subscription{ID: 3, ChannelID: "1234567890", RSSURL: "https://example.com/exclude", CreatedAt: now})
				db.Create(&model.Subscription{ID: 3, ChannelID: "0987654321", RSSURL: "https://example.com/exclude", CreatedAt: now})
			},
			want: []model.Subscription{
				{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now},
			},
		},
	}

	bfDbPath := os.Getenv("DB_PATH")
	os.Setenv("DB_PATH", "testdata/test.db")
	defer os.Setenv("DB_PATH", bfDbPath)

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			os.Remove("testdata/test.db")
			db := database.NewDB()
			defer database.CloseDB(db)
			sr := persistence.NewSubscriptionPersistence(db)

			// prepare
			tt.create(db)

			// test
			got, err := sr.FindByModel(tt.args)

			// assert
			if err != nil {
				t.Errorf("want: nil, got: %v", err)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestSubscriptionPersistenceFindAll(t *testing.T) {
	now := time.Now()
	test := []struct {
		name   string
		create func(*gorm.DB)
		want   []model.Subscription
	}{
		{
			name:   "empty",
			create: func(db *gorm.DB) {},
			want:   []model.Subscription{},
		},
		{
			name: "multiple",
			create: func(db *gorm.DB) {
				db.Create(&model.Subscription{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now})
				db.Create(&model.Subscription{ID: 2, ChannelID: "0987654321", RSSURL: "https://example.com", CreatedAt: now})
			},
			want: []model.Subscription{
				{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now},
				{ID: 2, ChannelID: "0987654321", RSSURL: "https://example.com", CreatedAt: now},
			},
		},
	}

	bfDbPath := os.Getenv("DB_PATH")
	os.Setenv("DB_PATH", "testdata/test.db")
	defer os.Setenv("DB_PATH", bfDbPath)

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			os.Remove("testdata/test.db")
			db := database.NewDB()
			defer database.CloseDB(db)
			sr := persistence.NewSubscriptionPersistence(db)

			// prepare
			tt.create(db)

			// test
			got, err := sr.FindAll()

			// assert
			if err != nil {
				t.Errorf("want: nil, got: %v", err)
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestSubscriptionPersistenceDelete(t *testing.T) {
	now := time.Now()
	test := []struct {
		name    string
		args    model.Subscription
		create  func(*gorm.DB)
		want    []model.Subscription
		withErr bool
	}{
		{
			name: "success",
			args: model.Subscription{ID: 1},
			create: func(db *gorm.DB) {
				db.Create(&model.Subscription{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now})
				db.Create(&model.Subscription{ID: 2, ChannelID: "0987654321", RSSURL: "https://example.com", CreatedAt: now})
			},
			want: []model.Subscription{
				{ID: 2, ChannelID: "0987654321", RSSURL: "https://example.com", CreatedAt: now},
			},
			withErr: false,
		},
		{
			name: "record not found",
			args: model.Subscription{ID: 3},
			create: func(db *gorm.DB) {
				db.Create(&model.Subscription{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now})
				db.Create(&model.Subscription{ID: 2, ChannelID: "0987654321", RSSURL: "https://example.com", CreatedAt: now})
			},
			want: []model.Subscription{
				{ID: 1, ChannelID: "1234567890", RSSURL: "https://example.com", CreatedAt: now},
				{ID: 2, ChannelID: "0987654321", RSSURL: "https://example.com", CreatedAt: now},
			},
			withErr: true,
		},
	}

	bfDbPath := os.Getenv("DB_PATH")
	os.Setenv("DB_PATH", "testdata/test.db")
	defer os.Setenv("DB_PATH", bfDbPath)

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			os.Remove("testdata/test.db")
			db := database.NewDB()
			defer database.CloseDB(db)
			sr := persistence.NewSubscriptionPersistence(db)

			// prepare
			tt.create(db)

			// test
			err := sr.Delete(tt.args)

			got := []model.Subscription{}
			db.Find(&got)

			// assert
			if tt.withErr && err == nil {
				t.Errorf("want: error, got: nil")
			} else if !tt.withErr && err != nil {
				t.Errorf("want: nil, got: %v)", err)
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}
