package main

type Room struct {
	Client1 *Client
	Client2 *Client
}

type RoomList map[*Room]bool

func NewRoom(s *Server) *Room {
	room := &Room{}
	s.rooms[room] = true
	return room
}

func (r *Room) available() bool {
	if r.Client1 == nil || r.Client2 == nil {
		return true
	}
	return false
}
