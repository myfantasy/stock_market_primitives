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
	Name string `json:"name"`

	InstrumentId string `json:"instrument_id"`
	Ticker       string `json:"ticker"`

	Volume                            int     `json:"volume"`
	LevelPriceUp                      float64 `json:"level_price_up"`
	LevelPriceDown                    float64 `json:"level_price_down"`
	LevelPriceOnTheMarketDown         float64 `json:"level_price_on_the_market_down"`
	LevelPriceOnTheMarketDownByMarket float64 `json:"level_price_on_the_market_down_by_market"`
	LevelPriceOnTheMarketUp           float64 `json:"level_price_on_the_market_up"`

	LevelPriceStopLoss float64 `json:"level_price_stop_loss"`

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

	StopLostBank *StopLostBank `json:"stop_lost_bank,omitempty"`

	StopLostTimes int `json:"stop_lost_times"`

	NeedSendCancel bool `json:"need_send_cancel"`
}

func (s *WingedSwing) Type() string {
	return "winged_swing"
}

func (s *WingedSwing) String() string {
	return "winged_swing"
}

func (s *WingedSwing) Status() smp.StartegyStatus {
	return smp.StartegyStatus{
		IsOnline: s.IsOnline,
	}
}
func (s *WingedSwing) Json() string {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
func (s *WingedSwing) Command(cmd smp.Command, params map[string]string) (res smp.CommandResult, ok bool, err *mft.Error) {
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

	if cmd == SetLevelPriceOnTheMarketUp {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000411, er0, cmd, fS)
			}
			s.LevelPriceOnTheMarketUp = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000410, cmd)
		}
	}
	if cmd == SetLevelPriceOnTheMarketDown {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000431, er0, cmd, fS)
			}
			s.LevelPriceOnTheMarketDown = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000430, cmd)
		}
	}
	if cmd == SetLevelPriceOnTheMarketDownByMarket {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000433, er0, cmd, fS)
			}
			s.LevelPriceOnTheMarketDownByMarket = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000432, cmd)
		}
	}

	if cmd == SetLevelUp {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000413, er0, cmd, fS)
			}
			s.LevelPriceUp = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000412, cmd)
		}
	}

	if cmd == SetLevelDown {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000415, er0, cmd, fS)
			}
			s.LevelPriceDown = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000414, cmd)
		}
	}

	if cmd == SetVolume {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseInt(fS, 10, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000421, er0, cmd, fS)
			}
			s.Volume = int(f)
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000420, cmd)
		}
	}

	return res, false, smp.GenerateError(500000400, cmd)
}

func (s *WingedSwing) AllowCommands() map[smp.Command]smp.CommandInfo {
	return map[smp.Command]smp.CommandInfo{
		smp.ShowCommand:  {0, "Отобразить", "", ""},
		smp.StartCommand: {1, "Старт", "", ""},
		smp.StopCommand:  {2, "Стоп", "", ""},
		SetLevelDown:     {3, "Установить уровень покупки", "уровень", "set_level_down 372.25"},
		SetLevelUp:       {4, "Установить уровень продажи", "уровень", "set_level_up 392.80"},
		SetVolume:        {5, "Установить объём:", "объём", "пример set_vol 25"},

		SetLevelPriceOnTheMarketUp: {6, "Установить уровень начала работы стратегии сверху " +
			"(с этого уровня происходит покупка", "уровень", "set_level 420.90"},
		SetLevelPriceOnTheMarketDown: {7, "Установить уровень начала работы стратегии снизу " +
			"(с этого уровня происходит покупка", "уровень", "set_level 345.67"},
		SetLevelPriceOnTheMarketDownByMarket: {8, "Установить уровень начала работы стратегии снизу по цене рынка " +
			"(до этого уровня происходит покупка по рынку)", "уровень", "set_level 370.10"},
	}
}

func (s *WingedSwing) Description() string {
	return `winged_swing - стратегия, качель`
}

func (s *WingedSwing) ComputeVolume() int {
	if s.StopLostTimes > 0 {
		if s.StopLostTimes > 3 {
			return s.Volume + 3
		}
		return s.Volume + s.StopLostTimes
	}
	return s.Volume
}

func (s *WingedSwing) Step(p smp.StepParams) (meta smp.MetaForStep, err *mft.Error) {
	panic("not implemet correct")

	meta.Name = s.Name
	ob, err := p.GetOrderBook(s.InstrumentId, s.Ticker)
	if err != nil {
		return meta, smp.GenerateErrorE(500000500, err)
	}

	computeVolume := s.ComputeVolume()

	if s.OrderIdBuy == "" && s.OrderIdSell == "" && s.Volume > 0 {
		// Заполняем если установили объёмы руками
		if s.IsBought && s.InMarket <= 0 {
			meta.HasChanges = true
			s.IsBought = false
			s.InMarketPrice = 0
		}

		if !s.IsBought && s.InMarket >= computeVolume {
			meta.HasChanges = true
			s.IsBought = true
			s.InMarketPrice = smp.Round(s.LevelPriceDown*float64(s.InMarket), 6)
		}
	}

	// Проверяем что Покупка прошла
	if s.OrderIdBuy != "" {
		status, prices, err := p.StatusBuyOrder(s.InstrumentId, s.Ticker, s.OrderIdBuy,
			&smp.MetaForOperations{
				NameOfStrategy: s.Name,
			},
		)
		if err != nil {
			return meta, smp.GenerateErrorE(500000502, err)
		}

		cnt := 0
		price := 0.0
		for _, pr := range prices {
			cnt += pr.Count
			price += pr.Price * float64(pr.Count)
		}

		if status == smp.Complete || status == smp.Canceled {
			s.InMarket += cnt
			s.InMarketPrice = smp.Round(s.InMarketPrice+price, 6)
			s.OrderIdBuy = ""
			s.InMarketWait = 0
			s.InMarketPriceWait = 0

			if s.InMarket >= computeVolume {
				s.IsBought = true
			}

		} else {
			s.InMarketWait = cnt
			s.InMarketPriceWait = price
		}
	}

	// Проверяем что покупка прошла
	if s.OrderIdSell != "" {
		status, prices, err := p.StatusSellOrder(s.InstrumentId, s.Ticker, s.OrderIdSell,
			&smp.MetaForOperations{
				NameOfStrategy: s.Name,
			},
		)
		if err != nil {
			return meta, smp.GenerateErrorE(500000503, err)
		}

		cnt := 0
		price := 0.0
		for _, pr := range prices {
			cnt += pr.Count
			price += pr.Price * float64(pr.Count)
		}

		if status == smp.Complete || status == smp.Canceled {
			s.InMarket -= cnt
			s.InMarketPrice = smp.Round(s.InMarketPrice-price, 6)
			s.OrderIdSell = ""
			s.InMarketWait = 0
			s.InMarketPriceWait = 0

			if s.InMarket <= 0 {
				s.Iteration++

				s.Profit = smp.Round(s.Profit-s.InMarketPrice, 6)
				s.InMarketPrice = 0
				s.IsBought = false
			}

			if cnt > s.Volume {
				s.StopLostTimes = -1
			}
		} else {
			s.InMarketWait = -cnt
			s.InMarketPriceWait = -price
		}
	}

	if ob.TradeStatus != smp.NormalTrading {
		return meta, nil
	}

	if s.NeedSendCancel {
		if s.OrderIdBuy != "" {
			_, err := p.CancelBuyOrder(s.InstrumentId, s.Ticker, s.OrderIdBuy,
				&smp.MetaForOperations{
					NameOfStrategy: s.Name,
				},
			)
			if err != nil {
				return meta, smp.GenerateErrorE(500000513, err)
			}
		}

		if s.OrderIdSell != "" {
			_, err := p.CancelSellOrder(s.InstrumentId, s.Ticker, s.OrderIdSell,
				&smp.MetaForOperations{
					NameOfStrategy: s.Name,
				},
			)
			if err != nil {
				return meta, smp.GenerateErrorE(500000514, err)
			}
		}

		meta.HasChanges = true
		s.NeedSendCancel = false
	}

	if !s.IsOnline {
		return meta, nil
	}

	if ob.SellPrice() < s.LevelPriceStopLoss && s.OrderIdSell != "" {
		meta.HasChanges = true
		s.NeedSendCancel = true
		meta.OpDescr = append(meta.OpDescr, "PriceDown Need cancel Sell")
		return meta, nil
	}

	if ob.SellPrice() < s.LevelPriceStopLoss && s.OrderIdBuy != "" {
		meta.HasChanges = true
		s.NeedSendCancel = true
		meta.OpDescr = append(meta.OpDescr, "PriceDown Need cancel Buy")
		return meta, nil
	}

	if s.StopLostBank != nil {
		if ob.SellPrice() < s.LevelPriceStopLoss && s.InMarket > 0 {
			err := s.StopLostBank.Sell(p,
				s.InstrumentId,
				s.Ticker,
				s.InMarket, s.LevelPriceUp)
			if err != nil {
				return meta, smp.GenerateErrorE(500000515, err)
			}
			s.StopLostTimes = s.StopLostTimes + 1
			s.InMarketPrice = smp.Round(s.InMarketPrice-float64(s.InMarket)*ob.SellPrice(), 6)
			s.InMarket = 0

			meta.HasChanges = true
			meta.OpDescr = append(meta.OpDescr, "Stop loss sale")
			return meta, nil
		}
	}

	if s.StopLostBank != nil {
		if s.IsBought && s.InMarket > 0 &&
			// Текущая цена покупки выше чем цена покупки стратегии
			ob.BuyPrice() >= s.LevelPriceDown {
			if req := s.StopLostBank.AllowRequestCount(s.InstrumentId, s.Ticker); req > 0 {
				if s.OrderIdSell != "" {
					meta.HasChanges = true
					s.NeedSendCancel = true
					meta.OpDescr = append(meta.OpDescr, "Need cancel Sell for Request from SLB")
					return meta, nil
				} else {
					if s.InMarket < req {
						req = s.InMarket
					}

					err = s.StopLostBank.Sell(p, s.InstrumentId, s.Ticker, req, s.LevelPriceUp)
					if err != nil {
						return meta, smp.GenerateErrorE(500000516, err)
					}

					s.StopLostTimes = s.StopLostTimes + 1
					s.InMarketPrice = smp.Round(s.InMarketPrice-float64(req)*ob.SellPrice(), 6)
					s.InMarket -= req

					meta.HasChanges = true
					meta.OpDescr = append(meta.OpDescr, "Requet to SLB sale")
					return meta, nil
				}
			}
		}
	}

	if s.IsBought {
		if s.InMarket > 0 {

			if ob.SellPrice() <= s.LevelPriceOnTheMarketUp {
				s.OrderIdSell, err = p.SellByPrice(s.InstrumentId, s.Ticker, s.InMarket, s.LevelPriceUp,
					&smp.MetaForOperations{
						NameOfStrategy: s.Name,
					},
				)
				if err != nil {
					s.OrderIdSell = ""
					return meta, smp.GenerateErrorE(500000504, err)
				}
				meta.HasChanges = true
				meta.OpDescr = append(meta.OpDescr,
					"Заявка на покупку",
				)
			}
			return meta, nil
		} else {
			return meta, smp.GenerateError(500000505, s.InMarket)
		}
	} else {
		if computeVolume-s.InMarket > 0 {
		} else {
			return meta, smp.GenerateError(500000506, s.InMarket, computeVolume)
		}
	}

	return meta, nil
}

/*

*
*
*
*
*

 */

func (s *WingedSwing) Step_old(p smp.StepParams) (meta smp.MetaForStep, err *mft.Error) {
	meta.Name = s.Name

	ob, err := p.GetOrderBook(s.InstrumentId, s.Ticker)
	if err != nil {
		return meta, smp.GenerateErrorE(500000500, err)
	}

	computeVolume := s.ComputeVolume()

	if s.OrderIdBuy == "" && s.OrderIdSell == "" && s.Volume > 0 {
		// Заполняем если установили объёмы руками
		if s.IsBought && s.InMarket <= 0 {
			s.IsBought = false
			s.InMarketPrice = 0
		}

		if !s.IsBought && s.InMarket >= computeVolume {
			s.IsBought = true
			s.InMarketPrice = smp.Round(s.LevelPriceDown*float64(s.InMarket), 6)
		}
	}

	if s.OrderIdBuy != "" {
		status, prices, err := p.StatusBuyOrder(s.InstrumentId, s.Ticker, s.OrderIdBuy,
			&smp.MetaForOperations{
				NameOfStrategy: s.Name,
			},
		)
		if err != nil {
			return meta, smp.GenerateErrorE(500000502, err)
		}

		cnt := 0
		price := 0.0
		for _, pr := range prices {
			cnt += pr.Count
			price += pr.Price * float64(pr.Count)
		}

		if status == smp.Complete || status == smp.Canceled {
			s.InMarket += cnt
			s.InMarketPrice = smp.Round(s.InMarketPrice+price, 6)
			s.OrderIdBuy = ""
			s.InMarketWait = 0
			s.InMarketPriceWait = 0

			if s.InMarket >= computeVolume {
				s.IsBought = true
			}

		} else {
			s.InMarketWait = cnt
			s.InMarketPriceWait = price
		}
	}

	if s.OrderIdSell != "" {
		status, prices, err := p.StatusSellOrder(s.InstrumentId, s.Ticker, s.OrderIdSell,
			&smp.MetaForOperations{
				NameOfStrategy: s.Name,
			},
		)
		if err != nil {
			return meta, smp.GenerateErrorE(500000503, err)
		}

		cnt := 0
		price := 0.0
		for _, pr := range prices {
			cnt += pr.Count
			price += pr.Price * float64(pr.Count)
		}

		if status == smp.Complete || status == smp.Canceled {
			s.InMarket -= cnt
			s.InMarketPrice = smp.Round(s.InMarketPrice-price, 6)
			s.OrderIdSell = ""
			s.InMarketWait = 0
			s.InMarketPriceWait = 0

			if s.InMarket <= 0 {
				s.Iteration++

				s.Profit = smp.Round(s.Profit-s.InMarketPrice, 6)
				s.InMarketPrice = 0
				s.IsBought = false
			}

			if cnt > s.Volume {
				s.StopLostTimes = -1
			}
		} else {
			s.InMarketWait = -cnt
			s.InMarketPriceWait = -price
		}
	}

	computeVolume = s.ComputeVolume()

	if s.Volume == 0 {
		return meta, nil
	}

	if !s.IsOnline {
		return meta, nil
	}

	if ob.TradeStatus != smp.NormalTrading {
		return meta, nil
	}

	if ob.BuyPrice() < s.LevelPriceStopLoss && s.OrderIdSell != "" {
		_, err := p.CancelSellOrder(s.InstrumentId, s.Ticker, s.OrderIdSell,
			&smp.MetaForOperations{
				NameOfStrategy: s.Name,
			},
		)
		if err != nil {
			return meta, smp.GenerateErrorE(500000507, err)
		}
	}

	if !s.IsBought &&
		ob.SellPrice() <= s.LevelPriceUp && ob.BuyPrice() >= s.LevelPriceDown &&
		s.StopLostBank.AllowCount(s.InstrumentId, s.Ticker) > 0 &&
		computeVolume-s.InMarket > 0 {

		success, err := s.StopLostBank.Buy(p, s.InstrumentId, s.Ticker, computeVolume-s.InMarket, s.LevelPriceDown)
		if err != nil {
			return meta, smp.GenerateErrorE(500000511, err)
		}
		s.InMarket += success
		s.InMarketPrice = smp.Round(s.InMarketPrice+float64(success)*s.LevelPriceDown, 6)

		if s.OrderIdBuy != "" {
			_, err = p.CancelBuyOrder(s.InstrumentId, s.Ticker, s.OrderIdBuy,
				&smp.MetaForOperations{
					NameOfStrategy: s.Name,
					IsStopLoss:     true,
				})
			if err != nil {
				return meta, smp.GenerateErrorE(500000512, err)
			}
		}
	}

	if s.OrderIdBuy != "" {
		return meta, nil
	}
	if s.OrderIdSell != "" {
		return meta, nil
	}

	if len(ob.Bids) < 1 || len(ob.Asks) < 1 {
		return meta, nil
	}

	if s.IsBought {
		if s.InMarket > 0 {
			if ob.BuyPrice() <= s.LevelPriceStopLoss {
				if s.StopLostBank != nil {
					err = s.StopLostBank.Sell(p, s.InstrumentId, s.Ticker, s.InMarket, s.LevelPriceUp)
					if err != nil {
						return meta, smp.GenerateErrorE(500000509, err)
					}
					meta.IsStopLoss = true
					s.StopLostTimes = s.StopLostTimes + 1
				} else {
					s.OrderIdSell, err = p.SellByMarket(s.InstrumentId, s.Ticker, s.InMarket,
						&smp.MetaForOperations{
							NameOfStrategy: s.Name,
							IsStopLoss:     true,
						},
					)
					if err != nil {
						s.OrderIdSell = ""
						meta.IsStopLoss = true
						return meta, smp.GenerateErrorE(500000508, err)
					}
					s.StopLostTimes = s.StopLostTimes + 1
				}
			} else {
				if ob.SellPrice() <= s.LevelPriceOnTheMarketUp {
					s.OrderIdSell, err = p.SellByPrice(s.InstrumentId, s.Ticker, s.InMarket, s.LevelPriceUp,
						&smp.MetaForOperations{
							NameOfStrategy: s.Name,
						},
					)
					if err != nil {
						s.OrderIdSell = ""
						return meta, smp.GenerateErrorE(500000504, err)
					}
				}
			}
		} else {
			return meta, smp.GenerateError(500000505, s.InMarket)
		}
	} else {
		if computeVolume-s.InMarket > 0 {
			if ob.BuyPrice() >= s.LevelPriceOnTheMarketDown && ob.SellPrice() <= s.LevelPriceOnTheMarketUp {
				if s.StopLostBank != nil {
					success, err := s.StopLostBank.Buy(p, s.InstrumentId, s.Ticker, computeVolume-s.InMarket, s.LevelPriceDown)
					if err != nil {
						return meta, smp.GenerateErrorE(500000510, err)
					}
					s.InMarket += success
					s.InMarketPrice = smp.Round(s.InMarketPrice+float64(success)*s.LevelPriceDown, 6)
				}

				if ob.BuyPrice() < s.LevelPriceOnTheMarketDownByMarket {
					s.OrderIdBuy, err = p.BuyByMarket(s.InstrumentId, s.Ticker, computeVolume-s.InMarket,
						&smp.MetaForOperations{
							NameOfStrategy: s.Name,
						},
					)
				} else {
					s.OrderIdBuy, err = p.BuyByPrice(s.InstrumentId, s.Ticker, computeVolume-s.InMarket, s.LevelPriceDown,
						&smp.MetaForOperations{
							NameOfStrategy: s.Name,
						},
					)
				}
				if err != nil {
					s.OrderIdBuy = ""
					return meta, smp.GenerateErrorE(500000501, err)
				}
			}
		} else {
			return meta, smp.GenerateError(500000506, s.InMarket, computeVolume)
		}
	}
	return meta, nil
}
