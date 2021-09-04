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

type StepParams interface {
	GetCandles(instrumentId string, ticker string, dateFrom time.Time, dateTo time.Time) (cs Candles, err *mft.Error)
	GetOrderBook(instrumentId string, ticker string) (ob *OrderBook, err *mft.Error)

	BuyByMarket(instrumentId string, ticker string, cnt int) (orderId string, err *mft.Error)
	SellByMarket(instrumentId string, ticker string, cnt int) (orderId string, err *mft.Error)

	BuyByPrice(instrumentId string, ticker string, cnt int, price float64) (orderId string, err *mft.Error)
	SellByPrice(instrumentId string, ticker string, cnt int, price float64) (orderId string, err *mft.Error)

	CancelBuyOrder(instrumentId string, ticker string, orderId string) (ok bool, err *mft.Error)
	CancelSellOrder(instrumentId string, ticker string, orderId string) (ok bool, err *mft.Error)

	StatusBuyOrder(instrumentId string, ticker string, orderId string) (status StatusOrder, err *mft.Error)
	StatusSellOrder(instrumentId string, ticker string, orderId string) (status StatusOrder, err *mft.Error)
}

type Position struct {
	InstrumentId string
	Ticker       string
	Cnt          int
	BuyPrice     float64
	OrderId      string
	IsComplete   bool
	IsSell       bool
}
type Log struct {
	InstrumentId string
	Ticker       string
	Cnt          int
	Price        float64
	Operation    Operation
}

type Positions []Position
type Logs []Log

type Status interface {
	Profit() float64
	Positions() Positions
}

type Strategy interface {
	Step(p StepParams)

	Status() Status

	Logs() Logs
}

func Round0(price float64) float64 {
	return math.Round(price*1000000) / 1000000
}

type VirtualMarket struct {
	Candles   Candles
	OrderBook *OrderBook
	Position  int
}

func (vm *VirtualMarket) GetCandles(instrumentId string, ticker string, dateFrom time.Time, dateTo time.Time) (cs Candles, err *mft.Error) {
	return vm.Candles.After(dateFrom).Before(dateTo).Before(vm.OrderBook.Time).Clone(),
		nil
}
func (vm *VirtualMarket) GetOrderBook(instrumentId string, ticker string) (ob *OrderBook, err *mft.Error) {
	return vm.OrderBook, nil
}
func (vm *VirtualMarket) BuyByMarket(instrumentId string, ticker string, cnt int) (orderId string, err *mft.Error) {
	return "", nil
}
func (vm *VirtualMarket) SellByMarket(instrumentId string, ticker string, cnt int) (orderId string, err *mft.Error) {
	return "", nil
}
func (vm *VirtualMarket) BuyByPrice(instrumentId string, ticker string, cnt int, price float64) (orderId string, err *mft.Error) {
	return "", nil
}
func (vm *VirtualMarket) SellByPrice(instrumentId string, ticker string, cnt int, price float64) (orderId string, err *mft.Error) {
	return "", nil
}
func (vm *VirtualMarket) CancelBuyOrder(instrumentId string, ticker string, orderId string) (ok bool, err *mft.Error) {
	return false, nil
}
func (vm *VirtualMarket) CancelSellOrder(instrumentId string, ticker string, orderId string) (ok bool, err *mft.Error) {
	return false, nil
}
func (vm *VirtualMarket) StatusBuyOrder(instrumentId string, ticker string, orderId string) (status StatusOrder, err *mft.Error) {
	return Complete, nil
}
func (vm *VirtualMarket) StatusSellOrder(instrumentId string, ticker string, orderId string) (status StatusOrder, err *mft.Error) {
	return Complete, nil
}

func (vm *VirtualMarket) DoStep() bool {
	if vm.Position >= vm.Candles.Len() {
		return false
	}
	vm.Position++
	if vm.Position >= vm.Candles.Len() {
		return false
	}

	vm.OrderBook = vm.Candles[vm.Position].OrderBook()
	return true
}
