package main

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Token struct {
	Key       string
	CreatedAt time.Time
}

type Client struct {
	Username   string
	Token      *Token
	Connection *websocket.Conn
	Room       *Room
	server     *Server
}

type ClientList map[*Client]bool

func NewToken() *Token {
	return &Token{
		Key:       uuid.NewString(),
		CreatedAt: time.Now(),
	}
}

func NewClient(username string, token *Token, server *Server) *Client {
	return &Client{
		Username: username,
		Token:    token,
		server:   server,
	}
}
