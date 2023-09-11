package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

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

func (s *Server) authenticateToken(ws *websocket.Conn) (*Client, error) {
	if token, ok := ws.Request().URL.Query()["token"]; ok {
		ok = false
		for client := range s.clients {
			if client.Token.Key == token[0] {
				client.Connection = ws
				ok = true
				return client, nil
			}
		}
		if !ok {
			return nil, errors.New("invalid token")
		}
	}
	return nil, errors.New("no token provided")
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new Connection from - ", ws.RemoteAddr())
	client, err := s.authenticateToken(ws)
	if err != nil {
		ws.Write([]byte(err.Error()))
		return
	}
	err = client.SetRoom()
	if err != nil {
		ws.Write([]byte(err.Error()))
		return
	}
	// implement ping pong
	buffer := make([]byte, 1024)
	for {
		n, err := ws.Read(buffer)
		if err != nil {
			if err == io.EOF {
				client.disconnect(s)
				break
			}
			fmt.Println(err)
			continue
		}
		msg := buffer[:n]
		fmt.Println("msg recieved in server - ", string(msg))
		client.Send(msg)
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization,access-control-allow-methods,access-control-allow-origin,access-control-allow-headers")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	type userLoginRequest struct {
		Username string `json:"username"`
	}
	var req userLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token := NewToken()
	client := NewClient(req.Username, token, s)
	s.clients[client] = true
	http.Error(w, token.Key, http.StatusOK)
}

func (s *Server) StatusLoggerListener() {
	var inp string
	for {
		fmt.Scanln(&inp)
		if inp == "s" {
			var activeWScount int
			for client := range s.clients {
				if client.Connection != nil {
					activeWScount++
				}
			}
			fmt.Println("-----------Server status-----------")
			fmt.Println(" clients count in server - ", len(s.clients))
			fmt.Println(" active websocket count - ", activeWScount)
			fmt.Println(" active rooms count - ", len(s.rooms))
			fmt.Println("-----------------------------------")
		}
		time.Sleep(time.Second * 1)
	}
}
