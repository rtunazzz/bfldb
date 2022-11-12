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
	res, err := u.GetOtherLeaderboardBaseInfo(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", res.Data)
}
