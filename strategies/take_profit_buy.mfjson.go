// Code generated by mfjson for marshaling/unmarshaling. DO NOT EDIT.
// https://github.com/myfantasy/json

package strategies

import (
	mfj "github.com/myfantasy/json"
)
func (obj *TakeProfitBuy) UnmarshalJSONTypeName() string {
	return "smp.strategies.take_profit_buy"
}

func init() {
	mfj.GlobalStructFactory.Add("smp.strategies.take_profit_buy", func() mfj.JsonInterfaceMarshaller { return &TakeProfitBuy{} })
	mfj.GlobalStructFactory.AddNil("smp.strategies.take_profit_buy", func() mfj.JsonInterfaceMarshaller {
		var out *TakeProfitBuy
		return out
	})
}
