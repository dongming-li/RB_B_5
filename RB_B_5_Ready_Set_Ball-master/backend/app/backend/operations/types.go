package operations

// result is the result of an operation
type result struct {
	meta map[string]string
	data interface{}
}

func (r result) GetMeta() map[string]string { return r.meta }

func (r result) GetData() interface{} { return r.data }

// owner represents the logged in user
type owner struct {
	id, name string
}

// ID returns the userid of the logged in user
func (o owner) ID() string { return o.id }

// Username returns the usernam of the logged in user
func (o owner) Username() string { return o.name }
