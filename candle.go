package smp

import (
	"math"
	"sort"
	"time"
)

type Candle struct {
	InstrumentId string `json:"instrument_id"`
	Ticker       string `json:"ticker"`

	// End date
	Date  time.Time `json:"date"`
	Start time.Time `json:"start"`
	Open  float64   `json:"open"`
	High  float64   `json:"high"`
	Low   float64   `json:"low"`
	Close float64   `json:"close"`
	Vol   int       `json:"vol"`

	// Last dividend amount
	LastDividend float64 `json:"last_dividend,omitempty"`
	// First date after [last dividend buy date]
	LastDividendDate time.Time `json:"last_dividend_date,omitempty"`
	// sum of all dividends amount
	AggDividend float64 `json:"agg_dividend,omitempty"`
}

type Candles []Candle

func (cs Candles) Sort()              { sort.Sort(cs) }
func (cs Candles) Len() int           { return len(cs) }
func (cs Candles) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs Candles) Less(i, j int) bool { return cs[i].Date.Before(cs[j].Date) }

// After returns candels after date include it (after or equal)
func (cs Candles) After(date time.Time) (out Candles) {
	if cs.Len() == 0 {
		return cs
	}

	ix := sort.Search(cs.Len(), func(i int) bool {
		return cs[i].Date.After(date) || cs[i].Date.Equal(date)
	})

	return cs[ix:]
}

// Befor returns candels befor date exclude it (befor strong)
func (cs Candles) Before(date time.Time) (out Candles) {
	if cs.Len() == 0 {
		return cs
	}

	ix := sort.Search(cs.Len(), func(i int) bool {
		return cs[i].Date.After(date) || cs[i].Date.Equal(date)
	})

	return cs[:ix]
}
func (cs Candles) Clone() (res Candles) {
	res = make(Candles, 0, cs.Len())
	for i := 0; i < cs.Len(); i++ {
		res = append(res, cs[i])
	}

	return res
}

// Aggregate - Aggregate Candles by aggFrame (cs should be sorted)
func (cs Candles) Aggregate(aggFrame time.Duration) (out Candles) {
	out = make(Candles, 0)

	current := Candle{
		InstrumentId:     cs[0].InstrumentId,
		Ticker:           cs[0].Ticker,
		Date:             cs[0].Start.Add(aggFrame),
		Start:            cs[0].Start,
		Open:             cs[0].Open,
		High:             cs[0].High,
		Low:              cs[0].Low,
		Close:            cs[0].Close,
		Vol:              cs[0].Vol,
		LastDividend:     cs[0].LastDividend,
		LastDividendDate: cs[0].LastDividendDate,
		AggDividend:      cs[0].AggDividend,
	}
	if cs.Len() == 0 {
		return out
	}

	for i, c := range cs {
		if i == 0 {
			continue
		}

		if c.Date.After(current.Date) {
			out = append(out, current)
			current = Candle{
				InstrumentId:     c.InstrumentId,
				Ticker:           c.Ticker,
				Date:             c.Start.Add(aggFrame),
				Start:            c.Start,
				Open:             c.Open,
				High:             c.High,
				Low:              c.Low,
				Close:            c.Close,
				Vol:              c.Vol,
				LastDividend:     c.LastDividend,
				LastDividendDate: c.LastDividendDate,
				AggDividend:      c.AggDividend,
			}
			continue
		}

		current.High = math.Max(current.High, c.High)
		current.Low = math.Max(current.Low, c.Low)
		current.Close = c.Close
		current.Vol += c.Vol

		if !current.LastDividendDate.Equal(c.LastDividendDate) {
			current.LastDividendDate = c.LastDividendDate
			current.AggDividend = c.AggDividend
			current.LastDividend = c.LastDividend
		}
	}

	out = append(out, current)
	return out
}

func (cs Candles) FillDividents(ds Dividends) (out Candles) {
	out = make(Candles, 0)

	for _, c := range cs {
		current := Candle{
			InstrumentId: c.InstrumentId,
			Ticker:       c.Ticker,
			Date:         c.Date,
			Start:        c.Start,
			Open:         c.Open,
			High:         c.High,
			Low:          c.Low,
			Close:        c.Close,
			Vol:          c.Vol,
		}

		current.LastDividend, current.LastDividendDate, current.AggDividend =
			ds.FindCandleDividend(current.Start, current.Date)

		out = append(out, current)
	}

	return out
}

// OrderBook Generates order_book
func (c *Candle) OrderBook() OrderBook {
	ob := OrderBook{
		InstrumentId: c.InstrumentId,
		Ticker:       c.Ticker,

		Time:  c.Date,
		Depth: 1,
		Bids: []RestPriceQuantity{{
			Price:    c.Low,
			Quantity: c.Vol,
		}},
		Asks: []RestPriceQuantity{{
			Price:    c.High,
			Quantity: c.Vol,
		}},
		TradeStatus:       NormalTrading,
		MinPriceIncrement: 0.0001,
		LastPrice:         c.Close,
		ClosePrice:        c.Close,
		LimitUp:           c.High,
		LimitDown:         c.Low,
	}
	return ob
}
