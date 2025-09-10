package types

type QueryDel struct {
	keys    []Key
	options Options
}

func (q *QueryDel) Keys() []Key {
	return q.keys
}

func (q *QueryDel) Options() Options {
	return q.options
}

func (q *QueryDel) Command() Command {
	return CommandDel
}

func NewQueryDel(keys []Key, options ...Option) *QueryDel {
	return &QueryDel{
		keys:    keys,
		options: options,
	}
}
