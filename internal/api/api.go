package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rjacobs31/trees-against-humanity-server/internal/middleware"
)

// Setup adds all API routes to given router.
func Setup(router *mux.Router, store sessions.Store) {
	mustAuth := middleware.MustAuth(store)
	rm := RoomManager{store: store}

	router.Handle("/test", http.HandlerFunc(handleTest))

	sesh := sessionHandler{store: store}
	router.HandleFunc("/login", sesh.handleLogin).Headers("Content-Type", "application/json")
	router.HandleFunc("/logout", sesh.handleLogout)

	router.HandleFunc("/games", rm.HandleGetRooms).Methods("GET")
	router.HandleFunc("/games", mustAuth(rm.HandleCreateRoom)).Methods("POST")
}

type sessionHandler struct {
	store sessions.Store
}

type loginRequest struct {
	Username string `json:"username"`
}

func (s *sessionHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	req := loginRequest{}
	body, err := r.GetBody()
	if err != nil {
		http.Error(w, "Could not process request", http.StatusInternalServerError)
		return
	}

	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(body)
	if err != nil {
		http.Error(w, "Could not read body", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		http.Error(w, "Could not deserialise", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username must be non-empty", http.StatusBadRequest)
		return
	}

	session, _ := s.store.Get(r, "session-name")
	session.Values["username"] = req.Username
	session.Save(r, w)

	http.NoBody.WriteTo(w)
}

func (s *sessionHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := s.store.Get(r, "session-name")
	name, ok := session.Values["username"].(string)

	if ok && name != "" {
		delete(session.Values, "username")
		session.Save(r, w)
	}

	http.NoBody.WriteTo(w)
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"message": "This was a triumph."}`))
}
