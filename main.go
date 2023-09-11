package main

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	server := NewServer()
	go server.StatusLoggerListener()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.HandleFunc("/login", server.handleLogin)
	http.ListenAndServe(":8080", nil)
}
