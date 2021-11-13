package myago

// ContextKey provides a scoped key used to persist data on contexts.
type ContextKey string

func (c ContextKey) String() string {
	return "myago:" + string(c)
}
