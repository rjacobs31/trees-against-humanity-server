package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rjacobs31/trees-against-humanity-server/auth"
	"github.com/urfave/negroni"
)

var addr = flag.String("addr", ":8000", "http service address")
var aud = flag.String("audience", "", "Auth0 audience")
var iss = flag.String("issuer", "", "Auth0 issuer")
var authDomain = flag.String("auth-domain", "", "Auth0 auhentication domain (defaults to iss)")

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

	apiRouter.HandleFunc("/ws", handleWebsocket)
	apiRouter.Handle("/test", negroni.New(
		negroni.HandlerFunc(authMiddleware.HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(handleTest)),
	))

	http.Handle("/", r)
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

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}
