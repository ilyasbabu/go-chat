package main

import (
	"net/http"

	"github.com/ilyasbabu/go-chat/app"
	"golang.org/x/net/websocket"
)

func main() {
	server := app.NewServer()
	go server.StatusLoggerListener()
	http.Handle("/ws", websocket.Handler(server.HandleWS))
	http.HandleFunc("/login", server.HandleLogin)
	http.ListenAndServe(":8080", nil)
}
