package strategies

import (
	"encoding/json"
	"strconv"

	"github.com/myfantasy/mft"
	smp "github.com/myfantasy/stock_market_primitives"
)

var (
	_ smp.Strategy = &TakeProfitSell{}
)

//go:generate mfjson take_profit_sell.go

//mfjson:interface smp.strategies.take_profit_sell
type TakeProfitSell struct {
	InstrumentId string `json:"instrument_id"`
	Ticker       string `json:"ticker"`

	Volume       int     `json:"volume"`
	LevelPrice   float64 `json:"level_price"`
	StayInMarket bool    `json:"stay_in_market"`

	IsOnline      bool    `json:"is_online"`
	InMarket      int     `json:"in_market"`
	InMarketPrice float64 `json:"in_market_price"`
	OrderId       string  `json:"order_id"`

	InMarketWait      int     `json:"in_market_wait"`
	InMarketPriceWait float64 `json:"in_market_price_wait"`
}

func (s *TakeProfitSell) Status() smp.StartegyStatus {
	return smp.StartegyStatus{
		IsOnline: s.IsOnline,
	}
}
func (s *TakeProfitSell) String() string {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
func (s *TakeProfitSell) Command(cmd smp.Command, params map[string]string) (ok bool, err *mft.Error) {
	if cmd == smp.ShowCommand {
		return true, nil
	}

	if cmd == smp.StartCommand {
		s.IsOnline = true
		return true, nil
	}

	if cmd == smp.StopCommand {
		s.IsOnline = false
		return true, nil
	}

	if cmd == SetLevel {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000211, er0, cmd, fS)
			}
			s.LevelPrice = f
		} else {
			return false, smp.GenerateError(500000210, cmd)
		}
	}

	if cmd == SetVolume {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseInt(fS, 10, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000221, er0, cmd, fS)
			}
			s.Volume = int(f)
		} else {
			return false, smp.GenerateError(500000220, cmd)
		}
	}
	if cmd == SetStayInMarket {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseBool(fS)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000231, er0, cmd, fS)
			}
			s.StayInMarket = f
		} else {
			return false, smp.GenerateError(500000230, cmd)
		}
	}

	return false, smp.GenerateError(500000200, cmd)
}

func (s *TakeProfitSell) AllowCommands() map[smp.Command]string {
	return map[smp.Command]string{
		smp.ShowCommand:  "Отобразить",
		smp.StartCommand: "Старт",
		smp.StopCommand:  "Стоп",
		SetLevel:         "Установить уровень (параметр: уровень) пример: set_level 345.67",
		SetVolume:        "Установить объём (параметр: объём) пример set_vol 25",
		SetStayInMarket:  "Оставаться в рынке [выставлять заявку не дожидаясь приближения цены] (параметр: true/false) пример stay_in_market true",
	}
}

func (s *TakeProfitSell) Step(p smp.StepParams) (err *mft.Error) {

	if s.OrderId != "" {
		status, prices, err := p.StatusSellOrder(s.InstrumentId, s.Ticker, s.OrderId)
		if err != nil {
			return smp.GenerateErrorE(500000302, err)
		}

		cnt := 0
		price := 0.0
		for _, pr := range prices {
			cnt += pr.Count
			price += pr.Price * float64(pr.Count)
		}

		if status == smp.Complete || status == smp.Canceled {
			s.InMarket += cnt
			s.InMarketPrice += price
			s.OrderId = ""
			s.InMarketWait = 0
			s.InMarketPriceWait = 0
		} else {
			s.InMarketWait = cnt
			s.InMarketPriceWait = price
		}
	}

	if !s.IsOnline {
		return nil
	}

	if s.OrderId != "" {
		return nil
	}

	ob, err := p.GetOrderBook(s.InstrumentId, s.Ticker)
	if err != nil {
		return smp.GenerateErrorE(500000300, err)
	}

	if ob.TradeStatus != smp.NormalTrading {
		return nil
	}

	if len(ob.Bids) < 1 || len(ob.Asks) < 1 {
		return nil
	}

	if s.InMarket < s.Volume && s.OrderId != "" {
		if s.StayInMarket ||
			ob.BuyPrice() >= s.LevelPrice {
			s.OrderId, err = p.SellByPrice(s.InstrumentId, s.Ticker, s.Volume-s.InMarket, s.LevelPrice)
			if err != nil {
				s.OrderId = ""
				return smp.GenerateErrorE(500000301, err)
			}
		}
	}
	return nil
}
