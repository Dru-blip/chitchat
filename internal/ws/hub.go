package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type Event int

const (
	EventConnected Event = iota
	EventDisconnected
	EventPing
	EventPong
	EventNewConversation
	EventMessage
	EventError
	EventQueryPresence
	EventPresenceResponse
)

type PresenceQuery struct {
	Users []string `json:"users"`
}

type Message struct {
	Event   Event           `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type Hub struct {
	lock    sync.RWMutex
	devices map[string]*Device
	clients map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		devices: make(map[string]*Device),
		clients: make(map[string]*Client),
	}
}

func (h *Hub) Register(userID, deviceID string, conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()

	if _, exists := h.devices[deviceID]; exists {
		return
	}

	device := &Device{
		Id:     deviceID,
		Conn:   conn,
		userID: userID,
	}

	h.devices[deviceID] = device

	client, ok := h.clients[userID]
	if !ok {
		client = NewClient(userID)
		h.clients[userID] = client
	}
	client.AddDevice(device)
}

func (h *Hub) Unregister(deviceID string) {
	h.lock.Lock()
	defer h.lock.Unlock()

	device, ok := h.devices[deviceID]
	if !ok {
		return
	}

	delete(h.devices, deviceID)
	if client, ok := h.clients[device.userID]; ok {
		client.RemoveDevice(deviceID)
		if client.IsEmpty() {
			delete(h.clients, device.userID)
		}
	}
}

func (h *Hub) GetUserDevices(userID string) ([3]*Device, bool) {
	if client, ok := h.clients[userID]; ok {
		return client.Devices, true
	}
	return [3]*Device{}, false
}

func (h *Hub) handleClient(conn *websocket.Conn, userID, deviceId string) {
	defer h.Unregister(deviceId)

	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			conn.WriteJSON(Message{Event: EventError, Payload: toRawMessage(err.Error())})
			continue
		}

		switch msg.Event {
		case EventPing:
			{
				conn.WriteJSON(Message{Event: EventPong})
			}
		case EventQueryPresence:
			{
				var query PresenceQuery

				if err := json.Unmarshal(msg.Payload, &query); err != nil {
					conn.WriteJSON(Message{Event: EventError, Payload: toRawMessage(err.Error())})
					continue
				}

				onlineUsers := make(map[string]bool, len(query.Users))
				for _, uID := range query.Users {
					if _, ok := h.clients[uID]; ok {
						onlineUsers[uID] = true
					}
				}
				conn.WriteJSON(Message{Event: EventPresenceResponse, Payload: toRawMessage(onlineUsers)})
			}
		default:
			{

			}
		}
	}
}

func (h *Hub) SendEvent(deviceId string, event Event, payload any) {
	device := h.devices[deviceId]
	if device != nil {
		device.Conn.WriteJSON(Message{Event: event, Payload: toRawMessage(payload)})
	}
}

func (h *Hub) SendConnectionEvent(deviceId string) {
	device := h.devices[deviceId]
	if device != nil {
		device.Conn.WriteJSON(Message{Event: EventConnected})
	}
}

func (h *Hub) SendToUser(userID string, event Event, payload any) {
	devices, ok := h.GetUserDevices(userID)
	if !ok {
		return
	}

	for _, device := range devices {
		if device != nil {
			device.Conn.WriteJSON(Message{Event: event, Payload: toRawMessage(payload)})
		}
	}
}

func toRawMessage(payload any) json.RawMessage {
	raw, err := json.Marshal(payload)
	if err != nil {
		return json.RawMessage([]byte(err.Error()))
	}
	return json.RawMessage(raw)
}
