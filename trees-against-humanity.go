package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rjacobs31/trees-against-humanity-server/auth"
	"github.com/urfave/negroni"
)

var addr = flag.String("addr", ":8000", "http service address")
var aud = flag.String("audience", "", "Auth0 audience")
var iss = flag.String("issuer", "", "Auth0 issuer")
var authDomain = flag.String("auth-domain", "", "Auth0 auhentication domain (defaults to iss)")
var allowedOrigins = flag.String("allowed-origins", "*", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	flag.Parse()
	if *authDomain == "" {
		*authDomain = *iss
	}

	authMiddleware := auth.Auth0Middleware(*aud, *iss, *authDomain)

	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()

	hub := &Hub{}
	hub.Run()
	apiRouter.HandleFunc("/ws", handleWebsocket(hub))

	apiRouter.Handle("/test", negroni.New(
		negroni.HandlerFunc(authMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(handleTest)),
	))

	customHandlers := handlers.CORS(
		handlers.AllowedOrigins(strings.Split(*allowedOrigins, ",")),
		handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Authorization"}),
	)
	http.Handle("/", customHandlers(r))
	log.Println("Starting server at:", *addr)
	err := http.ListenAndServe(*addr, nil)
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
