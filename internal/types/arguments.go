package types

import "fmt"

type Key string

type ValueType int

const (
	ValueTypeInt ValueType = iota
	ValueTypeString
	ValueTypeBool
)

type Value struct {
	Type ValueType
	Data any
}

func (v Value) String() string {
	return fmt.Sprintf("%v", v.Data)
}

type Option string
type Options []Option
