# bfldb
> Wrapper around Binance's Futures Leaderboard API, in Go.

[![GoDoc](https://godoc.org/github.com/rtunazzz/bfldb?status.svg)](https://godoc.org/github.com/rtunazzz/bfldb)

<p align="center"><img width=100% src="https://user-images.githubusercontent.com/38296319/200170420-0644f467-49ff-4ecd-8811-1bc939f84fea.png"></p>

# Installation
```bash
go get -u github.com/rtunazzz/bfldb
```

# Example

## Subscribing to user's positions
> **The user needs to have their position sharing enabled!** Otherwise you will not get any positions through the channel.

```golang
package main

import (
	"context"
	"fmt"

	"github.com/rtunazzz/bfldb"
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
			fmt.Println("error has occured:", err)
			break
		}
	}
}
```

## Searching for user by their nickname

```golang
package main

import (
	"context"
	"fmt"

	"github.com/rtunazzz/bfldb"
)

func main() {
	res, err := ftl.SearchNickname(context.Background(), "TreeOfAlpha")
	if err != nil {
		panic(err)
	}

	// res.Data is an array here so it can include more than one result
	fmt.Printf("%+v\n", res.Data)
}
```

## Getting profile details for an user

```golang
package main

import (
	"context"
	"fmt"

	"github.com/rtunazzz/bfldb"
)

func main() {
	// You can find this UID (encryptedUid) in the end of a leaderboard profile URL. For example:
	// https://www.binance.com/en/futures-activity/leaderboard/user?encryptedUid=47E6D002EBB1173967A6561F72B9395C
	u := ftl.NewUser("47E6D002EBB1173967A6561F72B9395C")
	res, err := u.GetOtherLeaderboardBaseInfo(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", res.Data)
}
```
