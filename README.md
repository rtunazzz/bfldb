# bfldb [![GoDoc](https://godoc.org/github.com/rtunazzz/bfldb?status.svg)](https://godoc.org/github.com/rtunazzz/bfldb)
> BFLDB is a wrapper around **B**inance's **F**utures **L**ea**d**er**b**oard API, in Golang.

<p align="center"><img width=100% src="https://user-images.githubusercontent.com/38296319/200170420-0644f467-49ff-4ecd-8811-1bc939f84fea.png"></p>

This library provides a convenient way to access Binance's futures leaderboard API, which allows you to subscribe to positions opened by leading futures traders (that are sharing their positions publicly via Binance) and query other leaderboard data. Keep in mind that this is not a publicly documented API.

## Installation
```bash
go get -u github.com/rtunazzz/bfldb
```

## Example usage

<details>
<summary>Subscribe to a trader's positions</summary>

```golang
package main

import (
	"context"
	"fmt"

	"github.com/rtunazzz/bfldb"
)

func main() {
	// You can find this UID (encryptedUid) in the end of a leaderboard profile URL. 
	// For example: https://www.binance.com/en/futures-activity/leaderboard/user?encryptedUid=47E6D002EBB1173967A6561F72B9395C
	u := bfldb.NewUser("47E6D002EBB1173967A6561F72B9395C")
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
</details>

For more examples, check out the [`examples/`](./examples) directory.
