package middleware

import (
	"net/http"

	"github.com/gorilla/sessions"
)

// Middleware represents a function with the purpose of
// wrapping an http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// MustAuth is a middleware to prevent users from accessing
// certain endpoints until they have authenticated.
func MustAuth(store sessions.Store) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "session-name")
			name, ok := session.Values["username"].(string)
			if !ok || name == "" {
				http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
				return
			}
			f.ServeHTTP(w, r)
		}
	}
}
