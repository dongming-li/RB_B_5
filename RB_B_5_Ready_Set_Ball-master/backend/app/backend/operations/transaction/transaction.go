package transaction

import (
	"fmt"
	"strconv"
)

// Result is the result of an transaction
// It is implememted by [operation.Result]
type Result interface {
	GetMeta() map[string]string
	GetData() interface{}
}

// Owner represents the identification information of the
// user/account that starts a transaction
type Owner interface {
	// ID returns the ID of the user
	ID() string

	// Username returns the username of the owner
	Username() string
}

// ReadableMap creates a [map[string]string] from a [map[string]interface{}] and returns an error if a type assertion fails
func ReadableMap(m interface{}) (map[string]string, error) {
	m1, _ := m.(map[string]interface{})
	m2 := make(map[string]string, len(m1))

	for key, value := range m1 {
		switch value := value.(type) {
		case string:
			m2[key] = value
		case bool:
			m2[key] = strconv.FormatBool(value)
		default:
			return nil, fmt.Errorf("Expected param type to be string/bool but recieved map[%v]=%v as %T", key, value, value)
		}
	}

	return m2, nil
}
