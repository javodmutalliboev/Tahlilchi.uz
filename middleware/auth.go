package middleware

import (
	"net/http"

	"Tahlilchi.uz/response"
	"github.com/gorilla/sessions"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("#Tahlilchi.uz#-$admin$-?secret?-%key%")
	Store = sessions.NewCookieStore(key)
)

func AdminAuth() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			session, _ := Store.Get(r, "Tahlilchi.uz-admin")

			// Check if user is authenticated
			if auth, ok := session.Values["#Tahlilchi.uz#-$admin$-?authenticated?"].(bool); !ok || !auth {
				response.Res(w, "error", http.StatusForbidden, "Forbidden")
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
