package main

import (
	"errors"
	"fmt"
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

func (c *Client) SetRoom() error {
	var roomAllocated bool
	if c.Room != nil {
		return errors.New("already in a Room")
	}
	for client := range c.server.clients {
		if client != c {
			if client.Room != nil {
				if client.Room.available() {
					c.Room = client.Room
					roomAllocated = true
					c.Connection.Write([]byte("Joined to a Room"))
					if c.Room.Client1 == nil {
						c.Room.Client1 = c
					} else {
						c.Room.Client2 = c
					}
					client.Connection.Write([]byte("User " + c.Username + " Joined"))
					c.Connection.Write([]byte("User " + client.Username + " Joined"))
				}
			}
		}
	}
	if !roomAllocated {
		room := NewRoom(c.server)
		c.Room = room
		c.Room.Client1 = c
		c.Connection.Write([]byte("New Room created"))
	}
	return nil
}

func (c *Client) Send(msg []byte) {
	if c.Room.Client1 == c {
		if c.Room.Client2 != nil {
			c.Room.Client2.Connection.Write(msg)
		} else {
			c.Connection.Write([]byte("No User Connected"))
		}
	} else {
		if c.Room.Client1 != nil {
			c.Room.Client1.Connection.Write(msg)
		} else {
			c.Connection.Write([]byte("No User Connected"))
		}
	}
}

func (c *Client) disconnect(s *Server) {
	fmt.Println("Disconnected  - ", c.Username)
	c.Connection = nil
	if c.Room.Client1 == c {
		c.Room.Client1 = nil
		if c.Room.Client2 != nil {
			c.Room.Client2.Connection.Write([]byte(c.Username + " left"))
		}
	} else {
		c.Room.Client2 = nil
		if c.Room.Client1 != nil {
			c.Room.Client1.Connection.Write([]byte(c.Username + " left"))
		}
	}
	if c.Room.Client1 == nil && c.Room.Client2 == nil {
		for r := range s.rooms {
			if r == c.Room {
				delete(s.rooms, c.Room)
			}
		}
	}
	c.Room = nil
}
