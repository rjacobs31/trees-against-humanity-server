package internal

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Serve(addr, allowedOrigins string) {
	r := mux.NewRouter()

	hub := &Hub{}
	go hub.Run()
	r.HandleFunc("/ws", handleWebsocket(hub))

	apiRouter := r.PathPrefix("/api").Subrouter()

	apiRouter.Handle("/test", http.HandlerFunc(handleTest))

	customHandlers := handlers.CORS(
		handlers.AllowedOrigins(strings.Split(allowedOrigins, ",")),
		handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Authorization"}),
	)
	http.Handle("/", customHandlers(r))

	log.Println("Starting server at:", addr)
	server := http.Server{
		Addr: addr,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"message": "This was a triumph."}`))
}

func handleWebsocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	}
}
