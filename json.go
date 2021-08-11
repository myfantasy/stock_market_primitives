package smp

import "encoding/json"

type JsonTypedContainer struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
