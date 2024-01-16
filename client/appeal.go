package client

import (
	"fmt"
	"io"
	"net/http"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
)

func Appeal(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Max memory 10MB
	if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}

	name := r.FormValue("name")
	if name == "" {
		response.Res(w, "error", http.StatusBadRequest, "The 'name' field is required.")
		return
	}

	surname := r.FormValue("surname")
	if surname == "" {
		response.Res(w, "error", http.StatusBadRequest, "The 'surname' field is required.")
		return
	}

	phoneNumber := r.FormValue("phone_number")
	if phoneNumber == "" {
		response.Res(w, "error", http.StatusBadRequest, "The 'phone_number' field is required.")
		return
	}

	message := r.FormValue("message")
	if message == "" {
		response.Res(w, "error", http.StatusBadRequest, "The 'message' field is required.")
		return
	}

	// For optional fields, check if they are empty
	pictureFile, pictureHeader, err := r.FormFile("picture")
	var picture []byte
	if err == http.ErrMissingFile {
		picture = []byte{}
	} else if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	} else {
		// Get the size of the file
		pictureSize := pictureHeader.Size

		// Convert 2MB to bytes
		maxSize := int64(2 << 20)

		if pictureSize > maxSize {
			response.Res(w, "error", http.StatusBadRequest, "The 'picture' field exceeds the 2MB size limit.")
			return
		}

		picture, err = io.ReadAll(pictureFile)
		if err != nil {
			fmt.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, err.Error())
			return
		}
	}

	videoFile, videoHeader, err := r.FormFile("video")
	var video []byte
	if err == http.ErrMissingFile {
		video = []byte{}
	} else if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	} else {
		// Get the size of the file
		videoSize := videoHeader.Size

		// Convert 6MB to bytes
		maxSize := int64(6 << 20)

		if videoSize > maxSize {
			response.Res(w, "success", http.StatusBadRequest, "The 'video' field exceeds the 6MB size limit.")
			return
		}

		video, err = io.ReadAll(videoFile)
		if err != nil {
			fmt.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, err.Error())
			return
		}
	}

	db, err := db.DB()
	if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}
	defer db.Close()

	sqlStatement := `INSERT INTO appeals (name, surname, phone_number, message, picture, video) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = db.Exec(sqlStatement, name, surname, phoneNumber, message, picture, video)
	if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}

	response.Res(w, "success", http.StatusCreated, "The appeal form has been submitted successfully.")
}

/*
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
*/
