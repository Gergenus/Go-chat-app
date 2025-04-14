package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Socket   *websocket.Conn
	Message  chan *Message
	ID       string
	RoomID   string
	Username string
}

type Message struct {
	Content  string
	RoomID   string
	Username string
}

func (c *Client) WriteMessage() {
	defer c.Socket.Close()
	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		c.Socket.WriteJSON(message)
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, m, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("error", err)
			}
			break
		}
		msg := &Message{
			Content:  string(m),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		hub.Broadcast <- msg
	}

}
