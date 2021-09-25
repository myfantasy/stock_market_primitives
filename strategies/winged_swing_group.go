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
	InstrumentId string `json:"instrument_id"`
	Ticker       string `json:"ticker"`

	Volume       int     `json:"volume"`
	LevelPrice   float64 `json:"level_price"`
	PriceUp      float64 `json:"price_up"`
	PriceDown    float64 `json:"price_down"`
	PriceBuy     float64 `json:"price_buy"`
	PriceBetween float64 `json:"price_between"`

	IsOnline bool `json:"is_online"`

	Swings []WingedSwing `json:"swings"`

	mx mfs.PMutex
}

func (s *WingedSwingGroup) Status() smp.StartegyStatus {
	return smp.StartegyStatus{
		IsOnline: s.IsOnline,
	}
}
func (s *WingedSwingGroup) String() string {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
func (s *WingedSwingGroup) Command(cmd smp.Command, params map[string]string) (ok bool, err *mft.Error) {
	s.mx.Lock()
	defer s.mx.Unlock()
	if cmd == smp.ShowCommand {
		return true, nil
	}

	if cmd == smp.StartCommand {
		s.IsOnline = true
		for i := range s.Swings {
			s.Swings[i].IsOnline = true
		}
		return true, nil
	}

	if cmd == smp.StopCommand {
		s.IsOnline = false
		for i := range s.Swings {
			s.Swings[i].IsOnline = false
		}
		return true, nil
	}

	if cmd == SetLevel {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000611, er0, cmd, fS)
			}
			s.LevelPrice = f
		} else {
			return false, smp.GenerateError(500000610, cmd)
		}
	}

	if cmd == SetPriceUp {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000613, er0, cmd, fS)
			}
			s.PriceUp = f
		} else {
			return false, smp.GenerateError(500000612, cmd)
		}
	}

	if cmd == SetPriceDown {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000615, er0, cmd, fS)
			}
			s.PriceDown = f
		} else {
			return false, smp.GenerateError(500000614, cmd)
		}
	}

	if cmd == SetPriceBuy {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000617, er0, cmd, fS)
			}
			s.PriceBuy = f
		} else {
			return false, smp.GenerateError(500000616, cmd)
		}
	}

	if cmd == SetPriceBetween {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseFloat(fS, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000619, er0, cmd, fS)
			}
			s.PriceBetween = f
		} else {
			return false, smp.GenerateError(500000618, cmd)
		}
	}

	if cmd == SetVolume {
		fS, ok := params[""]
		if ok {
			f, er0 := strconv.ParseInt(fS, 10, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000621, er0, cmd, fS)
			}
			s.Volume = int(f)
		} else {
			return false, smp.GenerateError(500000620, cmd)
		}
	}

	if cmd == Render {
		fS, ok := params["f"]
		if !ok {
			return false, smp.GenerateError(500000640, cmd, "f")
		}
		f, er0 := strconv.ParseInt(fS, 10, 64)
		if er0 != nil {
			return false, smp.GenerateErrorE(500000641, er0, cmd, fS, "f")
		}
		tS, ok := params["t"]
		if !ok {
			return false, smp.GenerateError(500000642, cmd, "t")
		}
		t, er0 := strconv.ParseInt(tS, 10, 64)
		if er0 != nil {
			return false, smp.GenerateErrorE(500000643, er0, cmd, tS, "t")
		}

		for i := f; i <= t; i++ {
			level := s.LevelPrice + float64(i)*s.PriceBetween
			s.Swings = append(s.Swings, WingedSwing{
				InstrumentId: s.InstrumentId,
				Ticker:       s.Ticker,
				Volume:       s.Volume,

				LevelPriceUp:   level + s.PriceUp,
				LevelPriceDown: level + s.PriceDown,
				LevelPriceBuy:  level + s.PriceBuy,

				IsOnline: s.IsOnline,

				Labels: map[string]string{"i": strconv.Itoa(int(i))},
			})
		}
	}

	if cmd == SetInMarket {
		fS, ok := params["f"]
		if !ok {
			return false, smp.GenerateError(500000644, cmd, "f")
		}
		f, er0 := strconv.ParseInt(fS, 10, 64)
		if er0 != nil {
			return false, smp.GenerateErrorE(500000645, er0, cmd, fS, "f")
		}
		tS, ok := params["t"]
		if !ok {
			return false, smp.GenerateError(500000646, cmd, "t")
		}
		t, er0 := strconv.ParseInt(tS, 10, 64)
		if er0 != nil {
			return false, smp.GenerateErrorE(500000647, er0, cmd, tS, "t")
		}

		for i := range s.Swings {
			iS, ok := s.Swings[i].Labels["i"]
			if !ok {
				return false, smp.GenerateError(500000650, cmd)
			}
			iVal, er0 := strconv.ParseInt(iS, 10, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000651, er0, cmd, tS)
			}

			if f <= iVal && t >= iVal {
				s.Swings[i].InMarket = s.Swings[i].Volume
			}
		}
	}

	if cmd == SetOutOfMarket {
		fS, ok := params["f"]
		if !ok {
			return false, smp.GenerateError(500000644, cmd, "f")
		}
		f, er0 := strconv.ParseInt(fS, 10, 64)
		if er0 != nil {
			return false, smp.GenerateErrorE(500000645, er0, cmd, fS, "f")
		}
		tS, ok := params["t"]
		if !ok {
			return false, smp.GenerateError(500000646, cmd, "t")
		}
		t, er0 := strconv.ParseInt(tS, 10, 64)
		if er0 != nil {
			return false, smp.GenerateErrorE(500000647, er0, cmd, tS, "t")
		}

		for i := range s.Swings {
			iS, ok := s.Swings[i].Labels["i"]
			if !ok {
				return false, smp.GenerateError(500000650, cmd)
			}
			iVal, er0 := strconv.ParseInt(iS, 10, 64)
			if er0 != nil {
				return false, smp.GenerateErrorE(500000651, er0, cmd, tS)
			}

			if f <= iVal && t >= iVal {
				s.Swings[i].InMarket = 0
			}
		}
	}

	return false, smp.GenerateError(500000600, cmd)
}

func (s *WingedSwingGroup) AllowCommands() map[smp.Command]string {
	return map[smp.Command]string{
		smp.ShowCommand:  "Отобразить",
		smp.StartCommand: "Старт",
		smp.StopCommand:  "Стоп",
		SetLevel:         "Установить уровень начала работы стратегии (с этого уровня происходит распределение стратегии) (параметр: уровень) пример: set_level 345.67",

		SetPriceUp:      "Установить шаг продажи (параметр: уровень) пример: set_price_up 12.25",
		SetPriceDown:    "Установить шаг покупки (параметр: уровень) пример: set_price_down -0.80",
		SetPriceBuy:     "Установить шаг входа в рынок (параметр: уровень) пример: set_price_buy -50.80",
		SetPriceBetween: "Установить шаг между стратегиями (параметр: уровень) пример: set_price_between 10.2",

		SetVolume: "Установить объём (параметр: объём) пример set_vol 25",

		Render: "Сгенерировать внутренние стратегии " +
			"(параметры: f=[от шагов] t=[до шагов]) пример: render f=-5 t=10",

		SetInMarket: "Установить кол-во акций купленных на рынке, по максимальному объёму " +
			"(параметры: f=[от шагов] t=[до шагов]) " +
			"пример: set_in_market f=-4 t=8",
		SetOutOfMarket: "Установить кол-во акций купленных на рынке, в 0 " +
			"(параметры: f=[от шагов] t=[до шагов]) " +
			"пример: set_out_of_market f=-4 t=8",
	}
}

func (s *WingedSwingGroup) Step(p smp.StepParams) (err *mft.Error) {
	s.mx.Lock()
	defer s.mx.Unlock()

	for i := range s.Swings {
		err.AppendList(s.Swings[i].Step(p))
	}

	if err != nil {
		return smp.GenerateErrorSubList(500000600, err.InternalErrors, len(err.InternalErrors), len(s.Swings))
	}

	return nil
}
