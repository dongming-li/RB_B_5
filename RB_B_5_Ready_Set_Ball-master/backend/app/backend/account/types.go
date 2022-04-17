package account

// operation represents an account operation
type operation int

// result is the result of a successful account operation
type result struct {
	meta map[string]string
	data interface{}
}

func (r result) GetMeta() map[string]string {
	return r.meta
}

func (r result) GetData() interface{} {
	return r.data
}
