package websocket

import (
	"log"

	"github.com/gofiber/contrib/websocket"
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
			h.handleRegister(conn)
		case conn := <-h.Unregister:
			h.handleUnregister(conn)
		case message := <-h.Broadcast:
			h.handleBroadcast(message)
		}
	}
}

func (h *Hub) handleRegister(conn *websocket.Conn) {
	if conn != nil {
		h.Clients[conn] = true
	}
}

func (h *Hub) handleUnregister(conn *websocket.Conn) {
	if conn != nil {
		if _, ok := h.Clients[conn]; ok {
			delete(h.Clients, conn)
			_ = conn.Close() // ปิด connection
		}
	}
}

func (h *Hub) handleBroadcast(message []byte) {
	for conn := range h.Clients {
		if conn == nil {
			continue
		}
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			h.Unregister <- conn // ถ้าเขียนไม่ได้ แสดงว่า connection ตายแล้ว
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
