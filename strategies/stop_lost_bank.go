package strategies

import (
	"github.com/myfantasy/mfs"
	"github.com/myfantasy/mft"
	smp "github.com/myfantasy/stock_market_primitives"
)

type InstrumentInfo struct {
	Cnt          int     `json:"cnt"`
	SellCnt      int     `json:"sell_cnt"`
	SellPrice    float64 `json:"sell_price"`
	BuyCnt       int     `json:"buy_cnt"`
	BuyPrice     float64 `json:"buy_price"`
	RequestToBuy int     `json:"req_to_buy"`
}

type StopLostBank struct {
	Items map[string]*InstrumentInfo `json:"instrument_id"`
	mx    mfs.PMutex
}

func MakeKey(instrumentId string, ticker string) string {
	return instrumentId + "-" + ticker
}

func (slb *StopLostBank) AllowCount(instrumentId string, ticker string) int {
	key := MakeKey(instrumentId, ticker)
	slb.mx.Lock()
	defer slb.mx.Unlock()
	if slb.Items == nil {
		return 0
	}
	ii, ok := slb.Items[key]
	if !ok {
		return 0
	}

	return ii.Cnt
}

func (slb *StopLostBank) AllowRequestCount(instrumentId string, ticker string) int {
	key := MakeKey(instrumentId, ticker)
	slb.mx.Lock()
	defer slb.mx.Unlock()
	if slb.Items == nil {
		return 0
	}
	ii, ok := slb.Items[key]
	if !ok {
		return 0
	}

	return ii.RequestToBuy
}

func (slb *StopLostBank) RequestToBuyDO(p smp.StepParams, instrumentId string, ticker string, cnt int, price float64) (err *mft.Error) {
	key := MakeKey(instrumentId, ticker)
	slb.mx.Lock()
	defer slb.mx.Unlock()
	if slb.Items == nil {
		slb.Items = make(map[string]*InstrumentInfo)
	}

	ii, ok := slb.Items[key]
	if !ok {
		ii = &InstrumentInfo{}
		slb.Items[key] = ii
	}

	ii.RequestToBuy += ii.RequestToBuy

	return nil
}

func (slb *StopLostBank) Sell(p smp.StepParams, instrumentId string, ticker string, cnt int, price float64) (err *mft.Error) {
	key := MakeKey(instrumentId, ticker)
	slb.mx.Lock()
	defer slb.mx.Unlock()
	if slb.Items == nil {
		slb.Items = make(map[string]*InstrumentInfo)
	}

	ii, ok := slb.Items[key]
	if !ok {
		ii = &InstrumentInfo{}
		slb.Items[key] = ii
	}

	ii.Cnt += cnt
	ii.SellCnt += cnt
	ii.SellPrice = smp.Round(ii.SellPrice+float64(cnt)*price, 6)

	return nil
}
func (slb *StopLostBank) Buy(p smp.StepParams, instrumentId string, ticker string, cnt int, price float64) (success int, err *mft.Error) {
	key := MakeKey(instrumentId, ticker)
	slb.mx.Lock()
	defer slb.mx.Unlock()
	if slb.Items == nil {
		slb.Items = make(map[string]*InstrumentInfo)
	}

	ii, ok := slb.Items[key]
	if !ok {
		return 0, nil
	}

	if ii.Cnt < cnt {
		cnt = ii.Cnt
	}
	ii.Cnt -= cnt
	ii.BuyCnt += cnt
	ii.BuyPrice = smp.Round(ii.BuyPrice+float64(cnt)*price, 6)
	return cnt, nil
}
