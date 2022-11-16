package bfldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func chanToArrays(cp chan Position, ce chan error) (pos []Position, errs []error) {
	select {
	case p := <-cp:
		pos = append(pos, p)
	case err := <-ce:
		errs = append(errs, err)
	}

	return
}

func TestLogic(t *testing.T) {
	rp1 := rawPosition{
		Symbol:          "SUSHIUSDT",
		EntryPrice:      1.886,
		MarkPrice:       1.85843264,
		Pnl:             -248.82299136,
		Roe:             -0.02261789,
		Amount:          9026,
		UpdateTimeStamp: 1667674507457,
		Leverage:        2,
	}
	p1 := newPosition(rp1)
	p1.Type = Opened

	p1C := p1
	p1C.Amount = 0
	p1C.setType(p1.Amount)

	rp1Added := rawPosition{
		Symbol:          "SUSHIUSDT",
		EntryPrice:      1.886,
		MarkPrice:       1.85843264,
		Pnl:             -248.82299136,
		Roe:             -0.02261789,
		Amount:          rp1.Amount + 1,
		UpdateTimeStamp: 1667674507457,
		Leverage:        2,
	}

	p1Added := newPosition(rp1Added)
	p1Added.setType(p1.Amount)

	tests := []struct {
		initPoss []rawPosition
		rawPoss  []rawPosition
		outPos   []Position
		endPos   []Position
		msg      string
	}{
		{
			msg:      "no initial positions",
			initPoss: []rawPosition{},
			rawPoss: []rawPosition{
				rp1,
			},
			outPos: []Position{
				p1,
			},
			endPos: []Position{p1},
		},
		{
			msg:      "added to position",
			initPoss: []rawPosition{rp1},
			rawPoss: []rawPosition{
				rp1Added,
			},
			outPos: []Position{
				p1Added,
			},
			endPos: []Position{p1Added},
		},
		{
			msg:      "position closed",
			initPoss: []rawPosition{rp1},
			rawPoss:  []rawPosition{},
			outPos:   []Position{p1C},
			endPos:   []Position{},
		},
	}

	for _, tt := range tests {
		u := NewUser("47E6D002EBB1173967A6561F72B9395C", WithLogging())
		cp := make(chan Position)
		ce := make(chan error)

		// load in initial positions
		u.handlePositions(tt.initPoss, cp, ce)
		t.Log("init positions:", u.pHashes)

		go func() {
			// handle positions
			u.handlePositions(tt.rawPoss, cp, ce)
		}()

		ops, errs := chanToArrays(cp, ce)
		t.Logf("positions: %+v\n", ops)
		t.Logf("errors: %+v\n", errs)

		t.Log("end positions:", u.pHashes)

		ep := make([]Position, 0, len(u.pHashes))
		for _, p := range u.pHashes {
			ep = append(ep, p)
		}

		close(cp)
		close(ce)

		assert.Equal(t, 0, len(errs), "there were some errors")
		assert.EqualValues(t, tt.outPos, ops, "expected different output positions for test "+tt.msg)
		assert.EqualValues(t, tt.endPos, ep, "expected different end positions for test "+tt.msg)
	}
}
