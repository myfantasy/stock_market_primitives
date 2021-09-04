package smp

import (
	"fmt"

	"github.com/myfantasy/mft"
)

// Errors codes and description
var Errors map[int]string = map[int]string{}

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
