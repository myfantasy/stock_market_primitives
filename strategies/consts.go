package strategies

import (
	smp "github.com/myfantasy/stock_market_primitives"
)

const (
	SetLevel        smp.Command = "set_level"
	SetVolume       smp.Command = "set_vol"
	SetStayInMarket smp.Command = "stay_in_market"

	SetLevelUp   smp.Command = "set_level_up"
	SetLevelDown smp.Command = "set_level_down"

	SetPriceUp      smp.Command = "set_price_up"
	SetPriceDown    smp.Command = "set_price_down"
	SetPriceBuy     smp.Command = "set_price_buy"
	SetPriceBetween smp.Command = "set_price_between"

	Render smp.Command = "render"

	SetInMarket    smp.Command = "set_in_market"
	SetOutOfMarket smp.Command = "set_out_of_market"
)
