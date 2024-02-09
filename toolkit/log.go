package toolkit

import (
	"log"
	"net/http"
)

// LogError is a function to log an error inside a route handler function
// it takes *http.Request, error as its parameters
// it returns nothing
func LogError(r *http.Request, err error) {
	// use log.Printf to log the error in format "%v: error: %v"
	log.Printf("%v: error: %v", r.URL, err)
}
