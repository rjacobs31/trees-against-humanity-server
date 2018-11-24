package internal

import (
	"html/template"
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

	r.Handle("/static", http.FileServer(http.Dir("./web/static/")))

	r.HandleFunc("/", rootHandler())

	customHandlers := handlers.CORS(
		handlers.AllowedOrigins(strings.Split(allowedOrigins, ",")),
		handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Authorization"}),
	)
	server := http.Server{
		Addr:    addr,
		Handler: customHandlers(r),
	}

	log.Println("Starting server at:", addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func rootHandler() http.HandlerFunc {
	tmpl, err := template.ParseFiles("./web/template/index.gohtml")
	if err != nil {
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		err := tmpl.Execute(w, struct{}{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
