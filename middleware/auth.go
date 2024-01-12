package middleware

import (
	"net/http"

	"Tahlilchi.uz/admin"
)

func Auth() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			session, _ := admin.Store.Get(r, "admin")

			// Check if user is authenticated
			if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}
