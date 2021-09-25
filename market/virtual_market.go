package market

import (
	"time"

	"github.com/myfantasy/mft"
	smp "github.com/myfantasy/stock_market_primitives"
)

var (
	_ smp.StepParams = &VirtualMarket{}
)

type VirtualMarket struct {
	Candles        smp.Candles
	OrderBook      *smp.OrderBook
	InstrumentInfo *smp.InstrumentInfo
	Position       int
}

func (vm *VirtualMarket) GetCandles(instrumentId string, ticker string, dateFrom time.Time, dateTo time.Time) (cs smp.Candles, err *mft.Error) {
	return vm.Candles[0:vm.Position].After(dateFrom).Before(dateTo).Before(vm.OrderBook.Time).Clone(),
		nil
}
func (vm *VirtualMarket) GetOrderBook(instrumentId string, ticker string) (ob *smp.OrderBook, err *mft.Error) {
	return vm.OrderBook, nil
}
func (vm *VirtualMarket) GetInstrumentInfo(instrumentId string, ticker string) (instrumentInfo *smp.InstrumentInfo, err *mft.Error) {
	return vm.InstrumentInfo, nil
}
func (vm *VirtualMarket) BuyByMarket(instrumentId string, ticker string, cnt int) (orderId string, err *mft.Error) {
	return "", nil
}
func (vm *VirtualMarket) SellByMarket(instrumentId string, ticker string, cnt int) (orderId string, err *mft.Error) {
	return "", nil
}
func (vm *VirtualMarket) BuyByPrice(instrumentId string, ticker string, cnt int, price float64) (orderId string, err *mft.Error) {
	return "", nil
}
func (vm *VirtualMarket) SellByPrice(instrumentId string, ticker string, cnt int, price float64) (orderId string, err *mft.Error) {
	return "", nil
}
func (vm *VirtualMarket) CancelBuyOrder(instrumentId string, ticker string, orderId string) (ok bool, err *mft.Error) {
	return false, nil
}
func (vm *VirtualMarket) CancelSellOrder(instrumentId string, ticker string, orderId string) (ok bool, err *mft.Error) {
	return false, nil
}
func (vm *VirtualMarket) StatusBuyOrder(instrumentId string, ticker string, orderId string) (status smp.StatusOrder, prices []smp.LotPrices, err *mft.Error) {
	return smp.Complete, make([]smp.LotPrices, 0), nil
}
func (vm *VirtualMarket) StatusSellOrder(instrumentId string, ticker string, orderId string) (status smp.StatusOrder, prices []smp.LotPrices, err *mft.Error) {
	return smp.Complete, make([]smp.LotPrices, 0), nil
}

func (vm *VirtualMarket) DoStep() bool {
	if vm.Position >= vm.Candles.Len() {
		return false
	}
	vm.Position++
	if vm.Position >= vm.Candles.Len() {
		return false
	}

	vm.OrderBook = vm.Candles[vm.Position].OrderBook()
	return true
}
