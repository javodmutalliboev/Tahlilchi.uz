package admin

import (
	"fmt"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
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
		response.Res(w, "error", http.StatusBadRequest, "Failed to parse form")
		return
	}
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Connect to the database
	db, err := db.DB()
	if err != nil {
		response.Res(w, "error", http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	// Query the database
	var dbPassword string
	err = db.QueryRow("SELECT password FROM public.admins WHERE email = $1", email).Scan(&dbPassword)
	if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "Failed to query database")
		return
	}

	// Check password
	authenticated := CheckPasswordHash(password, dbPassword)
	if !authenticated {
		response.Res(w, "error", http.StatusUnauthorized, "Invalid login credentials")
		return
	}

	// Set admin as authenticated
	session.Options.HttpOnly = true
	session.Options.MaxAge = 3600 * 24
	session.Values["#Tahlilchi.uz#-$admin$-?authenticated?"] = true
	session.Values["email"] = email
	session.Save(r, w)
	response.Res(w, "success", http.StatusOK, "Login successful")
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
