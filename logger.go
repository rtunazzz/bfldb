package bfldb

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "bfldb: ", log.Ldate|log.Ltime|log.Lshortfile)
}
