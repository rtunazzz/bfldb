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
	apiBase   = "https://www.binance.com/bapi/futures"
	apiBaseV1 = apiBase + "/v1/public/future/leaderboard"
	apiBaseV2 = apiBase + "/v2/public/future/leaderboard"
)

var (
	defaultHeaders = map[string]string{
		"authority":       "www.binance.com",
		"accept":          "*/*",
		"accept-language": "en-US,en;q=0.8",
		"cache-control":   "no-cache",
		"clienttype":      "web",
		"content-type":    "application/json",
		"lang":            "en",
		"origin":          "https://www.binance.com",
		"pragma":          "no-cache",
		"sec-fetch-dest":  "empty",
		"sec-fetch-mode":  "cors",
		"sec-fetch-site":  "same-origin",
		"sec-gpc":         "1",
		"user-agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
	}
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
	return res, doPost(ctx, http.DefaultClient, apiBaseV1, "/getOtherPosition", defaultHeaders, strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", UUID)), &res)
}

// GetOtherPosition gets all currently open positions for an user.
func (u User) GetOtherPosition(ctx context.Context) (LdbAPIRes[UserPositionData], error) {
	var res LdbAPIRes[UserPositionData]
	return res, doPost(ctx, u.c, apiBaseV1, "/getOtherPosition", u.headers, strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", u.UID)), &res)
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
	return res, doPost(ctx, http.DefaultClient, apiBaseV2, "/getOtherLeaderboardBaseInfo", defaultHeaders, strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", UUID)), &res)
}

// GetOtherLeaderboardBaseInfo gets information about an user.
func (u User) GetOtherLeaderboardBaseInfo(ctx context.Context) (LdbAPIRes[UserBaseInfo], error) {
	var res LdbAPIRes[UserBaseInfo]
	return res, doPost(ctx, u.c, apiBaseV2, "/getOtherLeaderboardBaseInfo", u.headers, strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", u.UID)), &res)
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
	return res, doPost(ctx, http.DefaultClient, apiBaseV1, "/searchNickname", defaultHeaders, strings.NewReader(fmt.Sprintf("{\"nickname\":\"%s\"}", nickname)), &res)
}

// ************************************************** Unexported **************************************************

// doPost POSTs the data passed in to the path on Binance's leaderboard API.
func doPost(ctx context.Context, c *http.Client, endpoint, path string, headers map[string]string, data io.Reader, resPtr any) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		endpoint+path,
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

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
