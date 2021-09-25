package smp

import (
	"time"
)

type TradingStatus string

const (
	BreakInTrading               TradingStatus = "break_in_trading"
	NormalTrading                TradingStatus = "normal_trading"
	NotAvailableForTrading       TradingStatus = "not_available_for_trading"
	ClosingAuction               TradingStatus = "closing_auction"
	ClosingPeriod                TradingStatus = "closing_period"
	DarkPoolAuction              TradingStatus = "dark_pool_auction"
	DiscreteAuction              TradingStatus = "discrete_auction"
	OpeningPeriod                TradingStatus = "opening_period"
	OpeningAuctionPeriod         TradingStatus = "opening_auction_period"
	TradingAtClosingAuctionPrice TradingStatus = "trading_at_closing_auction_price"
)

type RestPriceQuantity struct {
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type OrderBook struct {
	InstrumentId string `json:"instrument_id"`
	Ticker       string `json:"ticker"`

	Time  time.Time `json:"time"`
	Depth int       `json:"depth"`
	// under price
	Bids []RestPriceQuantity `json:"bids"`
	// over price
	Asks              []RestPriceQuantity `json:"asks"`
	TradeStatus       TradingStatus       `json:"trade_status"`
	MinPriceIncrement float64             `json:"min_price_increment"`
	LastPrice         float64             `json:"last_price,omitempty"`
	ClosePrice        float64             `json:"close_price,omitempty"`
	LimitUp           float64             `json:"limit_up,omitempty"`
	LimitDown         float64             `json:"limit_down,omitempty"`
}

func (ob *OrderBook) Price() float64 {
	return ob.LastPrice
}

func (ob *OrderBook) BuyPrice() float64 {
	if len(ob.Asks) > 0 {
		return ob.Asks[0].Price
	}
	return ob.LastPrice
}

func (ob *OrderBook) SellPrice() float64 {
	if len(ob.Bids) > 0 {
		return ob.Bids[0].Price
	}
	return ob.LastPrice
}

type InstrumentInfo struct {
	InstrumentId string `json:"instrument_id"`
	Ticker       string `json:"ticker"`

	LotSize int `json:"lot_size"`
}
