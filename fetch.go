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
func doPost[T UserPositionData | UserBaseInfo](path string, data io.Reader) (ldbres LdbAPIRes[T], err error) {
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

	res, err := http.DefaultClient.Do(req)
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

// ***************************** /getOtherPosition *****************************

// UserPositionData represents data about user's positions.
type UserPositionData struct {
	OtherPositionRetList []rawPosition `json:"otherPositionRetList"`
	UpdateTimeStamp      int64         `json:"updateTimeStamp"`
	UpdateTime           []int         `json:"updateTime"`
}

// rawPosition represent details of an individual position returned.
type rawPosition struct {
	Symbol          string  `json:"symbol"`
	EntryPrice      float64 `json:"entryPrice"`
	MarkPrice       float64 `json:"markPrice"`
	Pnl             float64 `json:"pnl"`
	Roe             float64 `json:"roe"`
	Amount          float64 `json:"amount"`
	UpdateTimeStamp int64   `json:"updateTimeStamp"`
	UpdateTime      []int   `json:"updateTime"`
	Yellow          bool    `json:"yellow"`
	TradeBefore     bool    `json:"tradeBefore"`
	Leverage        int     `json:"leverage"`
}

// GetOtherPosition gets all currently open positions for a user.
func GetOtherPosition(uid string) (upr LdbAPIRes[UserPositionData], err error) {
	return doPost[UserPositionData]("/getOtherPosition", strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", uid)))
}

// ***************************** /getOtherLeaderboardBaseInfo *****************************

// UserBaseInfo represents user's data.
type UserBaseInfo struct {
	NickName               string      `json:"nickName"`
	UserPhotoURL           string      `json:"userPhotoUrl"`
	PositionShared         bool        `json:"positionShared"`
	DeliveryPositionShared bool        `json:"deliveryPositionShared"`
	FollowingCount         int         `json:"followingCount"`
	FollowerCount          int         `json:"followerCount"`
	TwitterURL             string      `json:"twitterUrl"`
	Introduction           string      `json:"introduction"`
	TwShared               bool        `json:"twShared"`
	IsTwTrader             bool        `json:"isTwTrader"`
	OpenID                 interface{} `json:"openId"`
}

// GetOtherLeaderboardBaseInfo gets information about an user.
func GetOtherLeaderboardBaseInfo(uid string) (upr LdbAPIRes[UserBaseInfo], err error) {
	return doPost[UserBaseInfo]("/getOtherLeaderboardBaseInfo", strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", uid)))
}
