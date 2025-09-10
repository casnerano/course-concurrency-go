package types

type Command string

const (
	CommandUnknown Command = ""
	CommandGet     Command = "GET"
	CommandSet     Command = "SET"
	CommandDel     Command = "DEL"
	CommandClear   Command = "CLEAR"
)

func (c Command) String() string {
	return string(c)
}

func (c Command) Valid() bool {
	switch c {
	case CommandGet, CommandSet, CommandDel, CommandClear:
		return true
	default:
		return false
	}
}
