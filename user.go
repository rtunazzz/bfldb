package ftl

import (
	"log"
	"time"
)

// User represents one Binance leaderboard User
type User struct {
	UID  string // Encrypted User ID
	log  *log.Logger
	poss map[uint64]Position
	d    time.Duration
}

type UserOption func(*User)

// NewUser creates a new User with his encrypted UserID.
func NewUser(UID string, opts ...UserOption) User {
	u := User{
		UID:  UID,
		log:  log.Default(),
		poss: make(map[uint64]Position),
		d:    time.Second * 5,
	}

	for _, opt := range opts {
		opt(&u)
	}

	return u
}

func WithLogger(l *log.Logger) UserOption {
	return func(u *User) {
		u.log = l
	}
}

func WithCustomRefresh(d time.Duration) UserOption {
	return func(u *User) {
		u.d = d
	}
}
