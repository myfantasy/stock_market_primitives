package market

import (
	"strconv"
	"time"

	"github.com/myfantasy/mft"
	smp "github.com/myfantasy/stock_market_primitives"
)

var (
	_ smp.StepParams = &StepParamsDummy{}
)

type Action struct {
	Buy   bool
	Sell  bool
	Time  time.Time
	Cnt   int
	Price float64
}

type StepParamsDummy struct {
	Candles        smp.Candles
	OrderBook      *smp.OrderBook
	OrderBookNext  *smp.OrderBook
	InstrumentInfo *smp.InstrumentInfo
	Position       int

	Actions []Action

	nextId      int
	waitActions map[string]Action
}

func (sp *StepParamsDummy) GetCandles(instrumentId string, ticker string, dateFrom time.Time, dateTo time.Time) (cs smp.Candles, err *mft.Error) {
	return sp.Candles[0:sp.Position].After(dateFrom).Before(dateTo).Before(sp.OrderBook.Time).Clone(),
		nil
}
func (sp *StepParamsDummy) GetOrderBook(instrumentId string, ticker string) (ob *smp.OrderBook, err *mft.Error) {
	return sp.OrderBook, nil
}
func (sp *StepParamsDummy) GetInstrumentInfo(instrumentId string, ticker string) (instrumentInfo *smp.InstrumentInfo, err *mft.Error) {
	return sp.InstrumentInfo, nil
}

func (sp *StepParamsDummy) BuyByMarket(instrumentId string, ticker string, cnt int) (orderId string, err *mft.Error) {
	sp.Actions = append(sp.Actions, Action{
		Buy:   true,
		Cnt:   cnt,
		Price: sp.OrderBookNext.BuyPrice(),
		Time:  sp.OrderBookNext.Time,
	})
	return "dummy_order_id", nil
}
func (sp *StepParamsDummy) SellByMarket(instrumentId string, ticker string, cnt int) (orderId string, err *mft.Error) {
	sp.Actions = append(sp.Actions, Action{
		Sell:  true,
		Cnt:   cnt,
		Price: sp.OrderBookNext.SellPrice(),
		Time:  sp.OrderBookNext.Time,
	})
	return "dummy_order_id", nil
}

func (sp *StepParamsDummy) BuyByPrice(instrumentId string, ticker string, cnt int, price float64) (orderId string, err *mft.Error) {
	sp.nextId++
	orderId = strconv.Itoa(sp.nextId)
	sp.waitActions[orderId] = Action{
		Buy:   true,
		Cnt:   cnt,
		Price: price,
		Time:  sp.OrderBook.Time,
	}
	return orderId, nil
}
func (sp *StepParamsDummy) SellByPrice(instrumentId string, ticker string, cnt int, price float64) (orderId string, err *mft.Error) {
	sp.nextId++
	orderId = strconv.Itoa(sp.nextId)
	sp.waitActions[orderId] = Action{
		Sell:  true,
		Cnt:   cnt,
		Price: price,
		Time:  sp.OrderBook.Time,
	}
	return orderId, nil
}

func (sp *StepParamsDummy) CancelBuyOrder(instrumentId string, ticker string, orderId string) (ok bool, err *mft.Error) {
	_, ok = sp.waitActions[orderId]
	if ok {
		delete(sp.waitActions, orderId)
		return true, nil
	}
	return false, mft.ErrorS("Not found")
}
func (sp *StepParamsDummy) CancelSellOrder(instrumentId string, ticker string, orderId string) (ok bool, err *mft.Error) {
	_, ok = sp.waitActions[orderId]
	if ok {
		delete(sp.waitActions, orderId)
		return true, nil
	}
	return false, mft.ErrorS("Not found")
}

func (sp *StepParamsDummy) StatusBuyOrder(instrumentId string, ticker string, orderId string) (status smp.StatusOrder, prices []smp.LotPrices, err *mft.Error) {
	if orderId == "dummy_order_id" {
		return smp.Complete, []smp.LotPrices{{
			Count: 1,
			Price: sp.OrderBook.Price(),
		},
		}, nil
	}
	a, ok := sp.waitActions[orderId]
	if ok {
		if a.Time.Day() > sp.OrderBook.Time.Day() &&
			(sp.OrderBook.Time.Hour() > 5 || a.Time.Hour() <= 5) ||
			a.Time.Day() == sp.OrderBook.Time.Day() &&
				(a.Time.Hour() < 5 && sp.OrderBook.Time.Hour() >= 5) {
			delete(sp.waitActions, orderId)
			return smp.Canceled, []smp.LotPrices{}, nil
		}
		if a.Buy {
			if sp.OrderBook.LimitUp >= a.Price {
				delete(sp.waitActions, orderId)
				a.Time = sp.OrderBook.Time
				sp.Actions = append(sp.Actions, a)
				return smp.Complete, []smp.LotPrices{{
					Count: a.Cnt,
					Price: a.Price,
				},
				}, nil
			}
			return smp.Wait, make([]smp.LotPrices, 0), nil
		} else {
			if sp.OrderBook.LimitDown <= a.Price {
				delete(sp.waitActions, orderId)
				a.Time = sp.OrderBook.Time
				sp.Actions = append(sp.Actions, a)
				return smp.Complete, []smp.LotPrices{{
					Count: a.Cnt,
					Price: a.Price,
				},
				}, nil
			}
			return smp.Wait, make([]smp.LotPrices, 0), nil
		}
	}
	return smp.Canceled, make([]smp.LotPrices, 0), nil
}
func (sp *StepParamsDummy) StatusSellOrder(instrumentId string, ticker string, orderId string) (status smp.StatusOrder, prices []smp.LotPrices, err *mft.Error) {
	return sp.StatusBuyOrder(instrumentId, ticker, orderId)
}

func (sp *StepParamsDummy) DoStep() bool {
	if sp.Position >= sp.Candles.Len()-1 {
		return false
	}
	sp.Position++
	if sp.Position >= sp.Candles.Len()-1 {
		return false
	}

	sp.OrderBook = sp.Candles[sp.Position].OrderBook()
	sp.OrderBookNext = sp.Candles[sp.Position+1].OrderBook()
	return true
}
