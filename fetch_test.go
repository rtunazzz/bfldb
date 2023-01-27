package bfldb

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSearchNickname(t *testing.T) {
	tests := []struct {
		nickname  string
		assertion require.ErrorAssertionFunc
	}{
		{
			nickname:  "StellarMom",
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.nickname, func(t *testing.T) {
			_, err := SearchNickname(context.Background(), tt.nickname)
			tt.assertion(t, err)
		})
	}
}

func TestGetOtherLeaderboardBaseInfo(t *testing.T) {
	tests := []struct {
		uuid      string
		assertion require.ErrorAssertionFunc
	}{
		{
			uuid:      "3AFFCB67ED4F1D1D8437BA17F4E8E5ED",
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.uuid, func(t *testing.T) {
			_, err := GetOtherLeaderboardBaseInfo(context.Background(), tt.uuid)
			tt.assertion(t, err)
		})
	}
}

func TestUser_GetOtherLeaderboardBaseInfo(t *testing.T) {
	tests := []struct {
		uuid      string
		assertion require.ErrorAssertionFunc
	}{
		{
			uuid:      "3AFFCB67ED4F1D1D8437BA17F4E8E5ED",
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.uuid, func(t *testing.T) {
			u := NewUser(tt.uuid)
			_, err := u.GetOtherLeaderboardBaseInfo(context.Background())
			tt.assertion(t, err)
		})
	}
}

func TestGetOtherPosition(t *testing.T) {
	tests := []struct {
		uuid      string
		assertion require.ErrorAssertionFunc
	}{
		{
			uuid:      "3AFFCB67ED4F1D1D8437BA17F4E8E5ED",
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.uuid, func(t *testing.T) {
			_, err := GetOtherPosition(context.Background(), tt.uuid)
			tt.assertion(t, err)
		})
	}
}

func TestUser_GetOtherPosition(t *testing.T) {
	tests := []struct {
		uuid      string
		assertion require.ErrorAssertionFunc
	}{
		{
			uuid:      "3AFFCB67ED4F1D1D8437BA17F4E8E5ED",
			assertion: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.uuid, func(t *testing.T) {
			u := NewUser(tt.uuid)
			_, err := u.GetOtherPosition(context.Background())
			tt.assertion(t, err)
		})
	}
}
