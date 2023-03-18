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
	UID string // Encrypted User ID

	mtx     sync.RWMutex      // Synchronization for delay, apiBase and headers
	apiBase string            // API base used for requests
	delay   time.Duration     // duration between requests updating current positions
	headers map[string]string // headers

	positions  map[string]Position // map of positions user is currently in
	client     *http.Client        // http client
	log        *log.Logger         // Logger
	firstFetch bool                // indicating first fetch
}

type UserOption func(*User)

// NewUser creates a new User with his encrypted UserID.
func NewUser(UID string, opts ...UserOption) *User {
	u := User{
		UID:        UID,
		log:        logger,
		positions:  make(map[string]Position),
		delay:      time.Second * 5,
		client:     http.DefaultClient,
		firstFetch: true,
		headers:    defaultHeaders,
		apiBase:    defaultApiBase,
	}

	// disable logging by default
	u.log.SetOutput(io.Discard)

	for _, opt := range opts {
		opt(&u)
	}

	return &u
}

// SetAPIBase sets the API base used for requests.
func (u *User) SetAPIBase(s string) {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	u.apiBase = s
}

// APIBase returns the API base used for requests.
func (u *User) APIBase() string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()

	return u.apiBase
}

// SetDelay sets the delay between requests updating user's current positions.
func (u *User) SetDelay(d time.Duration) {
	u.mtx.Lock()
	defer u.mtx.Unlock()

	u.delay = d
}

// Delay returns the delay between requests updating user's current positions
func (u *User) Delay() time.Duration {
	u.mtx.RLock()
	defer u.mtx.RUnlock()

	return u.delay
}

// SetHeaders sets headers the client uses for every request.
func (u *User) SetHeaders(h map[string]string) {
	headers := make(map[string]string, len(h))
	// copy them so it doesn't matter if the input is modified by caller later
	for k, v := range h {
		headers[k] = v
	}

	u.mtx.Lock()
	defer u.mtx.Unlock()

	u.headers = h
}

// Headers returns headers the client uses for every request.
func (u *User) Headers() map[string]string {
	u.mtx.RLock()
	defer u.mtx.RUnlock()

	headers := make(map[string]string, len(u.headers))

	// copy them so it doesn't matter if they are modified by caller later
	for k, v := range u.headers {
		headers[k] = v
	}

	return headers
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
		u.delay = d
	}
}

// WithHTTPClient sets user's HTTP Client.
func WithHTTPClient(c *http.Client) UserOption {
	return func(u *User) {
		u.client = c
	}
}

// WithHeaders sets headers the client uses for every request.
func WithHeaders(h map[string]string) UserOption {
	return func(u *User) {
		u.headers = h
	}
}

// WithTestnet uses the testnet API
func WithTestnet() UserOption {
	return func(u *User) {
		u.SetAPIBase("https://testnet.binancefuture.com/bapi/future")
	}
}
