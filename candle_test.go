package smp

import (
	"math"
	"testing"
	"time"
)

func TestCandlesAggregate(t *testing.T) {
	cs := make(Candles, 0)

	tmStart, er0 := time.Parse("20060102 150405 MST", "20090101 000000 MSK")
	if er0 != nil {
		t.Fatal(er0)
	}

	step := 0
	for tm := tmStart; step < 100000; tm = tm.Add(1 * time.Minute) {
		cs = append(cs, Candle{
			Ticker: "TTTT",
			Date:   tm.Add(time.Minute),
			Start:  tm,
			Open:   3 + math.Cos(float64(step)/100*math.Pi),
			High:   3 + math.Cos(float64(step)/100*math.Pi) + 0.5,
			Low:    3 + math.Cos(float64(step)/100*math.Pi) - 0.5,
			Close:  3 + math.Cos(float64(step+1)/100*math.Pi),
			Vol:    1,
		})
		step++
	}

	cs2 := cs.Aggregate(time.Minute * 5)
	if len(cs2) != 20000 {
		t.Fatalf("Aggregate by 5 minute shoult contains 20000 candles (current %v)", len(cs2))
	}
}
