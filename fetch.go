package ftl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	apiBase = "https://www.binance.com/bapi/futures/v1/public/future/leaderboard" // Base endpoint
)

// LdbAPIRes represents a response from Binance's Futures LDB API.
type LdbAPIRes[T UserPositionData | UserBaseInfo] struct {
	Success       bool        `json:"success"`       // Whether or not the request was successful
	Code          string      `json:"code"`          // Error code, "000000" means success
	Message       string      `json:"message"`       // Error message
	Data          T           `json:"data"`          // Data
	MessageDetail interface{} `json:"messageDetail"` // ???
}

// doPost POSTs the data passed in to the path on Binance's leaderboard API.
func doPost[T UserPositionData | UserBaseInfo](c *http.Client, path string, data io.Reader) (ldbres LdbAPIRes[T], err error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	req, err := http.NewRequest(
		"POST",
		apiBase+path,
		data,
	)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		return
	}

	req.Header.Add("authority", "www.binance.com")
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.8")
	req.Header.Add("cache-control", "no-cache")
	req.Header.Add("clienttype", "web")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("lang", "en")
	req.Header.Add("origin", "https://www.binance.com")
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-gpc", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36")

	res, err := c.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to do request: %w", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("failed to read request body: %w", err)
		return
	}

	return ldbres, json.Unmarshal(body, &ldbres)
}

// ************************************************** /getOtherPosition **************************************************

// UserPositionData represents data about user's positions.
type UserPositionData struct {
	OtherPositionRetList []rawPosition `json:"otherPositionRetList"` // List of positions
	UpdateTimeStamp      int64         `json:"updateTimeStamp"`      // Timestamp
	UpdateTime           []int         `json:"updateTime"`           // Time array in the format of [YEAR, MONTH, DAY, HOUR, MINUTE, SECOND, ... ]
}

// rawPosition represent details of an individual position returned.
type rawPosition struct {
	Symbol          string  `json:"symbol"`          // Position symbol
	EntryPrice      float64 `json:"entryPrice"`      // Entry price
	MarkPrice       float64 `json:"markPrice"`       // Mark Price
	Pnl             float64 `json:"pnl"`             // PNL
	Roe             float64 `json:"roe"`             // ROE
	Amount          float64 `json:"amount"`          // Position size
	UpdateTimeStamp int64   `json:"updateTimeStamp"` // Timestamp
	UpdateTime      []int   `json:"updateTime"`      // Time array in the format of [YEAR, MONTH, DAY, HOUR, MINUTE, SECOND, ... ]
	Yellow          bool    `json:"yellow"`          // ???
	TradeBefore     bool    `json:"tradeBefore"`     // ???
	Leverage        int     `json:"leverage"`        // leverage used
}

// GetOtherPosition gets all currently open positions for a user.
func (u *User) GetOtherPosition() (upr LdbAPIRes[UserPositionData], err error) {
	return doPost[UserPositionData](u.c, "/getOtherPosition", strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", u.UID)))
}

// ************************************************** /getOtherLeaderboardBaseInfo **************************************************

// UserBaseInfo represents user's data.
type UserBaseInfo struct {
	NickName               string      `json:"nickName"`               // Nickname
	UserPhotoURL           string      `json:"userPhotoUrl"`           // Photo URL
	PositionShared         bool        `json:"positionShared"`         // true if user is sharing their positions, false otherwise
	DeliveryPositionShared bool        `json:"deliveryPositionShared"` // ???
	FollowingCount         int         `json:"followingCount"`         // How many people user follows
	FollowerCount          int         `json:"followerCount"`          // How many people follow user
	TwitterURL             string      `json:"twitterUrl"`             // Twitter URL
	Introduction           string      `json:"introduction"`           // Introduction (profile description)
	TwShared               bool        `json:"twShared"`               // ???
	IsTwTrader             bool        `json:"isTwTrader"`             // ???
	OpenID                 interface{} `json:"openId"`                 // ???
}

// GetOtherLeaderboardBaseInfo gets information about an user.
func (u *User) GetOtherLeaderboardBaseInfo() (upr LdbAPIRes[UserBaseInfo], err error) {
	return doPost[UserBaseInfo](u.c, "/getOtherLeaderboardBaseInfo", strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", u.UID)))
}
