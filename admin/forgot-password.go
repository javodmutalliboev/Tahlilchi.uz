package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"

	"Tahlilchi.uz/response"
	"github.com/gorilla/sessions"
)

func ForgotPasswordEmail(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	if email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
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
		code := Generate6drn()
		eSent := SendEmail("forgot-password", email, code)
		if !eSent {
			res = response.Response{
				Status:     "error",
				StatusCode: http.StatusInternalServerError,
				Message:    "unknown event",
			}

			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res)
		}

		session := saveIdentificationCode(r, code)
		session.Save(r, w)

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

func SendEmail(event string, to string, code int) bool {
	if event == "forgot-password" {
		// Set up authentication information.
		auth := EmailAuth()

		// Connect to the server, authenticate, set the sender and recipient,
		// and send the email all in one step.
		msg := []byte("To: " + to + "\r\n" +
			"Subject: Identification code\r\n" +
			"\r\n" +
			"The code: " + strconv.Itoa(code) + ". Please, do not give this code to anyone. Otherwise, your profile's security will go into risk" + "\r\n")
		err := smtp.SendMail(os.Getenv("SMTPSERVER")+":"+os.Getenv("SMTPPORT"), auth, os.Getenv("EMAILFROM"), []string{to}, msg)
		if err != nil {
			log.Println(err)
			return false
		}

		return true
	}
	return false
}

func Generate6drn() int {
	return rand.Intn(900000) + 100000
}

func EmailAuth() smtp.Auth {
	return smtp.PlainAuth("", os.Getenv("EMAILFROM"), os.Getenv("EMAILFROMPASSWORD"), os.Getenv("SMTPSERVER"))
}

func saveIdentificationCode(r *http.Request, code int) *sessions.Session {
	session, _ := Store.Get(r, "admin-forgot-password")

	session.Options.HttpOnly = true
	session.Options.MaxAge = 3600
	session.Values["identification-code"] = code
	return session
}
