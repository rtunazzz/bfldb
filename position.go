package bfldb

import (
	"github.com/mitchellh/hashstructure/v2"
)

// PositionDirection can be either LONG / SHORT
type PositionDirection int

const (
	Short PositionDirection = iota + 1
	Long
)

func (pd PositionDirection) String() string {
	if pd == Short {
		return "SHORT"
	}
	return "LONG"
}

// PositionType is a type of a position
type PositionType int

const (
	Opened          PositionType = iota + 1 // A completely new position
	Closed                                  // A completely closed position
	AddedTo                                 // A new position where there previously already was a position for the same direction and ticker + the amount increased
	PartiallyClosed                         // A new position where there previously already was a position for the same direction and ticker + the amount decreased
)

// Position represents a trade position.
type Position struct {
	// TODO: Verify that there can only be one position with the same ticker on Binance's leaderboard
	// so it's impossible to have for example one BTCUSDT LONG at 20k and ANOTHER at 30k
	// if that IS possible, adjust the hashing adequately

	Type       PositionType      `hash:"ignore"` // Type of the position
	Direction  PositionDirection // Direction of the position (e.g. LONG / SHORT)
	Ticker     string            // Ticker of the position (e.g. BTCUSDT)
	EntryPrice float64           `hash:"ignore"` // Entry price
	Amount     float64           `hash:"ignore"` // Amount
}

// parsePosition parseas a raw position into a position
func parsePosition(rp rawPosition) Position {
	return Position{
		Direction:  getPosDir(rp),
		Ticker:     rp.Symbol,
		EntryPrice: rp.EntryPrice,
		Amount:     rp.Amount,
	}
}

// hash hashes a position into an uint64
func (p Position) hash() (uint64, error) {
	return hashstructure.Hash(p, hashstructure.FormatV2, nil)
}

// getPosDir determines the position direction.
func getPosDir(rp rawPosition) PositionDirection {
	pd := Short

	// Long is either when:
	// - Entry is below mark price and PNL is positive
	// - Entry is above mark price and PNL is negative

	if rp.EntryPrice < rp.MarkPrice && rp.Pnl > 0 {
		pd = Long
	} else if rp.EntryPrice > rp.MarkPrice && rp.Pnl < 0 {
		pd = Long
	}

	return pd
}
