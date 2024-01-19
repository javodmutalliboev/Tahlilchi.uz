package admin

import (
	"fmt"
	"log"
	"net/http"

	"Tahlilchi.uz/authPackage"
	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := authPackage.Store.Get(r, "Tahlilchi.uz-admin")

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
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	// Query the database
	var dbName, dbEmail, dbRole, dbPassword string
	err = db.QueryRow("SELECT name, email, role, password FROM public.admins WHERE email = $1", email).Scan(&dbName, &dbEmail, &dbRole, &dbPassword)
	if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// Check password
	authenticated := checkPasswordHash(password, dbPassword)
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
	response.Res(w, "success", http.StatusOK, admin{Name: dbName, Email: dbEmail, Role: dbRole})
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type admin struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
