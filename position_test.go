package bfldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPosDir(t *testing.T) {
	tests := []struct {
		rp  rawPosition
		exp PositionDirection
		msg string
	}{
		{
			rp:  rawPosition{EntryPrice: 1000, MarkPrice: 2000, Pnl: 1000},
			exp: Long,
			msg: "should be long with positive PNL",
		},
		{
			rp:  rawPosition{EntryPrice: 1000, MarkPrice: 500, Pnl: -500},
			exp: Long,
			msg: "should be long with negative PNL",
		},
		{
			rp:  rawPosition{EntryPrice: 1000, MarkPrice: 500, Pnl: 500},
			exp: Short,
			msg: "should be short with positive PNL",
		},
		{
			rp:  rawPosition{EntryPrice: 1000, MarkPrice: 1500, Pnl: -500},
			exp: Short,
			msg: "should be short with negative PNL",
		},
	}

	for _, tc := range tests {
		dir := getPosDir(tc.rp)
		assert.Equal(t, tc.exp, dir, tc.msg)
	}
}

func TestHashEquality(t *testing.T) {
	tests := []struct {
		p1  Position
		p2  Position
		msg string
	}{
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			msg: "two exactly same positions",
		},
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: AddedTo, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20*10 ^ 3, Amount: 2},
			msg: "added to position at the same price",
		},
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: AddedTo, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 25 * 10e3, Amount: 2},
			msg: "added to position at a different price",
		},
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: PartiallyClosed, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 0.5},
			msg: "partially closed position",
		},
	}

	for _, tc := range tests {
		h1, err := tc.p1.hash()
		assert.Nil(t, err, "hashing errored out")

		h2, err := tc.p2.hash()
		assert.Nil(t, err, "hashing errored out")

		assert.Equal(t, h1, h2, tc.msg)
	}
}

func TestHashInequality(t *testing.T) {
	tests := []struct {
		p1  Position
		p2  Position
		msg string
	}{
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: Opened, Direction: Long, Ticker: "ETHSDT", EntryPrice: 20 * 10e3, Amount: 1},
			msg: "two differnt tickers",
		},
		{
			p1:  Position{Type: Opened, Direction: Long, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			p2:  Position{Type: Opened, Direction: Short, Ticker: "BTCUSDT", EntryPrice: 20 * 10e3, Amount: 1},
			msg: "two differnt directions",
		},
	}

	for _, tc := range tests {
		h1, err := tc.p1.hash()
		assert.Nil(t, err, "hashing errored out")

		h2, err := tc.p2.hash()
		assert.Nil(t, err, "hashing errored out")

		assert.NotEqual(t, h1, h2, tc.msg)
	}
}

func TestSetType(t *testing.T) {
	t.Parallel()

	type args struct {
	}
	tests := []struct {
		name string
		np   Position     // new position
		pp   Position     // previous position
		et   PositionType // expected type
	}{
		{
			name: "new opened position",
			pp:   Position{},
			np:   Position{Amount: 0.5},
			et:   Opened,
		},
		{
			name: "partially closed position",
			pp:   Position{Amount: 1},
			np:   Position{Amount: 0.5},
			et:   PartiallyClosed,
		},
		{
			name: "added to position",
			pp:   Position{Amount: 1},
			np:   Position{Amount: 1.5},
			et:   AddedTo,
		},
		{
			name: "same position",
			pp:   Position{Amount: 1, Type: Opened},
			np:   Position{Amount: 1},
			et:   Opened,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.np.setType(tt.pp)
			require.Equal(t, tt.et, tt.np.Type, "type missmatch")
		})
	}
}
