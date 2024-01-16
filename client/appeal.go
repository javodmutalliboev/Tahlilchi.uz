package client

import (
	"fmt"
	"net/http"
	"time"
)

func Appeal(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Max memory 10MB
	if err != nil {
		fmt.Println(err)
		return
	}

	name := r.FormValue("name")
	surname := r.FormValue("surname")
	phoneNumber := r.FormValue("phone_number")
	message := r.FormValue("message")

	// For optional fields, check if they are empty
	picture, _, err := r.FormFile("picture")
	if err != nil && err != http.ErrMissingFile {
		fmt.Println(err)
		return
	}

	video, _, err := r.FormFile("video")
	if err != nil && err != http.ErrMissingFile {
		fmt.Println(err)
		return
	}

	// TODO: Handle the form data
	fmt.Printf("Received form data - Name: %s, Surname: %s, Phone Number: %s, Message: %s\n", name, surname, phoneNumber, message)
	if picture != nil {
		fmt.Println("Received picture")
	}
	if video != nil {
		fmt.Println("Received video")
	}
}

type AppealModel struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	PhoneNumber string    `json:"phone_number"`
	Message     string    `json:"message"`
	Picture     []byte    `json:"picture"`
	Video       []byte    `json:"video"`
	CreatedAt   time.Time `json:"created_at"`
}
