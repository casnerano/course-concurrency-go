package types

type Command int

const (
	CommandGet Command = iota
	CommandSet
	CommandDel
)

func (c Command) String() string {
	switch c {
	case CommandGet:
		return "GET"
	case CommandSet:
		return "SET"
	case CommandDel:
		return "DEL"
	default:
		return "UNKNOWN"
	}
}
