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
		t := time.NewTicker(u.d)
		defer t.Stop()

		defer close(cp)
		defer close(ce)

		for {
			select {
			case <-ctx.Done():
				return

			case <-t.C:
				u.log.Printf("[%s] Checking for new positions\n", u.id)
				res, err := u.GetOtherPosition(ctx)
				if err != nil {
					ce <- fmt.Errorf("failed to fetch positions: %w", err)
					continue
				}

				if !res.Success {
					ce <- fmt.Errorf("failed to fetch positions, bad response message: %v", res.Message)
					continue
				}

				u.log.Printf("[%s] Updating %d positions\n", u.id, len(res.Data.OtherPositionRetList))
				u.handlePositions(res.Data.OtherPositionRetList, cp, ce)
			}
		}
	}()

	return cp, ce
}

// handlePositions parses raw positions, determines their type and sends the new ones through a channel.
func (u *User) handlePositions(rps []rawPosition, cp chan<- Position, ce chan<- error) {
	// used will be used for checking whether or not a position was already handled
	// (thus if it's a new position or if it hasn't been present in the latest fetch and thus been closed)
	used := make(map[uint64]struct{}, len(rps))

	for _, rp := range rps {
		p := parsePosition(rp)

		// check if there are any old positions or if it's a new one
		h, err := p.hash()
		if err != nil {
			ce <- fmt.Errorf("failed to hash a %q position: %w", p.Ticker, err)
			continue
		}

		// mark as used
		used[h] = struct{}{}

		pp, ok := u.pHashes[h]

		// if there's a record of the same position,
		// and the amount is the same, then nothing changed so skip
		if ok && pp.Amount == p.Amount {
			continue
		}

		// it's ok on !ok because pp will just be a Position{}
		p.setType(pp)

		u.log.Printf("[%s] Position change: %d %s %f %s @ %f\n", u.id, p.Type, p.Direction, p.Amount, p.Ticker, p.EntryPrice)

		// dont send the new position on first run (bc it's not really "new")
		if !u.isFirst {
			cp <- p
		}

		// update the old position to the current one
		u.pHashes[h] = p
	}

	// check which positions were not present in the latest fetch

	for h, p := range u.pHashes {
		if _, ok := used[h]; ok {
			continue
		}

		// position hasn't been updated (is not present in the leaderboard anymore)
		// thus it has been closed

		p.Type = Closed
		cp <- p

		// remove the position from user's positions
		delete(u.pHashes, h)
	}

	// set first run to false because we just completed it
	if u.isFirst {
		u.isFirst = false
	}
}
