package smp

import (
	"math"
	"time"

	"github.com/myfantasy/mft"
)

type StatusOrder string

const (
	Complete StatusOrder = "complete"
	Canceled StatusOrder = "canceled"
	Unknown  StatusOrder = "unknown"
	Wait     StatusOrder = "wait"
)

type Operation string

const (
	Buy             Operation = "buy"
	Sell            Operation = "sell"
	Tax             Operation = "tax"
	BuyOrderCreate  Operation = "buy_order_create"
	SellOrderCreate Operation = "sell_order_create"
)

type LotPrices struct {
	Count int
	Price float64
}

type StepParams interface {
	GetCandles(instrumentId string, ticker string, dateFrom time.Time, dateTo time.Time) (cs Candles, err *mft.Error)
	GetOrderBook(instrumentId string, ticker string) (ob *OrderBook, err *mft.Error)
	GetInstrumentInfo(instrumentId string, ticker string) (instrumentInfo *InstrumentInfo, err *mft.Error)

	BuyByMarket(instrumentId string, ticker string, cnt int) (orderId string, err *mft.Error)
	SellByMarket(instrumentId string, ticker string, cnt int) (orderId string, err *mft.Error)

	BuyByPrice(instrumentId string, ticker string, cnt int, price float64) (orderId string, err *mft.Error)
	SellByPrice(instrumentId string, ticker string, cnt int, price float64) (orderId string, err *mft.Error)

	CancelBuyOrder(instrumentId string, ticker string, orderId string) (ok bool, err *mft.Error)
	CancelSellOrder(instrumentId string, ticker string, orderId string) (ok bool, err *mft.Error)

	StatusBuyOrder(instrumentId string, ticker string, orderId string) (status StatusOrder, prices []LotPrices, err *mft.Error)
	StatusSellOrder(instrumentId string, ticker string, orderId string) (status StatusOrder, prices []LotPrices, err *mft.Error)
}

type StartegyStatus struct {
	IsOnline bool `json:"is_online"`
}

type Command string

const (
	StartCommand Command = "start"
	StopCommand  Command = "stop"
	ShowCommand  Command = "show"
)

type Strategy interface {
	Step(p StepParams) (err *mft.Error)
	Status() StartegyStatus
	String() string
	Command(cmd Command, params map[string]string) (ok bool, err *mft.Error)
	AllowCommands() map[Command]string
	Description() string
}

func Round(price float64, point int) float64 {
	if point == 0 {
		return math.Round(price)
	}
	mult := 1.0
	if point > 0 {
		for i := 0; i < point; i++ {
			mult = mult * 10
		}
		return math.Round(price*mult) / mult
	}
	if point < 0 {
		for i := 0; i < -point; i++ {
			mult = mult + 10
		}
		return math.Round(price/mult) * mult
	}
	return math.Round(price)
}
