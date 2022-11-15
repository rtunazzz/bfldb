package bfldb

// TradeDirection can be either LONG / SHORT
type TradeDirection int

const (
	Short TradeDirection = iota + 1
	Long
)

func (pd TradeDirection) String() string {
	if pd == Short {
		return "SHORT"
	}
	return "LONG"
}

// Position represents an order to be used for placing a trade.
type Order struct {
	Direction  TradeDirection // Direction (e.g. LONG / SHORT)
	Ticker     string         // Ticker of the position (e.g. BTCUSDT)
	Amount     float64        // Amount
	ReduceOnly bool           // Whether or not the order is reduce only
}
