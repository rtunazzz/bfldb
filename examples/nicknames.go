package main

import (
	"context"
	"fmt"

	"github.com/rtunazzz/binance-bfldb"
)

func main() {
	res, err := bfldb.SearchNickname(context.Background(), "TreeOfAlpha")
	if err != nil {
		panic(err)
	}

	// res.Data is an array here so it can include more than one result
	fmt.Printf("%+v\n", res.Data)
}
