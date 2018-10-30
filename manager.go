package main

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/rjacobs31/trees-against-humanity-server/game"
)

// Manager controls all of the active games and users
// connected to them.
type Manager struct {
	Users       map[int]User
	Games       map[int]game.Game
	userCounter int
	gameCounter int
}

// User represents a single user and their connection
// to the server.
//
// Currently, the user should die with the connection.
type User struct {
	ID         int
	Connection *websocket.Conn
	Username   string
}

// AddUser attempts to insert a user into the map of active
// users for the `Manager`.
//
// The username must not already be in use.
func (m *Manager) AddUser(username string, conn *websocket.Conn) (err error) {
	if conn == nil {
		return errors.New("user must have a websocket connection")
	}

	if len(username) < 4 {
		return errors.New("username must be at least 4 characters")
	}

	for _, u := range m.Users {
		if u.Username == username {
			return errors.New("username already taken")
		}
	}

	m.userCounter++
	m.Users[m.userCounter] = User{
		ID:         m.userCounter,
		Connection: conn,
		Username:   username,
	}

	return nil
}

// RemoveUser attempts to remove a user from the
// collection of active users.
func (m *Manager) RemoveUser(id int) (err error) {
	if id <= 0 {
		return errors.New("must specify an ID above 0")
	}

	user, ok := m.Users[id]

	if !ok {
		return errors.New("user to remove does not exist")
	}

	delete(m.Users, id)
	user.Connection.Close()

	return nil
}

// AddGame attempts to insert a game into the map of active
// games for the `Manager`.
//
// The game name must not already be in use.
func (m *Manager) AddGame(userID int, name, password string) (err error) {
	_, ok := m.Users[userID]
	if !ok {
		return errors.New("invalid owner ID")
	}

	m.gameCounter++
	m.Games[m.gameCounter] = game.Game{
		ID:       m.gameCounter,
		Name:     name,
		Password: password,
	}

	return nil
}

// RemoveGame attempts to remove a game from the
// collection of active games.
func (m *Manager) RemoveGame(id int) (err error) {
	_, ok := m.Games[id]
	if !ok {
		return errors.New("invalid game ID")
	}

	delete(m.Games, id)

	return nil
}
