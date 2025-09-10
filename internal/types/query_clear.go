package types

type QueryClear struct {
	options Options
}

func (q *QueryClear) Options() Options {
	return q.options
}

func (q *QueryClear) Command() Command {
	return CommandClear
}

func NewQueryClear(options ...Option) *QueryClear {
	return &QueryClear{
		options: options,
	}
}
