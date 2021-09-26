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
	SetPriceBetween smp.Command = "set_price_between"

	SetLevelPriceOnTheMarketUp           smp.Command = "set_level_on_the_market_up"
	SetLevelPriceOnTheMarketDown         smp.Command = "set_level_on_the_market_down"
	SetLevelPriceOnTheMarketDownByMarket smp.Command = "set_level_on_the_market_down_by_market"

	SetPriceOnTheMarketUp           smp.Command = "set_price_on_the_market_up"
	SetPriceOnTheMarketDown         smp.Command = "set_price_on_the_market_down"
	SetPriceOnTheMarketDownByMarket smp.Command = "set_price_on_the_market_down_by_market"

	Render smp.Command = "render"

	SetInMarket    smp.Command = "set_in_market"
	SetOutOfMarket smp.Command = "set_out_of_market"
)
