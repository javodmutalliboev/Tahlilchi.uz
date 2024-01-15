package admin

import (
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"golang.org/x/crypto/bcrypt"
)

func ForgotPasswordNewPassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	authentication := auth(r)
	if !authentication.status && authentication.message != "" {
		response.Res(w, "error", http.StatusForbidden, authentication.message)
		return
	}

	newPassword := r.Form.Get("new-password")
	if newPassword == "" {
		response.Res(w, "error", http.StatusBadRequest, "new-password not provided")
		return
	}

	hash, _ := HashPassword(newPassword)

	db, err := db.DB()
	if err != nil {
		response.Res(w, "error", http.StatusInternalServerError, "Failed to connect to database")
		return
	}
	defer db.Close()

	// Update statement
	stmt, err := db.Prepare("UPDATE public.admins SET password = $1 WHERE email = $2")
	if err != nil {
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}

	session, _ := Store.Get(r, "admin-forgot-password")
	email := session.Values["email"].(string)
	_, err = stmt.Exec(hash, email)
	if err != nil {
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}

	response.Res(w, "success", http.StatusOK, "New password has been set")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func auth(r *http.Request) authRT {
	session, _ := Store.Get(r, "admin-forgot-password")

	if auth, ok := session.Values["#i#-$code$-?authenticated?"].(bool); !ok || !auth {
		return authRT{status: false, message: "Forbidden"}
	}

	return authRT{
		status:  true,
		message: "",
	}
}

type authRT struct {
	status  bool
	message string
}
