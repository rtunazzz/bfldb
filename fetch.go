package bfldb

import (
	"context"
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
type LdbAPIRes[T UserPositionData | UserBaseInfo | []NicknameDetails] struct {
	Success       bool        `json:"success"`       // Whether or not the request was successful
	Code          string      `json:"code"`          // Error code, "000000" means success
	Message       string      `json:"message"`       // Error message
	Data          T           `json:"data"`          // Data
	MessageDetail interface{} `json:"messageDetail"` // ???
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

// GetOtherPosition gets all currently open positions for an user.
func GetOtherPosition(ctx context.Context, UUID string) (LdbAPIRes[UserPositionData], error) {
	var res LdbAPIRes[UserPositionData]
	return res, doPost(ctx, http.DefaultClient, "/getOtherPosition", strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", UUID)), &res)
}

// GetOtherPosition gets all currently open positions for an user.
func (u *User) GetOtherPosition(ctx context.Context) (LdbAPIRes[UserPositionData], error) {
	return GetOtherPosition(ctx, u.UID)
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

// GetOtherLeaderboardBaseInfo gets information for the uuid passed in.
func GetOtherLeaderboardBaseInfo(ctx context.Context, UUID string) (LdbAPIRes[UserBaseInfo], error) {
	var res LdbAPIRes[UserBaseInfo]
	return res, doPost(ctx, http.DefaultClient, "/getOtherLeaderboardBaseInfo", strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", UUID)), &res)
}

// GetOtherLeaderboardBaseInfo gets information about an user.
func (u *User) GetOtherLeaderboardBaseInfo(ctx context.Context) (LdbAPIRes[UserBaseInfo], error) {
	return GetOtherLeaderboardBaseInfo(ctx, u.UID)
}

// ************************************************** /searchNickname **************************************************

type NicknameDetails struct {
	EncryptedUID  string `json:"encryptedUid"`
	Nickname      string `json:"nickname"`
	FollowerCount int    `json:"followerCount"`
	UserPhotoURL  string `json:"userPhotoUrl"`
}

// SearchNickname searches for a nickname.
func SearchNickname(ctx context.Context, nickname string) (LdbAPIRes[[]NicknameDetails], error) {
	var res LdbAPIRes[[]NicknameDetails]
	return res, doPost(ctx, http.DefaultClient, "/searchNickname", strings.NewReader(fmt.Sprintf("{\"nickname\":\"%s\"}", nickname)), &res)
}

// ************************************************** Unexported **************************************************

// doPost POSTs the data passed in to the path on Binance's leaderboard API.
func doPost(ctx context.Context, c *http.Client, path string, data io.Reader, resPtr any) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		apiBase+path,
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
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
		return fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}

	// all of the endpoints are expected to return 200
	if res.StatusCode != http.StatusOK {
		return BadStatusError{
			Status:     res.Status,
			StatusCode: res.StatusCode,
			Body:       body,
		}
	}

	return json.Unmarshal(body, resPtr)
}
