package bfldb

import (
	"fmt"
)

type BadStatusError struct {
	Status     string
	StatusCode int
	Body       []byte
}

func (e BadStatusError) Error() string {
	return fmt.Sprintf("%s", e.Status)
}

var _ error = (*BadStatusError)(nil)
