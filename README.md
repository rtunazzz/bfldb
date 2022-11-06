# binance-ftl
> Wrapper around Binance's Futures Leaderboard API, in Go.]

[![GoDoc](https://godoc.org/github.com/rtunazzz/binance-ftl?status.svg)](https://godoc.org/github.com/rtunazzz/binance-ftl)

*THIS IS STILL WORK IN PROGRESS*

# Installation
```bash
go get -u github.com/rtunazzz/binance-ftl
```

## TODO
- [ ] Add test for `handlePositions`
- [ ] Add CI
- [ ] Complete TODO's

# Example

```golang
package main

import (
	"context"
	"fmt"

	"github.com/rtunazzz/binance-ftl"
)

func main() {
	// You can find this UID (encryptedUid) in the end of a leaderboard profile URL. For example:
	// https://www.binance.com/en/futures-activity/leaderboard/user?encryptedUid=47E6D002EBB1173967A6561F72B9395C
	u := ftl.NewUser("47E6D002EBB1173967A6561F72B9395C")
	cp, ce := u.SubscribePositions(context.Background())

	for {
		select {
		case position := <-cp:
			// Handle the new position as you need... Send a notification, copytrade...
			fmt.Printf("new position: %+v\n", position)
		case err := <-ce:
			fmt.Println(error has occured:", err)
			break
		}
	}
}
```

