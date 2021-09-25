package strategies

import (
	"encoding/json"
	"strconv"

	"github.com/myfantasy/mft"
	smp "github.com/myfantasy/stock_market_primitives"
)

//go:generate mfjson winged_swing.go

// Крылатые качели (WingedSwing)

var (
	_ smp.Strategy = &WingedSwing{}
)

//mfjson:interface smp.strategies.winged_swing
type WingedSwing struct {
	InstrumentId string `json:"instrument_id"`
	Ticker       string `json:"ticker"`

	Volume         int     `json:"volume"`
	LevelPriceUp   float64 `json:"level_price_up"`
	LevelPriceDown float64 `json:"level_price_down"`
	LevelPriceBuy  float64 `json:"level_price_buy"`

	IsOnline bool `json:"is_online"`

	InMarket      int     `json:"in_market"`
	InMarketPrice float64 `json:"in_market_price"`
	IsBought      bool    `json:"is_bought"`

	OrderIdSell string `json:"order_id_sell"`
	OrderIdBuy  string `json:"order_id_buy"`

	InMarketWait      int     `json:"in_market_wait"`
	InMarketPriceWait float64 `json:"in_market_price_wait"`

	Profit    float64 `json:"profit"`
	Iteration int     `json:"iteration"`

	Labels map[string]string `json:"labels"`
}

func (s *WingedSwing) Status() smp.StartegyStatus {
	return smp.StartegyStatus{
		IsOnline: s.IsOnline,
	}
}
func (s *WingedSwing) String() string {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
func (s *WingedSwing) Command(cmd smp.Command, params map[string]string) (ok bool, err *mft.Error) {
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
				return false, smp.GenerateErrorE(500000411, er0, cmd, fS)
			}
			s.LevelPriceBuy = f
		} else {
			return false, smp.GenerateError(500000410, cmd)
		}
	}

	if cmd == SetLevelUp {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000413, er0, cmd, fS)
			}
			s.LevelPriceUp = f
		} else {
			return false, smp.GenerateError(500000412, cmd)
		}
	}

	if cmd == SetLevelDown {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000415, er0, cmd, fS)
			}
			s.LevelPriceDown = f
		} else {
			return false, smp.GenerateError(500000414, cmd)
		}
	}

	if cmd == SetVolume {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseInt(fS, 10, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000421, er0, cmd, fS)
			}
			s.Volume = int(f)
		} else {
			return false, smp.GenerateError(500000420, cmd)
		}
	}

	return false, smp.GenerateError(500000400, cmd)
}

func (s *WingedSwing) AllowCommands() map[smp.Command]string {
	return map[smp.Command]string{
		smp.ShowCommand:  "Отобразить",
		smp.StartCommand: "Старт",
		smp.StopCommand:  "Стоп",
		SetLevel:         "Установить уровень начала работы стратегии (с этого уровня происходит покупка) (параметр: уровень) пример: set_level 345.67",
		SetLevelDown:     "Установить уровень покупки (параметр: уровень) пример: set_level_down 372.25",
		SetLevelUp:       "Установить уровень продажи (параметр: уровень) пример: set_level_up 392.80",
		SetVolume:        "Установить объём (параметр: объём) пример set_vol 25",
	}
}

func (s *WingedSwing) Step(p smp.StepParams) (err *mft.Error) {

	if s.OrderIdBuy == "" && s.OrderIdSell == "" {
		// Заполняем если установили объёмы руками
		if s.IsBought && s.InMarket <= 0 {
			s.IsBought = false
			s.InMarketPrice = 0
		}

		if !s.IsBought && s.InMarket >= s.Volume {
			s.IsBought = true
			s.InMarketPrice = s.LevelPriceDown * float64(s.InMarket)
		}
	}

	if s.OrderIdBuy != "" {
		status, prices, err := p.StatusBuyOrder(s.InstrumentId, s.Ticker, s.OrderIdBuy)
		if err != nil {
			return smp.GenerateErrorE(500000502, err)
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
			s.OrderIdBuy = ""
			s.InMarketWait = 0
			s.InMarketPriceWait = 0

			if s.InMarket >= s.Volume {
				s.IsBought = true
			}

		} else {
			s.InMarketWait = cnt
			s.InMarketPriceWait = price
		}
	}

	if s.OrderIdSell != "" {
		status, prices, err := p.StatusSellOrder(s.InstrumentId, s.Ticker, s.OrderIdSell)
		if err != nil {
			return smp.GenerateErrorE(500000503, err)
		}

		cnt := 0
		price := 0.0
		for _, pr := range prices {
			cnt += pr.Count
			price += pr.Price * float64(pr.Count)
		}

		if status == smp.Complete || status == smp.Canceled {
			s.InMarket -= cnt
			s.InMarketPrice -= price
			s.OrderIdSell = ""
			s.InMarketWait = 0
			s.InMarketPriceWait = 0

			if s.InMarket <= 0 {
				s.Profit = -s.InMarketPrice
				s.InMarketPrice = 0
				s.IsBought = false
			}
		} else {
			s.InMarketWait = -cnt
			s.InMarketPriceWait = -price
		}
	}

	if !s.IsOnline {
		return nil
	}

	if s.OrderIdBuy != "" {
		return nil
	}
	if s.OrderIdSell != "" {
		return nil
	}

	ob, err := p.GetOrderBook(s.InstrumentId, s.Ticker)
	if err != nil {
		return smp.GenerateErrorE(500000500, err)
	}

	if ob.TradeStatus != smp.NormalTrading {
		return nil
	}

	if len(ob.Bids) < 1 || len(ob.Asks) < 1 {
		return nil
	}

	if s.IsBought {
		if s.InMarket > 0 {
			s.OrderIdSell, err = p.SellByPrice(s.InstrumentId, s.Ticker, s.InMarket, s.LevelPriceDown)
			if err != nil {
				s.OrderIdBuy = ""
				return smp.GenerateErrorE(500000504, err)
			}
		} else {
			return smp.GenerateError(500000505, s.InMarket)
		}
	} else {
		if s.Volume-s.InMarket > 0 {
			if ob.BuyPrice() <= s.LevelPriceBuy {
				if ob.BuyPrice() < s.LevelPriceDown {
					s.OrderIdBuy, err = p.BuyByPrice(s.InstrumentId, s.Ticker, s.Volume-s.InMarket, s.LevelPriceDown)
				} else if ob.SellPrice() < s.LevelPriceUp {
					s.OrderIdBuy, err = p.BuyByPrice(s.InstrumentId, s.Ticker, s.Volume-s.InMarket, s.LevelPriceDown)
				}
				if err != nil {
					s.OrderIdBuy = ""
					return smp.GenerateErrorE(500000501, err)
				}
			}
		} else {
			return smp.GenerateError(500000506, s.InMarket, s.Volume)
		}
	}
	return nil
}
