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

func (s *TakeProfitSell) Type() string {
	return "take_profit_sell"
}

func (s *TakeProfitSell) String() string {
	return "take_profit_sell"
}

func (s *TakeProfitSell) Status() smp.StartegyStatus {
	return smp.StartegyStatus{
		IsOnline: s.IsOnline,
	}
}
func (s *TakeProfitSell) Json() string {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
func (s *TakeProfitSell) Command(cmd smp.Command, params map[string]string) (res smp.CommandResult, ok bool, err *mft.Error) {
	if cmd == smp.ShowCommand {
		res.Message = s.String()
		return res, true, nil
	}

	if cmd == smp.StartCommand {
		s.IsOnline = true
		return res, true, nil
	}

	if cmd == smp.StopCommand {
		s.IsOnline = false
		return res, true, nil
	}

	if cmd == SetLevel {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000211, er0, cmd, fS)
			}
			s.LevelPrice = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000210, cmd)
		}
	}

	if cmd == SetVolume {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseInt(fS, 10, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000221, er0, cmd, fS)
			}
			s.Volume = int(f)
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000220, cmd)
		}
	}
	if cmd == SetStayInMarket {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseBool(fS)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000231, er0, cmd, fS)
			}
			s.StayInMarket = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000230, cmd)
		}
	}

	return res, false, smp.GenerateError(500000200, cmd)
}

func (s *TakeProfitSell) AllowCommands() map[smp.Command]smp.CommandInfo {
	return map[smp.Command]smp.CommandInfo{
		smp.ShowCommand:  {0, "Отобразить", "", ""},
		smp.StartCommand: {1, "Старт", "", ""},
		smp.StopCommand:  {2, "Стоп", "", ""},
		SetLevel:         {3, "Установить уровень", "параметр: уровень", "set_level 345.67"},
		SetVolume:        {4, "Установить объём", "параметр: объём", "set_vol 25"},
		SetStayInMarket:  {5, "Оставаться в рынке [выставлять заявку не дожидаясь приближения цены]", "параметр: true/false", "stay_in_market true"},
	}
}

func (s *TakeProfitSell) Description() string {
	return `take_profit_sell - стратегия, продажа профита
	При достижении цены указанного уровня (цена продажи на рынке) выставляется заявка на продажу по указанной цене
	При цене ниже указанной не происходит ничего`
}

func (s *TakeProfitSell) Step(p smp.StepParams) (meta smp.MetaForStep, err *mft.Error) {
	if s.Volume == 0 {
		return meta, nil
	}

	if s.OrderId != "" {
		status, prices, err := p.StatusSellOrder(s.InstrumentId, s.Ticker, s.OrderId, nil)
		if err != nil {
			return meta, smp.GenerateErrorE(500000302, err)
		}

		cnt := 0
		price := 0.0
		for _, pr := range prices {
			cnt += pr.Count
			price += pr.Price * float64(pr.Count)
		}

		meta.HasChanges = true
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
		return meta, nil
	}

	if s.OrderId != "" {
		return meta, nil
	}

	ob, err := p.GetOrderBook(s.InstrumentId, s.Ticker)
	if err != nil {
		return meta, smp.GenerateErrorE(500000300, err)
	}

	if ob.TradeStatus != smp.NormalTrading {
		return meta, nil
	}

	if len(ob.Bids) < 1 || len(ob.Asks) < 1 {
		return meta, nil
	}

	if s.InMarket < s.Volume && s.OrderId != "" {
		if s.StayInMarket ||
			ob.BuyPrice() >= s.LevelPrice {
			meta.HasChanges = true
			s.OrderId, err = p.SellByPrice(s.InstrumentId, s.Ticker, s.Volume-s.InMarket, s.LevelPrice, nil)
			if err != nil {
				s.OrderId = ""
				return meta, smp.GenerateErrorE(500000301, err)
			}
		}
	}
	return meta, nil
}
