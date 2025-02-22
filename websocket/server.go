package websocket

import (
	"net/http"

	"github.com/c0d3-5t3w/GoNetPad/internal/config"
	"github.com/c0d3-5t3w/GoNetPad/internal/logger"
	"github.com/gorilla/websocket"
)

type Server struct {
	clients   map[*websocket.Conn]bool
	broadcast chan string
	upgrader  websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan string),
		upgrader:  websocket.Upgrader{},
	}
}

func (s *Server) Start() {
	http.HandleFunc("/ws", s.handleConnections)
	go s.handleMessages()

	logger.InfoLogger.Printf("Starting WebSocket server on port%s\n", config.WebSocketPort)
	if err := http.ListenAndServe("127.0.0.1"+config.WebSocketPort, nil); err != nil {
		logger.ErrorLogger.Printf("WebSocket server error: %v\n", err)
	}
}

func (s *Server) handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.ErrorLogger.Printf("Error upgrading connection: %v\n", err)
		return
	}
	defer conn.Close()

	s.clients[conn] = true

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			logger.ErrorLogger.Printf("Error reading message: %v\n", err)
			delete(s.clients, conn)
			break
		}
		s.broadcast <- string(msg)
	}
}

func (s *Server) handleMessages() {
	for {
		msg := <-s.broadcast
		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				logger.ErrorLogger.Printf("Error writing message: %v\n", err)
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}
