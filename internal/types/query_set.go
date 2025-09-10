package types

type QuerySet struct {
	key     Key
	value   *Value
	options Options
}

func (q *QuerySet) Key() Key {
	return q.key
}

func (q *QuerySet) Value() *Value {
	return q.value
}

func (q *QuerySet) Options() Options {
	return q.options
}

func (q *QuerySet) Command() Command {
	return CommandSet
}

func NewQuerySet(key Key, value *Value, options ...Option) *QuerySet {
	return &QuerySet{
		key:     key,
		value:   value,
		options: options,
	}
}
