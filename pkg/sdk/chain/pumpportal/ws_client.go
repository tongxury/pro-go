package pumpportal

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	endpoint string
}

func CreateWSClient() *WSClient {
	return &WSClient{
		endpoint: "wss://pumpportal.fun/api/data",
	}
}

type Connection struct {
	*websocket.Conn
}

func (c *WSClient) Connect() (*Connection, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.endpoint, nil)
	if err != nil {
		return nil, err
	}
	return &Connection{conn}, nil
}

func (t *Connection) SubscribeNewToken() error {
	return t.sendPayload("subscribeNewToken")
}

func (t *Connection) SubscribeTokenTrade(mint string) error {
	return t.sendPayload("subscribeTokenTrade", mint)
}

func (t *Connection) UnsubscribeTokenTrade(mint string) error {
	return t.sendPayload("unsubscribeTokenTrade", mint)
}

func (t *Connection) sendPayload(method string, keys ...string) error {
	msg, err := json.Marshal(map[string]interface{}{
		"method": method,
		"keys":   keys,
	})
	if err != nil {
		return err
	}

	return t.WriteMessage(websocket.TextMessage, msg)
}
