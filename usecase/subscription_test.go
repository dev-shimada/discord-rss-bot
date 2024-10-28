package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/dev-shimada/discord-rss-bot/domain/model"
	"github.com/dev-shimada/discord-rss-bot/domain/repository"
	"github.com/dev-shimada/discord-rss-bot/usecase"
	"github.com/google/go-cmp/cmp"
)

type mockSubscription struct {
	repository.SubscriptionRepository
	mockCreate func() error
}

func (m mockSubscription) Create(sub model.Subscription) error {
	return m.mockCreate()
}

func TestCreateSubscription(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name   string
		args   model.Subscription
		create func() error
		want   string
	}{
		{
			name: "success",
			args: model.Subscription{ID: 1, ChannelID: "123", RSSURL: "https://example.com", CreatedAt: now},
			create: func() error {
				return nil
			},
			want: "Successfully subscribed to RSS feed.",
		},
		{
			name: "error",
			args: model.Subscription{ID: 1, ChannelID: "123", RSSURL: "https://example.com", CreatedAt: now},
			create: func() error {
				return errors.New("error")
			},
			want: "Failed to subscribe to RSS feed.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			sr := mockSubscription{mockCreate: tt.create}
			s := usecase.NewSubscriptionUsecase(sr)

			// test
			got := s.Create(tt.args)

			// assert
			if !cmp.Equal(got, tt.want) {
				t.Errorf("Diff: %v", cmp.Diff(got, tt.want))
			}
		})
	}
}
