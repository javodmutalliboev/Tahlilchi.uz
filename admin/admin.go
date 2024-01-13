package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("Tahlilchi.uz-admin-secret-key")
	Store = sessions.NewCookieStore(key)
)

func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "Tahlilchi.uz-admin")

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

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

	// Query the database
	var dbPassword string
	err = db.QueryRow("SELECT password FROM public.admins WHERE email = $1", email).Scan(&dbPassword)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		return
	}

	// Check password
	authenticated := CheckPasswordHash(password, dbPassword)
	if !authenticated {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	// Set admin as authenticated
	session.Options.HttpOnly = true
	session.Options.MaxAge = 3600 * 24
	session.Values["authenticated"] = true
	session.Values["email"] = email
	session.Save(r, w)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginMessage{Message: "Login successful"})
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type LoginMessage struct {
	Message string `json:"message"`
}
