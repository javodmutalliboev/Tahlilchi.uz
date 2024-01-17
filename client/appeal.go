package client

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/telegramBot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

	sqlStatement := `INSERT INTO appeals (name, surname, phone_number, message, picture, video) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var id int64
	err = db.QueryRow(sqlStatement, name, surname, phoneNumber, message, picture, video).Scan(&id)
	if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, err.Error())
		return
	}

	response.Res(w, "success", http.StatusCreated, "The appeal form has been submitted successfully.")

	err = sendToTBot(db, id)

	if err != nil {
		fmt.Println(err)
	}
}

func sendToTBot(db *sql.DB, id int64) error {
	// Query the database
	row := db.QueryRow("SELECT name, surname, phone_number, message FROM appeals WHERE id = $1", id)
	rowPV := db.QueryRow("SELECT picture, video FROM appeals WHERE id = $1", id)

	var name, surname, phoneNumber, message string
	var picture, video []byte
	err := row.Scan(&name, &surname, &phoneNumber, &message)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows returned
		} else {
			return err
		}
	}
	err = rowPV.Scan(&picture, &video)
	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows returned
		} else {
			return err
		}
	}

	// Create a new Telegram bot
	bot, err := telegramBot.TBot()
	if err != nil {
		return err
	}

	// Send a message to the Telegram bot
	chatID, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Yangi murojaat keldi | Янги мурожаат келди\nMurojaatchining ismi | Мурожаатчининг исми: %s\nFamiliyasi | Фамилияси: %s\nTelefon raqami | Телефон рақами: %s\nXabar | Хабар: %s", name, surname, phoneNumber, message))

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	// Send the picture to the Telegram Bot, if it exists
	if picture != nil {
		pic := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Name: "picture.jpg", Bytes: picture})
		_, err = bot.Send(pic)
		if err != nil {
			return err
		}
	}

	// Send the video to the Telegram Bot, if it exists
	if video != nil {
		vid := tgbotapi.NewVideoUpload(chatID, tgbotapi.FileBytes{Name: "video.mp4", Bytes: video})
		_, err = bot.Send(vid)
		if err != nil {
			return err
		}
	}

	return nil
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
