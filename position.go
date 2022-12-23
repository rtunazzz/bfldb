package bfldb

import (
	"errors"
	"math"

	"github.com/mitchellh/hashstructure/v2"
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
	PrevAmount float64 `hash:"ignore"` // previous amount, used for converting into an order

	Type       PositionType   `hash:"ignore"` // Type of the position
	Direction  TradeDirection // Direction (e.g. LONG / SHORT)
	Ticker     string         // Ticker of the position (e.g. BTCUSDT)
	EntryPrice float64        `hash:"ignore"` // Entry price
	MarkPrice  float64        `hash:"ignore"` // Entry price
	Amount     float64        `hash:"ignore"` // Amount
	Leverage   int            `hash:"ignore"` // Position leverage
	Pnl        float64        `hash:"ignore"` // PNL
	Roe        float64        `hash:"ignore"` // ROE
}

// ToOrder converts a position into an Order.
func (p Position) ToOrder() Order {
	// == 0 means no position type
	if p.Type == 0 {
		p.setType(p.PrevAmount)
	}

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

	if p.Type == Closed && p.Amount == 0 {
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

// hash hashes a position into an uint64
func (p Position) hash() (uint64, error) {
	return hashstructure.Hash(p, hashstructure.FormatV2, nil)
}

// setType sets the type of the position in accordance with the previous position.
//
// This is because Binance API only returns the CURRENT position details so
// we need to keep track of the previous position manually and detect changes that way.
func (p *Position) setType(pa float64) {
	p.PrevAmount = pa

	if pa == 0 {
		// no previous position, so it's a new one
		p.Type = Opened
	} else if pa > p.Amount {
		// previously saved amount is BIGGER than the current amount
		// meaning the amount has DECREASED thus the position
		// has been (partially) closed
		if p.Amount == 0 {
			p.Type = Closed
		} else {
			p.Type = PartiallyClosed
		}
	} else if pa < p.Amount {
		// previously saved amount is SMALLER than the current amount
		// meaning the amount has INCREASED thus the position
		// has been added to
		p.Type = AddedTo
	}
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
