package smp

import (
	"sort"
	"time"
)

const H24 = time.Hour * 24

type Dividend struct {
	InstrumentId string `json:"instrument_id"`
	Ticker       string `json:"ticker"`

	Amount float64 `json:"amount"`
	// LastDate - Last date for Buy
	LastDate time.Time `json:"last_date"`

	// sum of all dividends amount
	AggDividend float64 `json:"agg_dividend"`
}

type Dividends []Dividend

func (ds Dividends) Sort()              { sort.Sort(ds) }
func (ds Dividends) Len() int           { return len(ds) }
func (ds Dividends) Swap(i, j int)      { ds[i], ds[j] = ds[j], ds[i] }
func (ds Dividends) Less(i, j int) bool { return ds[i].LastDate.Before(ds[j].LastDate) }

func (d Dividend) AfterOrEqual24(t time.Time) bool {
	return d.LastDate.Add(H24).After(t) || d.LastDate.Add(H24).Equal(t)
}

func (ds Dividends) FindCandleDividend(
	start time.Time, end time.Time,
) (
	lastAmount float64, lastDate time.Time, agg float64,
) {
	if ds.Len() == 0 {
		return // 0
	}
	ix := sort.Search(ds.Len(), func(i int) bool {
		return ds[i].AfterOrEqual24(start)
	})

	for {
		if ds.Len() == ix {
			break
		}
		if ds[ix].AfterOrEqual24(end) {
			break
		}
		ix++
	}
	if ix == 0 {
		return // 0
	}

	return ds[ix-1].Amount, ds[ix-1].LastDate, ds[ix-1].AggDividend
}
