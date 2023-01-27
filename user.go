package bfldb

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// User represents one Binance leaderboard User
type User struct {
	UID     string              // Encrypted User ID
	id      string              // identified used in logging
	log     *log.Logger         // Logger
	pHashes map[string]Position // map of positions user is currently in
	c       *http.Client        // http client
	isFirst bool                // indicating first fetch

	mtx     sync.Mutex        // Synchronization for d and headers
	d       time.Duration     // duration between requests updating current positions
	headers map[string]string // headers
}

type UserOption func(*User)

// NewUser creates a new User with his encrypted UserID.
func NewUser(UID string, opts ...UserOption) *User {
	u := User{
		id:      UID,
		UID:     UID,
		log:     logger,
		pHashes: make(map[string]Position),
		d:       time.Second * 5,
		c:       http.DefaultClient,
		isFirst: true,
		headers: defaultHeaders,
	}

	// disable logging by default
	u.log.SetOutput(io.Discard)

	for _, opt := range opts {
		opt(&u)
	}

	return &u
}

// SetDelay sets the delay between requests updating user's current positions.
func (u *User) SetDelay(d time.Duration) {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	u.d = d
}

// GetDelay returns the delay between requests updating user's current positions
func (u *User) GetDelay() time.Duration {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	return u.d
}

// WithID sets user's logging id.
func WithID(id string) UserOption {
	return func(u *User) {
		u.id = id
	}
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

// WithHTTPClient sets user's HTTP Client.
func WithHTTPClient(c *http.Client) UserOption {
	return func(u *User) {
		u.c = c
	}
}

// WithHeaders sets headers the client uses for every request.
func WithHeaders(h map[string]string) UserOption {
	return func(u *User) {
		u.headers = h
	}
}
