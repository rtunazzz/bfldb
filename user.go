package ftl

import (
	"io"
	"log"
	"os"
	"time"
)

// User represents one Binance leaderboard User
type User struct {
	UID  string              // Encrypted User ID
	log  *log.Logger         // Logger
	poss map[uint64]Position // hashmap of positions user is currently in
	d    time.Duration       // duration between requests updating current positions
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

	u.log.SetOutput(io.Discard)

	for _, opt := range opts {
		opt(&u)
	}

	return u
}

// WithCustomLogger writes all user logs using the logger provided.
func WithCustomLogger(l *log.Logger) UserOption {
	return func(u *User) {
		u.log = l
	}
}

// WithLogging writes all user logs to STDOUT.
func WithLogging() UserOption {
	return func(u *User) {
		u.log.SetOutput(os.Stdout)
	}
}

// WithCustomRefresh sets the duration between requests updating user's current positions.
func WithCustomRefresh(d time.Duration) UserOption {
	return func(u *User) {
		u.d = d
	}
}
