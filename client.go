package main

import "golang.org/x/net/websocket"

type Client struct {
	Username   string
	Connection *websocket.Conn
	Room       *Room
	Server     *Server
}

type ClientList map[*Client]bool

func NewClient(username string, server *Server) *Client {
	return &Client{
		Username: username,
		Server:   server,
	}
}
