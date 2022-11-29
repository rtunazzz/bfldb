package bfldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			name: "closed position",
			pp:   Position{Amount: 1},
			np:   Position{},
			et:   Closed,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.np.setType(tt.pp.Amount)
			require.Equal(t, tt.et, tt.np.Type, "type missmatch")
			require.Equal(t, tt.pp.Amount, tt.np.PrevAmount, "previous amount missmatch")
		})
	}
}

func TestToOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		p    Position
		want Order
	}{
		{
			name: "opened position",
			p:    Position{Direction: Long, Amount: 1, PrevAmount: 0, Type: Opened},
			want: Order{Direction: Long, Amount: 1, ReduceOnly: false},
		},
		{
			name: "closed position",
			p:    Position{Direction: Long, Amount: 0, PrevAmount: 1, Type: Closed},
			want: Order{Direction: Short, Amount: 1, ReduceOnly: true},
		},
		{
			name: "added to position",
			p:    Position{Direction: Long, Amount: 1, PrevAmount: 0.5, Type: AddedTo},
			want: Order{Direction: Long, Amount: 0.5, ReduceOnly: false},
		},
		{
			name: "partially closed position",
			p:    Position{Direction: Long, Amount: 0.1, PrevAmount: 1, Type: PartiallyClosed},
			want: Order{Direction: Short, Amount: 0.9, ReduceOnly: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.ToOrder()
			require.Equalf(t, tt.want, got, "Position.ToOrder() = %v, want %v", got, tt.want)
		})
	}
}
