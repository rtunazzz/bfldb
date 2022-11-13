package bfldb

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type aggregatedNickname struct {
	Nickname string   // Nickmane this struct belongs to
	UIDs     []string // User IDs associated with the nickname
}

// NicknamesToUIDs gets a list of UIDs for the nicknames provided.
// Returns a map with nicknames mapped to the UIDs and also any errors that might've occured.
//
// It fires up one goroutine for each nickname and fetches the UIDs.
func NicknamesToUIDs(pCtx context.Context, nicks []string) (map[string][]string, error) {
	idC := make(chan aggregatedNickname)
	g, ctx := errgroup.WithContext(pCtx)

	for _, n := range nicks {
		n := n
		g.Go(func() error {
			res, err := SearchNickname(ctx, n)

			if err == nil {

				// map the response to the UIDs only
				nIds := make([]string, 0, len(res.Data))
				for _, data := range res.Data {
					nIds = append(nIds, data.EncryptedUID)
				}

				idC <- aggregatedNickname{Nickname: n, UIDs: nIds}
			}

			return err
		})
	}

	// close the channel once the requests are done
	go func() {
		g.Wait()
		close(idC)
	}()

	// agregate it into a map of usernames to the ids
	aRes := make(map[string][]string, len(nicks))
	for id := range idC {
		aRes[id.Nickname] = id.UIDs
	}

	return aRes, g.Wait()
}
