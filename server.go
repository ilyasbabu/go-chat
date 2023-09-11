package main

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

type Server struct {
	clients ClientList
	rooms   RoomList
}

func NewServer() *Server {
	return &Server{
		clients: make(ClientList),
		rooms:   make(RoomList),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new Connection from - ", ws.RemoteAddr())

	ws.Write([]byte("Hello from server"))
	buffer := make([]byte, 1024)
	for {
		n, err := ws.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Disconnected ")
				break
			}
			fmt.Println(err)
			continue
		}
		msg := buffer[:n]
		fmt.Println("msg recieved in server - ", string(msg))
	}
}
