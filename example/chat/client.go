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

//	chat.html 에서 socket.send(msgBox.val());를 통해 서버에 메시지를 전송하면,
//	c.socket.ReadMessage()를 통해 메시지를 수신
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		// 수신한 메시지를 room으로 전송
		c.room.forward <- msg
	}
}

// 	c.send를 통해 수신한 메시지가 있다면,
//	c.socket.WwriteMessage(websocket.TextMessage, msg)를 통해 메시지 전송
//	프론트 쪽 소켓이 메시지를 받으면 socket.onmessage를 통해 사용자에게 추가된 메시지를 보여줌
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
