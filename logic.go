package bfldb

import (
	"context"
	"fmt"
	"time"
)

// SubscribePositions subscribes to user's potition details in a new goroutine.
//
// Returns two read-only channels, one with user's positions, other with any errors occured during the subsription.
func (u *User) SubscribePositions(ctx context.Context) (<-chan Position, <-chan error) {
	cp := make(chan Position)
	ce := make(chan error)

	go func() {
		defer close(cp)
		defer close(ce)

		for {
			select {
			case <-ctx.Done():
				return

			default:
				// u.log.Printf("[%s] Checking for new positions\n", u.id)
				res, err := u.GetOtherPosition(ctx)
				if err != nil {
					ce <- fmt.Errorf("failed to fetch positions: %w", err)
					time.Sleep(u.Delay())
					continue
				}

				if !res.Success {
					ce <- fmt.Errorf("failed to fetch positions, bad response message: %v", res.Message)
					time.Sleep(u.Delay())
					continue
				}

				// u.log.Printf("[%s] Updating %d positions\n", u.id, len(res.Data.OtherPositionRetList))
				u.handlePositions(res.Data.OtherPositionRetList, cp, ce)
				time.Sleep(u.Delay())
			}
		}
	}()

	return cp, ce
}

// handlePositions parses raw positions, determines their type and sends the new ones through a channel.
func (u *User) handlePositions(rps []rawPosition, cp chan<- Position, ce chan<- error) {
	// used will be used for checking whether or not a position was already handled
	// (thus if it's a new position or if it hasn't been present in the latest fetch and thus been closed)
	used := make(map[string]struct{}, len(rps))

	for _, rp := range rps {
		p := newPosition(rp)

		// mark as used
		used[p.Ticker] = struct{}{}

		// retrieve old position
		pp := u.positions[p.Ticker]

		// amount is the same, so we dont want to send the update
		if pp.Amount == p.Amount {

			// update the values that change on every refresh
			pp.MarkPrice = p.MarkPrice
			pp.Pnl = p.Pnl
			pp.Roe = p.Roe

			u.positions[p.Ticker] = pp

			continue
		}

		// record the previous amount on the new position
		p.PrevAmount = pp.Amount

		// determine the current position type and assign
		p.Type = DeterminePositionType(p.Amount, pp.Amount)

		u.log.Printf("[%s] {send: %t} Position change: %d %s %f -> %f %s @ %f\n", u.UID, !u.firstFetch, p.Type, p.Direction, p.PrevAmount, p.Amount, p.Ticker, p.EntryPrice)

		// dont send the new position on first run (bc it's not really "new")
		if !u.firstFetch {
			cp <- p
		}

		// add/update the old position to the current one
		u.positions[p.Ticker] = p
	}

	// check which positions were not present in the latest fetch
	for h, p := range u.positions {
		// position still in the leaderboard
		if _, ok := used[h]; ok {
			continue
		}

		// position hasn't been updated (is not present in the leaderboard anymore)
		// thus it has been closed

		p.Type = Closed
		p.PrevAmount = p.Amount
		p.Amount = 0

		cp <- p

		// remove the position from user's positions
		delete(u.positions, h)
	}

	// set first run to false because we just completed it
	if u.firstFetch {
		u.firstFetch = false
	}
}
