package bfldb

import (
	"errors"
	"math"
)

var (
	ErrNoPreviousPosition = errors.New("no previous position amount")
)

// PositionType is a type of a position
type PositionType int

const (
	Opened          PositionType = iota + 1 // A completely new position
	Closed                                  // A completely closed position
	AddedTo                                 // A new position where there previously already was a position for the same direction and ticker + the amount increased
	PartiallyClosed                         // A new position where there previously already was a position for the same direction and ticker + the amount decreased
)

func (pt PositionType) String() string {
	switch pt {
	default:
		return ""
	case Opened:
		return "opened"
	case Closed:
		return "closed"
	case AddedTo:
		return "added to"
	case PartiallyClosed:
		return "partially closed"
	}
}

// Position represents a position user is in.
type Position struct {
	Type       PositionType   // Type of the position
	Direction  TradeDirection // Direction (e.g. LONG / SHORT)
	Ticker     string         // Ticker of the position (e.g. BTCUSDT)
	EntryPrice float64        // Entry price
	MarkPrice  float64        // Entry price
	Amount     float64        // Amount
	PrevAmount float64        // previous amount, used for determining position type
	Leverage   int            // Position leverage
	Pnl        float64        // PNL
	Roe        float64        // ROE
}

// ToOrder converts a position into an Order.
func (p Position) ToOrder() Order {
	o := Order{
		Ticker: p.Ticker,

		ReduceOnly: false,
		Direction:  p.Direction,
		Amount:     p.Amount,
		Leverage:   p.Leverage,
	}

	if p.Type == Closed || p.Type == PartiallyClosed {
		o.ReduceOnly = true

		// inverse the direction since the new order should be closing the
		// position
		if p.Direction == Long {
			o.Direction = Short
		} else {
			o.Direction = Long
		}
	}

	// opened = nothing changes

	if p.Type == Closed {
		o.Amount = p.PrevAmount
	}

	if p.Type == PartiallyClosed {
		o.Amount = p.PrevAmount - p.Amount
	}

	if p.Type == AddedTo {
		o.Amount = p.Amount - p.PrevAmount
	}

	return o
}

// DeterminePositionType determines the type of a position,
// based on the current and previous position size.
func DeterminePositionType(amt float64, prevAmt float64) PositionType {
	if prevAmt == 0 {
		// no previous amount, so it's a freshly opened position
		return Opened
	}

	if prevAmt < amt {
		// amount increased
		return AddedTo
	}

	if prevAmt > amt {
		// amount decreased
		return PartiallyClosed
	}

	if amt == 0 {
		// no amount, so position is closed
		return Closed
	}

	// should never get here...
	return 0
}

// newPosition creates a new Position from a rawPosition
func newPosition(rp rawPosition) Position {
	// Amount is negative on short positions

	dir := Long
	if rp.Amount < 0 {
		dir = Short
		rp.Amount = math.Abs(rp.Amount)
	}

	return Position{
		Direction:  dir,
		Ticker:     rp.Symbol,
		EntryPrice: rp.EntryPrice,
		MarkPrice:  rp.MarkPrice,
		Amount:     rp.Amount,
		Leverage:   rp.Leverage,
		Pnl:        rp.Pnl,
		Roe:        rp.Roe,
	}
}
