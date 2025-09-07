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

var commandOpts = map[Command]struct {
	NeedKeyArg   bool
	NeedValueArg bool
}{
	CommandGet: {NeedKeyArg: true},
	CommandSet: {NeedKeyArg: true, NeedValueArg: true},
	CommandDel: {NeedKeyArg: true},
}

func (c Command) NeedKeyArg() bool {
	return commandOpts[c].NeedKeyArg
}

func (c Command) NeedValueArg() bool {
	return commandOpts[c].NeedValueArg
}
