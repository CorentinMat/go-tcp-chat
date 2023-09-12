package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	room     map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		room:     make(map[string]*room),
		commands: make(chan command),
	}
}
func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRoom(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}
func (s *server) newClient(conn net.Conn) {
	log.Printf("new client connected : %s\n ", conn.RemoteAddr().String())
	c := &client{
		conn:     conn,
		nick:     "Anonymous",
		commands: s.commands,
	}
	c.readInput()
}
func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	fmt.Printf("Welcome to the server %s\n", c.nick)
}
func (s *server) join(c *client, args []string) {
	roomName := args[1]
	publicKey := args[2]
	if publicKey == "" {
		fmt.Println("Please enter a public key ...")
		return
	}
	fmt.Println("Public Key: ", publicKey)
	r, ok := s.room[roomName]
	if !ok {
		fmt.Println("No room existing, we are creating a new one ...")
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.room[roomName] = r
		fmt.Println("Room created !")
	}
	r.members[c.conn.RemoteAddr()] = c
	s.quitCurrentRoom(c)
	c.room = r
	r.broadcast(c, fmt.Sprintf("%s has joinded the room ! ", c.nick))
	c.msg(fmt.Sprintf("Welcome to %s", r.name))
}

func (s *server) listRoom(c *client, args []string) {
	var rooms []string
	for name := range s.room {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("available rooms: %s", strings.Join(rooms, ", ")))
}
func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("you must provide a room first 🫠"))
		return
	}
	c.room.broadcast(c, c.nick+": "+strings.Join(args[1:], " "))

}
func (s *server) quit(c *client, args []string) {
	log.Printf("%s has left the room", c.conn.RemoteAddr().String())
	s.quitCurrentRoom(c)
	c.msg("see you soon 🫡")
	c.conn.Close()
}
func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room ...👋", c.nick))
	}
}
