package ftl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// rawPosition represent details of an individual position returned from Binance's public leaderboard API.
type rawPosition struct {
	Symbol          string  `json:"symbol"`
	EntryPrice      float64 `json:"entryPrice"`
	MarkPrice       float64 `json:"markPrice"`
	Pnl             float64 `json:"pnl"`
	Roe             float64 `json:"roe"`
	Amount          float64 `json:"amount"`
	UpdateTimeStamp int64   `json:"updateTimeStamp"`
	Yellow          bool    `json:"yellow"`
	TradeBefore     bool    `json:"tradeBefore"`
	Leverage        int     `json:"leverage"`

	// UpdateTime      []int   `json:"updateTime"`
}

// userPosRes represents a response containing all position details for a particular user.
//
// This struct appears in the response of the getOtherPosition route on Binance's public leaderboard API.
type userPosRes struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		OtherPositionRetList []rawPosition `json:"otherPositionRetList"`
		UpdateTimeStamp      int64         `json:"updateTimeStamp"`
		// UpdateTime           []int             `json:"updateTime"`
	} `json:"data"`

	// MessageDetail interface{} `json:"messageDetail"`
}

// getOpenPositions gets all currently open positions for a user.
func getOpenPositions(uid string) (upr userPosRes, err error) {
	url := "https://www.binance.com/bapi/futures/v1/public/future/leaderboard/getOtherPosition"

	req, err := http.NewRequest("POST", url,
		strings.NewReader(fmt.Sprintf("{\"encryptedUid\":\"%s\",\"tradeType\":\"PERPETUAL\"}", uid)),
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

	return upr, json.Unmarshal(body, &upr)
}
