package main

import (
	"encoding/json"
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
					fmt.Println("Joined to a room")
					if c.Room.Client1 == nil {
						c.Room.Client1 = c
					} else {
						c.Room.Client2 = c
					}
					client.writeJSON("INFO", "User "+c.Username+" Joined")
					client.writeJSON("CONN", c.Username)
					c.writeJSON("INFO", "User "+client.Username+" Joined")
					c.writeJSON("CONN", client.Username)
				}
			}
		}
	}
	if !roomAllocated {
		room := NewRoom(c.server)
		c.Room = room
		c.Room.Client1 = c
		fmt.Println("New Room created")
	}
	return nil
}

func (c *Client) Send(msg []byte) {
	if c.Room.Client1 == c {
		if c.Room.Client2 != nil {
			c.Room.Client2.writeJSON("MSG", string(msg))
			fmt.Println("MESSAGE - (" + string(msg) + ") from - " + c.Username + " to - " + c.Room.Client2.Username)
		} else {
			c.writeJSON("INFO", "No User Connected")
		}
	} else {
		if c.Room.Client1 != nil {
			c.Room.Client1.writeJSON("MSG", string(msg))
			fmt.Println("MESSAGE - (" + string(msg) + ") from - " + c.Username + " to - " + c.Room.Client1.Username)
		} else {
			c.writeJSON("INFO", "No User Connected")
		}
	}
}

func (c *Client) disconnect(s *Server) {
	fmt.Println("Disconnected  - ", c.Username)
	c.Connection = nil
	if c.Room.Client1 == c {
		c.Room.Client1 = nil
		if c.Room.Client2 != nil {
			c.Room.Client2.writeJSON("INFO", c.Username+" left")
			c.Room.Client2.writeJSON("DISCONN", "")
		}
	} else {
		c.Room.Client2 = nil
		if c.Room.Client1 != nil {
			c.Room.Client1.writeJSON("INFO", c.Username+" left")
			c.Room.Client1.writeJSON("DISCONN", "")
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

func (c *Client) writeJSON(_type string, _data string) {
	obj := map[string]interface{}{
		"type": _type,
		"data": _data,
	}
	jsonBytes, _ := json.Marshal(obj)
	c.Connection.Write(jsonBytes)
}
