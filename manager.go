package main

import (
	"errors"

	"github.com/gorilla/websocket"
	"github.com/rjacobs31/trees-against-humanity-server/game"
)

// Hub controls all of the active games and users
// connected to them.
type Hub struct {
	Users       map[int]User
	clients     map[*Client]User
	Games       map[int]game.Game
	userCounter int
	gameCounter int

	// Registers clients.
	register chan *Client

	// Unregisters cliens.
	unregister chan *Client
}

// NewHub creates a Hub instance to manage clients.
func NewHub() (hub *Hub) {
	return &Hub{
		clients:    make(map[*Client]User),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts up the Hub instance and listens for
// client requests.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = User{Client: client}
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client.send)
		}
	}
}

// AddUser attempts to insert a user into the map of active
// users for the `Manager`.
//
// The username must not already be in use.
func (h *Hub) AddUser(username string, conn *websocket.Conn) (err error) {
	if conn == nil {
		return errors.New("user must have a websocket connection")
	}

	if len(username) < 4 {
		return errors.New("username must be at least 4 characters")
	}

	for _, u := range h.Users {
		if u.Username == username {
			return errors.New("username already taken")
		}
	}

	h.userCounter++
	h.Users[h.userCounter] = User{
		Username: username,
	}

	return nil
}

// RemoveUser attempts to remove a user from the
// collection of active users.
func (h *Hub) RemoveUser(id int) (err error) {
	if id <= 0 {
		return errors.New("must specify an ID above 0")
	}

	_, ok := h.Users[id]

	if !ok {
		return errors.New("user to remove does not exist")
	}

	delete(h.Users, id)

	return nil
}

// AddGame attempts to insert a game into the map of active
// games for the `Manager`.
//
// The game name must not already be in use.
func (h *Hub) AddGame(userID int, name, password string) (err error) {
	_, ok := h.Users[userID]
	if !ok {
		return errors.New("invalid owner ID")
	}

	h.gameCounter++
	h.Games[h.gameCounter] = game.Game{
		ID:       h.gameCounter,
		Name:     name,
		Password: password,
	}

	return nil
}

// RemoveGame attempts to remove a game from the
// collection of active games.
func (h *Hub) RemoveGame(id int) (err error) {
	_, ok := h.Games[id]
	if !ok {
		return errors.New("invalid game ID")
	}

	delete(h.Games, id)

	return nil
}
