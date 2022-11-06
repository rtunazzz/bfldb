package ftl

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
				u.log.Printf("[%s] Checking for new positions\n", u.UID)
				res, err := u.GetOtherPosition()
				if err != nil {
					ce <- fmt.Errorf("failed to fetch positions: %w", err)
					continue
				}

				if !res.Success {
					ce <- fmt.Errorf("failed to fetch positions: %v", res.Message)
					continue
				}

				u.log.Printf("[%s] Updating %d positions\n", u.UID, len(res.Data.OtherPositionRetList))
				u.handlePositions(res.Data.OtherPositionRetList, cp, ce)
			}
		}
	}()

	return cp, ce
}

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

		pp, ok := u.poss[h]
		if !ok {
			// no previous position, so it's a new one
			p.Type = Opened
			u.log.Printf("[%s] New opened position: %s %f %s @ %f\n", u.UID, p.Direction, p.Amount, p.Ticker, p.EntryPrice)
		} else if pp.Amount > p.Amount {
			// previously saved amount is BIGGER than the current amount
			// meaning the amount has DECREASED thus the position
			// has been partially closed

			u.log.Printf("[%s] Position partially closed: %s %f -> %f %s @ %f\n", u.UID, p.Direction, pp.Amount, p.Amount, p.Ticker, p.EntryPrice)
			p.Type = PartiallyClosed
		} else if pp.Amount < p.Amount {
			// previously saved amount is SMALLER than the current amount
			// meaning the amount has INCREASED thus the position
			// has been added to

			u.log.Printf("[%s] Position added to: %s %f -> %f %s @ %f\n", u.UID, p.Direction, pp.Amount, p.Amount, p.Ticker, p.EntryPrice)
			p.Type = AddedTo
		} else {
			// nothing changed
			continue
		}

		// something changed

		// dont send the new position on first run (bc it's not really "new")
		if !u.ff {
			cp <- p
		}

		// update the old position to the current one
		u.poss[h] = p
	}

	// TODO: rework this logic so we aren't looping twice over the maps when it could be done in one loop
	for h, p := range u.poss {
		if _, ok := used[h]; ok {
			continue
		}

		// position hasn't been updated (is not present in the leaderboard anymore)
		// thus it has been closed

		p.Type = Closed

		// dont send a new position on first run
		if !u.ff {
			cp <- p
		}
	}

	// set first run to false because we just completed it
	if u.ff {
		u.ff = false
	}
}
