package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
)

type RoomManager struct {
	roomCounter int
	rooms       map[int]room
	store       sessions.Store
}

func (rm *RoomManager) GetRooms() (infos []RoomInfo) {
	infos = make([]RoomInfo, 0, len(rm.rooms))
	for _, v := range rm.rooms {
		infos = append(infos, *v.toRoomInfo())
	}
	return
}

func (rm *RoomManager) CreateRoom(name, password string) (info *RoomInfo, err error) {
	if len(name) < 4 {
		return nil, errors.New("room name too short")
	}

	if rm.rooms == nil {
		rm.rooms = make(map[int]room)
	}

	for _, v := range rm.rooms {
		if v.name == name {
			return nil, errors.New("room name exists")
		}
	}

	rm.roomCounter++
	rm.rooms[rm.roomCounter] = room{
		id:       rm.roomCounter,
		name:     name,
		password: password,
	}

	return rm.rooms[rm.roomCounter].toRoomInfo(), nil
}

func (rm *RoomManager) HandleGetRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rooms := rm.GetRooms()
	body, err := json.Marshal(rooms)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(body)
}

func (rm *RoomManager) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req := createRoomRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info, err := rm.CreateRoom(req.Name, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, _ := json.Marshal(info)
	w.Write(res)
}

type room struct {
	id       int
	name     string
	password string
}

func (r room) toRoomInfo() *RoomInfo {
	return &RoomInfo{ID: r.id, Name: r.name}
}

type RoomInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type createRoomRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
