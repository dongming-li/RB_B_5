package sockets

import (
	"encoding/json"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda/yerr"
	"github.com/gorilla/websocket"
)

// Client represents a client with a socket connection
type Client struct {
	id      string
	socket  *websocket.Conn
	send    chan interface{}
	manager *Manager
}

// NewClient returns a newly created a client with the id
// and registers it
func NewClient(id string, socket *websocket.Conn, manager *Manager) *Client {
	c := &Client{
		id:      id,
		socket:  socket,
		send:    make(chan interface{}),
		manager: manager,
	}

	manager.register <- c
	return c
}

// Read gets a message from th client and sends it to the manager
func (c *Client) Read() {
	defer func() {
		c.manager.unregister <- c
		c.socket.Close()
	}()

	for {
		// read message from connection
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			c.manager.unregister <- c
			c.socket.Close()
			break
		}

		// if message is a ping, reply a pong
		if string(message) == "ping" {
			c.send <- "pong"
			continue
		}

		// unmarshal message to [Message]
		msg := Message{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			// tell client the message wasn't understood
			c.send <- map[string]interface{}{"code": yerr.InvalidRequest, "message": "the request was not understood"}
			continue
		}

		// inform manager of message
		msg.Sender = c.id
		c.manager.message <- msg

	}
}

// Write sends a message to the client from the manager
func (c *Client) Write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// send pong to client
			if message == "pong" {
				c.socket.WriteMessage(websocket.TextMessage, []byte("pong"))
				continue
			}

			err := c.socket.WriteJSON(message)
			if err != nil {
				c.manager.unregister <- c
				c.socket.Close()
				break
			}
		}

	}

}
