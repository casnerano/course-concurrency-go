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

func IntValue(value int) Value {
	return Value{
		Type: ValueTypeInt,
		Data: value,
	}
}

func StringValue(value string) Value {
	return Value{
		Type: ValueTypeString,
		Data: value,
	}
}

func BoolValue(value bool) Value {
	return Value{
		Type: ValueTypeBool,
		Data: value,
	}
}
