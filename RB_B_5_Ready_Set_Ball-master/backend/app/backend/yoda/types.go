package yoda

import (
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/sockets"
)

// ClientError represents a custom server error sent to the client
type ClientError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Response represents a response from the server to the client
type Response struct {
	Meta   map[string]string `json:"meta,omitempty"`
	Result interface{}       `json:"result,omitempty"`
}

// Request is a client request formatted for Yoda to undertand
type Request struct {
	Meta          map[string]string           `json:"meta"`
	Result        interface{}                 `json:"result,omitempty"`
	Params        map[string]string           `json:"-"`
	Session       map[interface{}]interface{} `json:"-"`
	SocketManager *sockets.Manager             `json:"-"`
}
