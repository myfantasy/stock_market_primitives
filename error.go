package smp

import (
	"fmt"

	"github.com/myfantasy/mft"
)

// Errors codes and description
var Errors map[int]string = map[int]string{
	500000000: "strategies.TakeProfitBuy: Command: `%v` does not exists",
	500000010: "strategies.TakeProfitBuy: Command: `%v` param not set",
	500000011: "strategies.TakeProfitBuy: Command: `%v` param `%v` is not float64",
	500000020: "strategies.TakeProfitBuy: Command: `%v` param not set",
	500000021: "strategies.TakeProfitBuy: Command: `%v` param `%v` is not int",
	500000030: "strategies.TakeProfitBuy: Command: `%v` param not set",
	500000031: "strategies.TakeProfitBuy: Command: `%v` param `%v` is not bool",

	500000100: "strategies.TakeProfitBuy: Step: fail order book get",
	500000101: "strategies.TakeProfitBuy: Step: fail buy by price",
	500000102: "strategies.TakeProfitBuy: Step: fail get info about buy order",

	500000200: "strategies.TakeProfitSell: Command: `%v` does not exists",
	500000210: "strategies.TakeProfitSell: Command: `%v` param not set",
	500000211: "strategies.TakeProfitSell: Command: `%v` param `%v` is not float64",
	500000220: "strategies.TakeProfitSell: Command: `%v` param not set",
	500000221: "strategies.TakeProfitSell: Command: `%v` param `%v` is not int",
	500000230: "strategies.TakeProfitSell: Command: `%v` param not set",
	500000231: "strategies.TakeProfitSell: Command: `%v` param `%v` is not bool",

	500000300: "strategies.TakeProfitSell: Step: fail order book get",
	500000301: "strategies.TakeProfitSell: Step: fail sell by price",
	500000302: "strategies.TakeProfitSell: Step: fail get info about sell order",

	500000400: "strategies.WingedSwing: Command: `%v` does not exists",
	500000410: "strategies.WingedSwing: Command: `%v` param not set",
	500000411: "strategies.WingedSwing: Command: `%v` param `%v` is not float64",
	500000412: "strategies.WingedSwing: Command: `%v` param not set",
	500000413: "strategies.WingedSwing: Command: `%v` param `%v` is not float64",
	500000414: "strategies.WingedSwing: Command: `%v` param not set",
	500000415: "strategies.WingedSwing: Command: `%v` param `%v` is not float64",
	500000420: "strategies.WingedSwing: Command: `%v` param not set",
	500000421: "strategies.WingedSwing: Command: `%v` param `%v` is not int",

	500000500: "strategies.WingedSwing: Step: fail order book get",
	500000501: "strategies.WingedSwing: Step: fail buy by price",
	500000502: "strategies.WingedSwing: Step: fail get info about buy order",
	500000503: "strategies.WingedSwing: Step: fail get info about sell order",
	500000504: "strategies.WingedSwing: Step: fail sell by price",
	500000505: "strategies.WingedSwing: Step: unexpected situation InMarket: %v",
	500000506: "strategies.WingedSwing: Step: unexpected situation InMarket: %v; s.Volume: %v",

	500000600: "strategies.WingedSwing: Command: `%v` does not exists",
	500000610: "strategies.WingedSwing: Command: `%v` param not set",
	500000611: "strategies.WingedSwing: Command: `%v` param `%v` is not float64",
	500000612: "strategies.WingedSwing: Command: `%v` param not set",
	500000613: "strategies.WingedSwing: Command: `%v` param `%v` is not float64",
	500000614: "strategies.WingedSwing: Command: `%v` param not set",
	500000615: "strategies.WingedSwing: Command: `%v` param `%v` is not float64",
	500000616: "strategies.WingedSwing: Command: `%v` param not set",
	500000617: "strategies.WingedSwing: Command: `%v` param `%v` is not float64",
	500000618: "strategies.WingedSwing: Command: `%v` param not set",
	500000619: "strategies.WingedSwing: Command: `%v` param `%v` is not float64",
	500000620: "strategies.WingedSwing: Command: `%v` param not set",
	500000621: "strategies.WingedSwing: Command: `%v` param `%v` is not int",
	500000640: "strategies.WingedSwing: Command: `%v` param `%v` not set",
	500000641: "strategies.WingedSwing: Command: `%v` value `%v` is not int (param `%v`)",
	500000642: "strategies.WingedSwing: Command: `%v` param `%v` not set",
	500000643: "strategies.WingedSwing: Command: `%v` value `%v` is not int (param `%v`)",
	500000644: "strategies.WingedSwing: Command: `%v` param `%v` not set",
	500000645: "strategies.WingedSwing: Command: `%v` value `%v` is not int (param `%v`)",
	500000646: "strategies.WingedSwing: Command: `%v` param `%v` not set",
	500000647: "strategies.WingedSwing: Command: `%v` value `%v` is not int (param `%v`)",
	500000650: "strategies.WingedSwing: Command: `%v` lavel `i` does not exists",
	500000651: "strategies.WingedSwing: Command: `%v` lavel `i` value `%v` is not int (element `%v`)",

	500000700: "strategies.WingedSwing: Step: fail do some nested steps faild: %v of %v",
}

// GenerateError -
func GenerateError(key int, a ...interface{}) *mft.Error {
	if text, ok := Errors[key]; ok {
		return mft.ErrorCS(key, fmt.Sprintf(text, a...))
	}
	panic(fmt.Sprintf("smp.GenerateError, error not found code:%v", key))
}

// GenerateErrorE -
func GenerateErrorE(key int, err error, a ...interface{}) *mft.Error {
	if text, ok := Errors[key]; ok {
		return mft.ErrorCSEf(key, err, text, a...)
	}
	panic(fmt.Sprintf("smp.GenerateErrorE, error not found code:%v error:%v", key, err))
}

// GenerateError -
func GenerateErrorSubList(key int, sub []*mft.Error, a ...interface{}) *mft.Error {
	if text, ok := Errors[key]; ok {
		err := mft.ErrorCS(key, fmt.Sprintf(text, a...))
		err.InternalErrors = sub
		return err
	}
	panic(fmt.Sprintf("smp.GenerateError, error not found code:%v", key))
}
