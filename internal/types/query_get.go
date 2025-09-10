package types

type QueryGet struct {
	key     Key
	options Options
}

func (q *QueryGet) Key() Key {
	return q.key
}

func (q *QueryGet) Options() Options {
	return q.options
}

func (q *QueryGet) Command() Command {
	return CommandGet
}

func NewQueryGet(key Key, options ...Option) *QueryGet {
	return &QueryGet{
		key:     key,
		options: options,
	}
}
