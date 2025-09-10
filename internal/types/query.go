package types

type Query interface {
	Options() Options
	Command() Command
}
