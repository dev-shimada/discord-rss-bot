package fetch_test

import (
	"testing"

	"github.com/dev-shimada/discord-rss-bot/infrastructure/fetch"
)

func TestRssFetch(t *testing.T) {
	got, err := fetch.NewRss().Fetch("https://example.com/")
	if err == nil {
		t.Errorf("want: error, got: nil")
	}
	if got != nil {
		t.Errorf("want: nil, got: %v", got)
	}
}
