package admin

import (
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/middleware"
	"Tahlilchi.uz/response"
	"github.com/gorilla/sessions"
)

func forgotPasswordEmail(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	email := r.Form.Get("email")
	if email == "" {
		response.Res(w, "error", http.StatusBadRequest, "email is required")
		return
	}

	// Connect to the database
	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT exists (SELECT 1 FROM public.admins WHERE email=$1)", email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if exists {
		iCode := generate6drn()
		eSent := sendEmail("forgot-password", email, iCode)
		if !eSent.Status {
			response.Res(w, "error", http.StatusInternalServerError, eSent.Message)
			return
		}

		session := saveIdentificationCode(r, email, iCode)
		session.Save(r, w)

		response.Res(w, "success", http.StatusOK, "Code sent to the email")
	} else {
		response.Res(w, "error", http.StatusNotFound, "Email not found")
	}
}

type emailStatus struct {
	Status  bool
	Message string
}

func sendEmail(event string, to string, code int) emailStatus {
	if event == "forgot-password" {
		// Set up authentication information.
		auth := emailAuth()

		// Connect to the server, authenticate, set the sender and recipient,
		// and send the email all in one step.
		msg := []byte("To: " + to + "\r\n" +
			"Subject: Tahlilchi.uz administratori parolini unutgan vaqt uchun administrator identifikatsiya kodi | Tahlilchi.uz администратори паролини унутган вақт учун администратор идентификация коди\r\n" +
			"\r\n" +
			"Kod | Код: " + strconv.Itoa(code) + ". Iltimos, bu kodni hech kimga bermang. Aks holda, profilingiz xavfsizligi xavf ostida qoladi. | Илтимос, бу кодни ҳеч кимга берманг. Акс ҳолда, профилингиз хавфсизлиги хавф остида қолади." + "\r\n")
		err := smtp.SendMail(os.Getenv("SMTPSERVER")+":"+os.Getenv("SMTPPORT"), auth, os.Getenv("EMAILFROM"), []string{to}, msg)
		if err != nil {
			log.Println(err)
			return emailStatus{Status: false, Message: err.Error()}
		}

		return emailStatus{Status: true, Message: ""}
	}
	return emailStatus{Status: false, Message: "unknown event"}
}

func generate6drn() int {
	return rand.Intn(900000) + 100000
}

func emailAuth() smtp.Auth {
	return smtp.PlainAuth("", os.Getenv("EMAILFROM"), os.Getenv("EMAILFROMPASSWORD"), os.Getenv("SMTPSERVER"))
}

func saveIdentificationCode(r *http.Request, email string, iCode int) *sessions.Session {
	session, _ := middleware.Store.Get(r, "admin-forgot-password")

	session.Options.HttpOnly = true
	session.Options.MaxAge = 3600
	session.Values["email"] = email
	session.Values["i-code"] = iCode
	return session
}
