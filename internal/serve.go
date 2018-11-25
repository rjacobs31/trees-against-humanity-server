package internal

import (
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/boltdb/bolt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rjacobs31/trees-against-humanity-server/internal/api"
	"github.com/yosssi/boltstore/store"
)

const templateDir string = "./web/template/"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeConfig specifies options for configuring
// the Serve command.
type ServeConfig struct {
	Address        string
	AllowedOrigins []string
	SessionSecret  string
}

// Serve initialises a TAH server instance at the
// indicated address and accepting the indicated
// origins
func Serve(config ServeConfig) {
	db, err := bolt.Open("./sessions.db", 0660, nil)
	if err != nil {
		log.Fatal("Open BoltDB: ", err)
	}

	str, err := store.New(db, store.Config{}, []byte(config.SessionSecret))
	if err != nil {
		log.Fatal("Open session store: ", err)
	}

	r, err := mainRouter(str)
	if err != nil {
		log.Fatal("Open router: ", err)
	}

	customHandlers := handlers.CORS(
		handlers.AllowedOrigins(config.AllowedOrigins),
		handlers.AllowedHeaders([]string{"Access-Control-Allow-Headers", "Authorization"}),
	)
	server := http.Server{
		Addr:    config.Address,
		Handler: customHandlers(r),
	}

	log.Println("Starting server at:", config.Address)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func mainRouter(str *store.Store) (r *mux.Router, err error) {
	r = mux.NewRouter()

	apiRouter := r.PathPrefix("/api").Subrouter()
	api.Setup(apiRouter, str)

	hub := &Hub{}
	go hub.Run()
	r.HandleFunc("/ws", handleWebsocket(hub))

	r.Handle("/static", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	r.HandleFunc("/login", loginHandler(str))
	r.HandleFunc("/logout", logoutHandler(str))

	r.HandleFunc("/", rootHandler(str))

	return r, nil
}

func rootHandler(str *store.Store) http.HandlerFunc {
	tmpl, err := parseTemplate("index.gohtml")
	if err != nil {
		log.Fatal("Could not parse root template: ", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := str.Get(r, "session-name")
		name, _ := session.Values["username"].(string)

		w.Header().Set("Content-Type", "text/html")
		data := templateValues{
			Username: name,
		}
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func loginHandler(str *store.Store) http.HandlerFunc {
	tmpl, err := parseTemplate("login.gohtml")
	if err != nil {
		log.Fatal("Could not parse login template: ", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := str.Get(r, "session-name")
		name, ok := session.Values["username"].(string)

		if r.Method == "POST" {
			enteredName := r.FormValue("username")
			if enteredName == "" {
				session.AddFlash("Invalid name")
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}

			session.Values["username"] = enteredName
			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		} else if r.Method == "GET" && ok && name != "" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		data := templateValues{
			Title:    "Login",
			Username: name,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func logoutHandler(str *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := str.Get(r, "session-name")
		name, ok := session.Values["username"].(string)
		if ok && name != "" {
			delete(session.Values, "username")
			session.Save(r, w)
		}
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
}

func handleWebsocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	}
}

func parseTemplate(fileName string) (tmpl *template.Template, err error) {
	filePath := path.Join(templateDir, fileName)
	basePath := path.Join(templateDir, "base.gohtml")
	return tmpl.ParseFiles(basePath, filePath)
}

type templateValues struct {
	Title    string
	Username string
}
