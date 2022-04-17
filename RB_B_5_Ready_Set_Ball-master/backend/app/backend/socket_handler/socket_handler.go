package sockethandler

import (
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/sockets"
)

// Subscriber is a function pairing to a subscriber name
// whenever a socket message with the [name] is received,
// the function is called
type Subscriber func(from, to string, data map[string]interface{}) (recepients []string, result interface{})

// SocketHandler wraps the [HandleMessage] for handling incoming messages
type SocketHandler struct {
	subscribers map[string]Subscriber
}

// NewSocketHandler returns a new [SocketHandler]
func NewSocketHandler() *SocketHandler {
	return &SocketHandler{
		subscribers: map[string]Subscriber{
		// "game": game.SocketSubscibe,
		},
	}
}

// HandleMessage handles an incoming message from the socket manager
func (sh *SocketHandler) HandleMessage(message sockets.Message) ([]string, interface{}) {
	to, data := sh.subscribers[message.Name](message.Sender, message.Recipient, message.Content)
	return to, sh.Format(data)
}

// Format converts the message into a yoda response
func (sh *SocketHandler) Format(message interface{}) interface{} {
	var name string
	if m, ok := message.(map[string]interface{}); ok {
		name = m["name"].(string)
	}

	return yoda.NewResponse(map[string]string{"error": "false", "name": name}, message, false)
}
