package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"Tahlilchi.uz/response"
)

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Missing email parameter", http.StatusBadRequest)
		return
	}

	// Connect to the database
	port, _ := strconv.Atoi(os.Getenv("DBPORT"))
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DBHOST"), port, os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBNAME"))
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT exists (SELECT 1 FROM public.admins WHERE email=$1)", email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	var res response.Response
	if exists {
		eSent := SendEmail("forgot-password")
		if !eSent {
			res = response.Response{
				Status:     "error",
				StatusCode: http.StatusInternalServerError,
				Message:    "unknown event",
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res)
		}

		res = response.Response{
			Status:     "success",
			StatusCode: http.StatusOK,
			Message:    "Code sent to the email",
		}

		w.WriteHeader(http.StatusOK)   // Set the HTTP status code to 200
		json.NewEncoder(w).Encode(res) // Encode the response into JSON
	} else {
		res = response.Response{
			Status:     "error",
			StatusCode: http.StatusNotFound,
			Message:    "Email not found",
		}

		w.WriteHeader(http.StatusNotFound) // Set the HTTP status code to 404
		json.NewEncoder(w).Encode(res)     // Encode the response into JSON
	}
}

func SendEmail(event string) bool {
	if event == "forgot-password" {
		sixDrn := Generate6drn()
		fmt.Println(sixDrn)
		return true
	}
	return false
}

func Generate6drn() int {
	return rand.Intn(900000) + 100000
}
