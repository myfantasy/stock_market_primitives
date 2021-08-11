package smp

import (
	"encoding/json"
	"math"
	"time"

	"github.com/myfantasy/mfs"
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

// DefaultStrategyGenerator used for marshal and unmarhal to Json StrategyStorage
var DefaultStrategyGenerator = StrategyGeneratorCreate()

type StrategyCreator func() (s Strategy)
type StrategyLoader func(data json.RawMessage) (s Strategy, err *mft.Error)

type StrategyGenerator struct {
	Creators map[string]StrategyCreator
	Loaders  map[string]StrategyLoader

	Mx mfs.PMutex
}

func StrategyGeneratorCreate() (sg *StrategyGenerator) {
	return &StrategyGenerator{
		Creators: make(map[string]StrategyCreator),
		Loaders:  make(map[string]StrategyLoader),
	}
}

func (sg *StrategyGenerator) Add(typeStrategy string, creator StrategyCreator, loader StrategyLoader) {
	sg.Mx.Lock()
	defer sg.Mx.Unlock()

	sg.Creators[typeStrategy] = creator
	sg.Loaders[typeStrategy] = loader
}

func (sg *StrategyGenerator) Exists(typeStrategy string) (ok bool) {
	sg.Mx.RLock()
	defer sg.Mx.RUnlock()

	_, ok = sg.Creators[typeStrategy]
	return ok
}
func (sg *StrategyGenerator) Create(typeStrategy string) (s Strategy, err *mft.Error) {
	sg.Mx.RLock()
	defer sg.Mx.RUnlock()

	cr, ok := sg.Creators[typeStrategy]
	if !ok {
		return nil, GenerateError(25000000, typeStrategy)
	}

	return cr(), nil
}
func (sg *StrategyGenerator) Load(typeStrategy string, data json.RawMessage) (s Strategy, err *mft.Error) {
	sg.Mx.RLock()
	defer sg.Mx.RUnlock()

	ld, ok := sg.Loaders[typeStrategy]
	if !ok {
		return nil, GenerateError(25000000, typeStrategy)
	}

	s, err = ld(data)

	if err != nil {
		return nil, GenerateErrorE(25000001, err, typeStrategy)
	}

	return s, nil
}

type StrategyStorage struct {
	Starategy map[string]Strategy

	Mx mfs.PMutex
}

func StrategyStorageCreate() (ss *StrategyStorage) {
	return &StrategyStorage{
		Starategy: make(map[string]Strategy),
	}
}

func (ss *StrategyStorage) MarshalJSON() ([]byte, error) {
	res := make(map[string]JsonTypedContainer)

	for k, s := range ss.Starategy {
		res[k] = JsonTypedContainer{
			Type: s.Type(),
			Data: s.Marshal(),
		}
	}

	return json.Marshal(res)
}

func (ss *StrategyStorage) UnmarshalJSON(b []byte) (err error) {
	var cont map[string]JsonTypedContainer

	err = json.Unmarshal(b, &cont)
	if err != nil {
		return GenerateErrorE(25000020, err)
	}

	ss.Starategy = make(map[string]Strategy)

	for k, data := range cont {
		s, errM := DefaultStrategyGenerator.Load(data.Type, data.Data)
		if errM != nil {
			return GenerateErrorE(25000021, errM, data.Type)
		}
		ss.Starategy[k] = s
	}

	return nil
}

type StepParams interface {
	GetCandles(instrumentId string, ticker string, dateFrom time.Time, dateTo time.Time) Candles
	GetOrderBook(instrumentId string, ticker string) OrderBook

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
	TnstrumentId string
	Ticker       string
	Cnt          int
	BuyPrice     float64
	OrderId      string
	IsComplete   bool
	IsSell       bool
}
type Log struct {
	TnstrumentId string
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
	Type() string
	Marshal() json.RawMessage
	Step(p StepParams)

	Status() Status

	Logs() Logs
}

func Round0(price float64) float64 {
	return math.Round(price*1000000) / 1000000
}

type VirtualMarket struct {
	Candles   Candles
	OrderBook OrderBook
	Position  int
}

func (vm *VirtualMarket) GetCandles(instrumentId string, ticker string, dateFrom time.Time, dateTo time.Time) Candles {
	return vm.Candles.After(dateFrom).Before(dateTo).Before(vm.OrderBook.Time).Clone()
}
func (vm *VirtualMarket) GetOrderBook(instrumentId string, ticker string) OrderBook {
	return vm.OrderBook
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
