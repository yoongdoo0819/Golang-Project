package chat

import (
	"github.com/gorilla/websocket"
)

type client struct {
	// 클라이언트의 웹 소켓
	socket *websocket.Conn

	// 메시지가 전송되는 채널
	send chan []byte

	// 클라이언트가 채팅하는 방
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
