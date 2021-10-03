package strategies

import (
	"encoding/json"
	"strconv"

	"github.com/myfantasy/mfs"
	"github.com/myfantasy/mft"
	smp "github.com/myfantasy/stock_market_primitives"
)

//go:generate mfjson winged_swing_group.go

// Крылатые качели группа (WingedSwingGroup)

var (
	_ smp.Strategy = &WingedSwingGroup{}
)

//mfjson:interface smp.strategies.winged_swing_group
type WingedSwingGroup struct {
	Name string `json:"name"`

	InstrumentId string `json:"instrument_id"`
	Ticker       string `json:"ticker"`

	Volume       int     `json:"volume"`
	LevelPrice   float64 `json:"level_price"`
	PriceUp      float64 `json:"price_up"`
	PriceDown    float64 `json:"price_down"`
	PriceBetween float64 `json:"price_between"`

	PriceOnTheMarketUp           float64 `json:"price_on_the_market_up"`
	PriceOnTheMarketDown         float64 `json:"price_on_the_market_down"`
	PriceOnTheMarketDownByMarket float64 `json:"price_on_the_market_down_by_market"`

	IsOnline bool `json:"is_online"`

	Swings []WingedSwing `json:"swings"`

	StopLostBank *StopLostBank `json:"stop_lost_bank,omitempty"`

	mx mfs.PMutex
}

func (s *WingedSwingGroup) Type() string {
	return "winged_swing_group"
}

func (s *WingedSwingGroup) String() string {
	return "winged_swing_group"
}

func (s *WingedSwingGroup) Status() smp.StartegyStatus {
	return smp.StartegyStatus{
		IsOnline: s.IsOnline,
	}
}
func (s *WingedSwingGroup) Json() string {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
func (s *WingedSwingGroup) Command(cmd smp.Command, params map[string]string) (res smp.CommandResult, ok bool, err *mft.Error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if cmd == smp.ShowCommand {
		res.Message = s.String()
		return res, true, nil
	}

	if cmd == smp.StartCommand {
		s.IsOnline = true
		for i := range s.Swings {
			s.Swings[i].IsOnline = true
		}
		return res, true, nil
	}

	if cmd == smp.StopCommand {
		s.IsOnline = false
		for i := range s.Swings {
			s.Swings[i].IsOnline = false
		}
		return res, true, nil
	}

	if cmd == SetLevel {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000611, er0, cmd, fS)
			}
			s.LevelPrice = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000610, cmd)
		}
	}

	if cmd == SetPriceUp {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000613, er0, cmd, fS)
			}
			s.PriceUp = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000612, cmd)
		}
	}

	if cmd == SetPriceDown {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000615, er0, cmd, fS)
			}
			s.PriceDown = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000614, cmd)
		}
	}

	if cmd == SetPriceOnTheMarketDownByMarket {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000617, er0, cmd, fS)
			}
			s.PriceOnTheMarketDownByMarket = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000616, cmd)
		}
	}

	if cmd == SetPriceOnTheMarketDown {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000661, er0, cmd, fS)
			}
			s.PriceOnTheMarketDown = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000660, cmd)
		}
	}

	if cmd == SetPriceOnTheMarketUp {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000663, er0, cmd, fS)
			}
			s.PriceOnTheMarketUp = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000662, cmd)
		}
	}

	if cmd == SetPriceBetween {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000619, er0, cmd, fS)
			}
			s.PriceBetween = f
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000618, cmd)
		}
	}

	if cmd == SetVolume {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseInt(fS, 10, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000621, er0, cmd, fS)
			}
			s.Volume = int(f)
			return res, true, nil
		} else {
			return res, false, smp.GenerateError(500000620, cmd)
		}
	}

	if cmd == Render {
		fS, ok := params["f"]
		if !ok {
			return res, false, smp.GenerateError(500000640, cmd, "f")
		}
		f, er0 := strconv.ParseInt(fS, 10, 64)
		if er0 != nil {
			return res, false, smp.GenerateErrorE(500000641, er0, cmd, fS, "f")
		}
		tS, ok := params["t"]
		if !ok {
			return res, false, smp.GenerateError(500000642, cmd, "t")
		}
		t, er0 := strconv.ParseInt(tS, 10, 64)
		if er0 != nil {
			return res, false, smp.GenerateErrorE(500000643, er0, cmd, tS, "t")
		}

		for i := f; i <= t; i++ {
			level := s.LevelPrice + float64(i)*s.PriceBetween
			s.Swings = append(s.Swings, WingedSwing{
				Name:         s.Name + "[" + strconv.Itoa(int(i)) + "]",
				InstrumentId: s.InstrumentId,
				Ticker:       s.Ticker,
				Volume:       s.Volume,

				LevelPriceUp:   smp.Round(level+s.PriceUp, 6),
				LevelPriceDown: smp.Round(level+s.PriceDown, 6),

				LevelPriceOnTheMarketUp:           smp.Round(level+s.PriceOnTheMarketUp, 6),
				LevelPriceOnTheMarketDown:         smp.Round(level+s.PriceOnTheMarketDown, 6),
				LevelPriceOnTheMarketDownByMarket: smp.Round(level+s.PriceOnTheMarketDownByMarket, 6),

				IsOnline: s.IsOnline,

				Labels: map[string]string{"i": strconv.Itoa(int(i))},
			})
		}

		return res, true, nil
	}

	if cmd == SetInMarket {
		fS, ok := params["f"]
		if !ok {
			return res, false, smp.GenerateError(500000644, cmd, "f")
		}
		f, er0 := strconv.ParseInt(fS, 10, 64)
		if er0 != nil {
			return res, false, smp.GenerateErrorE(500000645, er0, cmd, fS, "f")
		}
		tS, ok := params["t"]
		if !ok {
			return res, false, smp.GenerateError(500000646, cmd, "t")
		}
		t, er0 := strconv.ParseInt(tS, 10, 64)
		if er0 != nil {
			return res, false, smp.GenerateErrorE(500000647, er0, cmd, tS, "t")
		}

		for i := range s.Swings {
			iS, ok := s.Swings[i].Labels["i"]
			if !ok {
				return res, false, smp.GenerateError(500000650, cmd)
			}
			iVal, er0 := strconv.ParseInt(iS, 10, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000651, er0, cmd, tS)
			}

			if f <= iVal && t >= iVal {
				s.Swings[i].InMarket = s.Swings[i].Volume
			}
		}

		return res, true, nil
	}

	if cmd == SetOutOfMarket {
		fS, ok := params["f"]
		if !ok {
			return res, false, smp.GenerateError(500000644, cmd, "f")
		}
		f, er0 := strconv.ParseInt(fS, 10, 64)
		if er0 != nil {
			return res, false, smp.GenerateErrorE(500000645, er0, cmd, fS, "f")
		}
		tS, ok := params["t"]
		if !ok {
			return res, false, smp.GenerateError(500000646, cmd, "t")
		}
		t, er0 := strconv.ParseInt(tS, 10, 64)
		if er0 != nil {
			return res, false, smp.GenerateErrorE(500000647, er0, cmd, tS, "t")
		}

		for i := range s.Swings {
			iS, ok := s.Swings[i].Labels["i"]
			if !ok {
				return res, false, smp.GenerateError(500000650, cmd)
			}
			iVal, er0 := strconv.ParseInt(iS, 10, 64)
			if er0 != nil {
				return res, false, smp.GenerateErrorE(500000651, er0, cmd, tS)
			}

			if f <= iVal && t >= iVal {
				s.Swings[i].InMarket = 0
			}
		}

		return res, true, nil
	}

	return res, false, smp.GenerateError(500000600, cmd)
}

func (s *WingedSwingGroup) AllowCommands() map[smp.Command]smp.CommandInfo {
	return map[smp.Command]smp.CommandInfo{
		smp.ShowCommand:  {0, "Отобразить", "", ""},
		smp.StartCommand: {1, "Старт", "", ""},
		smp.StopCommand:  {2, "Стоп", "", ""},
		SetLevel:         {3, "Установить уровень начала работы стратегии (с этого уровня происходит распределение стратегии)", "уровень", "set_level 345.67"},

		SetPriceUp:   {4, "Установить шаг продажи", "шаг цены", "set_price_up 12.25"},
		SetPriceDown: {5, "Установить шаг покупки", "шаг цены", "set_price_down -0.80"},

		SetPriceOnTheMarketUp: {6, "Установить шаг входа в рынок", "шаг цены",
			"set_price_on_the_market_up 60.20"},
		SetPriceOnTheMarketDown: {7, "Установить шаг входа в рынок", "шаг цены",
			"set_price_on_the_market_down -50.80"},
		SetPriceOnTheMarketDownByMarket: {8, "Установить шаг входа в рынок", "шаг цены",
			"set_price_on_the_market_down_by_market -10.80"},

		SetPriceBetween: {9, "Установить шаг между стратегиями", "уровень", "set_price_between 10.2"},

		SetVolume: {10, "Установить объём", "объём (кол-во лотов)", "set_vol 25"},

		Render: {11, "Сгенерировать внутренние стратегии",
			"f=[от шагов] t=[до шагов]", "render f=-5 t=10"},

		SetInMarket: {12, "Установить кол-во акций купленных на рынке, по максимальному объёму ",
			"f=[от шагов] t=[до шагов]) ",
			"set_in_market f=-4 t=8"},
		SetOutOfMarket: {13, "Установить кол-во акций купленных на рынке, в 0 ",
			"f=[от шагов] t=[до шагов]",
			"set_out_of_market f=-4 t=8"},
	}
}

func (s *WingedSwingGroup) Description() string {
	return "`winged_swing_group`" + ` - стратегия, группы качель`
}

func (s *WingedSwingGroup) Step(p smp.StepParams) (meta smp.MetaForStep, err *mft.Error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	meta.Name = s.Name

	for i := range s.Swings {
		s.Swings[i].StopLostBank = s.StopLostBank
		mt, er := s.Swings[i].Step(p)
		if mt.HasChanges {
			meta.HasChanges = true
		}
		meta.SubMeta = append(meta.SubMeta, mt)
		err.AppendList(er)
		s.Swings[i].StopLostBank = nil
	}

	if err != nil {
		return meta, smp.GenerateErrorSubList(500000700, err.InternalErrors, len(err.InternalErrors), len(s.Swings))
	}

	return meta, nil
}
