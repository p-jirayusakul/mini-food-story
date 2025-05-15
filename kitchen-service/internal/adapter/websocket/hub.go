package websocket

import (
	"github.com/gofiber/contrib/websocket"
	"log"
)

type Hub struct {
	Clients    map[*websocket.Conn]bool // เก็บ client ทั้งหมด
	Broadcast  chan []byte              // ช่องทางรับข้อความที่ต้องส่งให้ทุกคน
	Register   chan *websocket.Conn     // มี client ใหม่เข้ามา
	Unregister chan *websocket.Conn     // มี client ออกไป
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.Clients[conn] = true // เก็บ connection ใหม่
		case conn := <-h.Unregister:
			delete(h.Clients, conn) // ลบตอน disconnect
			err := conn.Close()
			if err != nil {
				log.Fatal(err)
				return
			}
		case message := <-h.Broadcast:
			for conn := range h.Clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Error sending message:", err)
					return
				} // ส่งข้อความให้ทุกคน
			}
		}
	}
}

func (h *Hub) Shutdown() {
	for conn := range h.Clients {
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server shutting down"))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
		err = conn.Close()
		if err != nil {
			log.Println("Error closing connection:", err)
			return
		}
	}
}
