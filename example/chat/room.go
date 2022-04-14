package chat

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type room struct {
	// 수신 메시지를 보관하는 채널
	// 수신한 메시지는 다른 클라이언트로 전달돼야 함
	forward chan []byte

	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward: //	room으로 전송된 메시지가 있다면, room 내의 모든 클라이언트에게 전송
			for client := range r.clients {
				client.send <- msg // 클라이언트의 write() 메서드 내의 c.send 실행
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}

func NewRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
