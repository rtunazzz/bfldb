package ftl

import (
	"golang.org/x/sync/errgroup"
)

// NicknamesToUIDs converts a list of nicknames to a list of UIDs,
// returns any errors that might've occured.
//
// It fires up one goroutine for each nickname and fetches the UIDs.
func NicknamesToUIDs(nicks []string) ([]string, error) {
	var g errgroup.Group
	idC := make(chan string)

	for _, n := range nicks {
		n := n
		g.Go(func() error {
			res, err := SearchNickname(n)
			if err == nil {
				for _, data := range res.Data {
					idC <- data.EncryptedUID
				}
			}
			return err
		})
	}

	// close the channel once the requests are done
	go func() {
		g.Wait()
		close(idC)
	}()

	// agregate it into a slice
	uids := make([]string, 0, len(nicks))
	for id := range idC {
		uids = append(uids, id)
	}

	return uids, g.Wait()
}
