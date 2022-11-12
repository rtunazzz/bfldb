package main

import (
	"context"
	"fmt"

	"github.com/rtunazzz/binance-ftl"
)

func main() {
	res, err := ftl.SearchNickname(context.Background(), "TreeOfAlpha")
	if err != nil {
		panic(err)
	}

	// res.Data is an array here so it can include more than one result
	fmt.Printf("%+v\n", res.Data)
}
