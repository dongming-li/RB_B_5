package sockets

import (
	"log"
)

// Manager is the manger of all socket connections from clients
type Manager struct {
	clients    map[*Client]struct{}
	broadcast  chan Message
	message    chan Message
	register   chan *Client
	unregister chan *Client
	receiver   Receiver

	// started marks when a manager has started i.e.
	// when manager.Start() has been called
	started bool
}

// Receiver can handle messages from socket connection
type Receiver interface {

	// HandleMessage is called whenever the manager receives a message from a client
	// if there is [data] returned, it sends the data to all clients with the id
	HandleMessage(message Message) (recipients []string, data interface{})

	// Format reformats the [message] and returns a [formatted] version
	Format(message interface{}) (formatted interface{})
}

// Message represent the structure of information sent across socket connections
type Message struct {
	Name      string                 `json:"name"`
	Sender    string                 `json:"-"`
	Recipient string                 `json:"recipient,omitempty"`
	Content   map[string]interface{} `json:"content"`
	Type      string                 `json:"type"`
}

// NewManager returns a new socket maneger
func NewManager(r Receiver) *Manager {
	return &Manager{
		broadcast:  make(chan Message),
		message:    make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]struct{}),
		receiver:   r,
	}
}

// registerClient registers a client to the [manager]
func (m *Manager) registerClient(conn *Client) {
	log.Println("Registered socket connection:  ", conn.id)
	m.clients[conn] = struct{}{}
}

// unregisterClient removes a registered client and closes the connection
// with that client if it exists
func (m *Manager) unregisterClient(conn *Client) {
	if _, ok := m.clients[conn]; ok {
		close(conn.send)
		delete(m.clients, conn)
	}
}

// receiveMessage informs the manager of a new message from a client
func (m *Manager) receiveMessage(message Message) {
	to, data := m.receiver.HandleMessage(message)
	if data != nil && to != nil {
		m.SendMessage(data, to)
	}
}

// send sends a message to a connection or removes it if it is inactive
func (m *Manager) send(data interface{}, conn *Client) {
	if _, ok := m.clients[conn]; !ok {
		return
	}
	select {
	case conn.send <- data:
	default:
		close(conn.send)
		delete(m.clients, conn)
	}
}

// broadcastMessage [send]s [data] to all current clients
func (m *Manager) broadcastMessage(data map[string]interface{}) {
	for conn := range m.clients {
		m.send(data, conn)
	}
}

// SendMessage sends the [message] to all client sockets with the id of [id]
func (m *Manager) sendMessage(message interface{}, id string) {
	for conn := range m.clients {
		if conn.id == id {
			m.send(message, conn)
		}
	}
}

// Start initialites a socket listener that listens for connections from clients
// add writes to client connections
func (m *Manager) Start() {
	if m.started {
		return
	}
	m.started = true
	for {
		select {
		// register new client
		case conn := <-m.register:
			m.registerClient(conn)

		// remove client
		case conn := <-m.unregister:
			m.unregisterClient(conn)

		// handle incoming message from a client
		case message := <-m.message:
			m.receiveMessage(message)

		// handle message from cleint to be broadcasted to all clients
		case message := <-m.broadcast:
			m.broadcastMessage(message.Content)

		}
	}
}

// SendMessage sends [message] to all clients with ids in the ids slice
func (m *Manager) SendMessage(message interface{}, ids []string) {
	for i := range ids {
		go m.sendMessage(m.receiver.Format(message), ids[i])
	}
}
