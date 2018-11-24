package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Setup adds all API routes to given router.
func Setup(router *mux.Router) {
	router.Handle("/test", http.HandlerFunc(handleTest))

	router.HandleFunc("/games", getGames).Methods("GET")
	router.HandleFunc("/games", createGame).Methods("POST")
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
