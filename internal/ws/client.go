package ws

import "github.com/gorilla/websocket"

type Device struct {
	Id     string
	userID string
	Conn   *websocket.Conn
}

type Client struct {
	Id      string
	Devices [3]*Device
}

func NewClient(id string) *Client {
	return &Client{
		Id:      id,
		Devices: [3]*Device{nil, nil, nil},
	}
}

func (c *Client) AddDevice(device *Device) {
	for i := range c.Devices {
		if c.Devices[i] == nil {
			c.Devices[i] = device
			return
		}
	}

}

func (c *Client) RemoveDevice(deviceId string) {
	for i := range c.Devices {
		if c.Devices[i] != nil && c.Devices[i].Id == deviceId {
			c.Devices[i] = nil
			return
		}
	}
}

func (c *Client) IsEmpty() bool {
	for _, d := range c.Devices {
		if d != nil {
			return false
		}
	}
	return true
}
