package smp

import (
	"fmt"

	"github.com/myfantasy/mft"
)

// Errors codes and description
var Errors map[int]string = map[int]string{
	25000000: "StrategyGenerator.Create: strategy type `%v` not found",
	25000001: "StrategyGenerator.Load: strategy type `%v` load error",

	25000020: "StrategyStorage.UnmarshalJSON: fail to unmarshal map of JsonTypedContainer",
	25000021: "StrategyStorage.UnmarshalJSON: fail to load strategy type `%v` (used smp.DefaultStrategyGenerator)",
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
