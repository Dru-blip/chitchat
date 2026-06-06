package ws

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Event int

const (
	EventConnected Event = iota
	EventDisconnected
	EventPing
	EventPong
	EventMessage
)

type Message struct {
	Event   Event  `json:"event"`
	Payload []byte `json:"payload"`
}

type Hub struct {
	lock    sync.RWMutex
	clients map[string]*websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*websocket.Conn),
	}
}

func (h *Hub) Register(deviceID string, conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.clients[deviceID] = conn
}

func (h *Hub) Unregister(deviceID string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	if conn, ok := h.clients[deviceID]; ok {
		conn.Close()
		delete(h.clients, deviceID)
	}
}

func (h *Hub) handleClient(conn *websocket.Conn, deviceId string) {
	defer h.Unregister(deviceId)

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("bad message:", err)
			continue
		}

		switch msg.Event {
		case EventPing:
			{
				conn.WriteJSON(Message{Event: EventPong})
			}
		default:
			{

			}
		}

	}
}

func (h *Hub) SendConnectionEvent(deviceId string) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	conn := h.clients[deviceId]
	conn.WriteJSON(Message{Event: EventConnected})
}
