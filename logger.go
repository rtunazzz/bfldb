package bfldb

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "binance-ftl: ", log.Ldate|log.Ltime|log.Lshortfile)
}
