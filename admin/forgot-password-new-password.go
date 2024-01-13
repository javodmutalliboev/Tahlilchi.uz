package admin

import (
	"encoding/json"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"golang.org/x/crypto/bcrypt"
)

func ForgotPasswordNewPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		res := response.Response{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	authentication := auth(r)
	if !authentication.status && authentication.message != "" {
		res := response.Response{
			Status:     "error",
			StatusCode: http.StatusForbidden,
			Message:    authentication.message,
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	newPassword := r.Form.Get("new-password")
	if newPassword == "" {
		res := response.Response{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Message:    "new-password not provided",
		}

		json.NewEncoder(w).Encode(res)
	}

	hash, _ := HashPassword(newPassword)

	db, err := db.DB()
	if err != nil {
		// http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		res := response.Response{
			Status:     "error",
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to connect to database",
		}

		json.NewEncoder(w).Encode(res)
		return
	}
	defer db.Close()

	// Update statement
	stmt, err := db.Prepare("UPDATE public.admins SET password = $1 WHERE email = $2")
	if err != nil {
		res := response.Response{
			Status:     "error",
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	session, _ := Store.Get(r, "admin-forgot-password")
	email := session.Values["email"].(string)
	_, err = stmt.Exec(hash, email)
	if err != nil {
		res := response.Response{
			Status:     "error",
			StatusCode: http.StatusInternalServerError,
			Message:    err.Error(),
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	res := response.Response{
		Status:     "success",
		StatusCode: http.StatusOK,
		Message:    "New password has been set",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
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
