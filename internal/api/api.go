package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rjacobs31/trees-against-humanity-server/internal/middleware"
)

// Setup adds all API routes to given router.
func Setup(router *mux.Router, store sessions.Store) {
	mustAuth := middleware.MustAuth(store)

	router.Handle("/test", http.HandlerFunc(handleTest))

	router.HandleFunc("/games", mustAuth(getGames)).Methods("GET")
	router.HandleFunc("/games", mustAuth(createGame)).Methods("POST")
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"message": "This was a triumph."}`))
}

func getGames(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("[]"))
}

func createGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"error":"Could not create game"}`, http.StatusBadRequest)
}
