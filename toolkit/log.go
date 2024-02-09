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

// LogInfo is a function to log an info message inside a route handler function
// it takes *http.Request, string as its parameters
// it returns nothing
func LogInfo(r *http.Request, message string) {
	// use log.Printf to log the message in format "%v: info: %v"
	log.Printf("%v: info: %v", r.URL, message)
}
