package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/telegramBot"
	"Tahlilchi.uz/toolkit"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func addAppeal(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1 << 30) // Max memory 1GB
	if err != nil {
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	name := r.FormValue("name")
	if name == "" {
		toolkit.LogError(r, fmt.Errorf("the 'name' field is required"))
		response.Res(w, "error", http.StatusBadRequest, "The 'name' field is required.")
		return
	}

	surname := r.FormValue("surname")
	if surname == "" {
		toolkit.LogError(r, fmt.Errorf("the 'surname' field is required"))
		response.Res(w, "error", http.StatusBadRequest, "The 'surname' field is required.")
		return
	}

	phoneNumber := r.FormValue("phone_number")
	if phoneNumber == "" {
		toolkit.LogError(r, fmt.Errorf("the 'phone_number' field is required"))
		response.Res(w, "error", http.StatusBadRequest, "The 'phone_number' field is required.")
		return
	}

	message := r.FormValue("message")
	if message == "" {
		toolkit.LogError(r, fmt.Errorf("the 'message' field is required"))
		response.Res(w, "error", http.StatusBadRequest, "The 'message' field is required.")
		return
	}

	database, err := db.DB()
	if err != nil {
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// Insert the appeal form into the database returning the id. it is integer
	var id int
	err = database.QueryRow("INSERT INTO appeals (name, surname, phone_number, message) VALUES ($1, $2, $3, $4) RETURNING id", name, surname, phoneNumber, message).Scan(&id)
	if err != nil {
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// insert the picture into the database
	pictureFile, _, err := r.FormFile("picture")
	if err != nil {
		if err != http.ErrMissingFile {
			toolkit.LogError(r, fmt.Errorf("appeal id: %v: error reading the picture file: %v", id, err))
			// delete the appeal by id from the database
			_, err = database.Exec("DELETE FROM appeals WHERE id = $1", id)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("appeal id: %v: error deleting the appeal from the database: %v", id, err))
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	} else {
		defer pictureFile.Close()

		// read the picture file
		picture, err := io.ReadAll(pictureFile)
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("appeal id: %v: error reading the picture file: %v", id, err))
			// delete the appeal by id from the database
			_, err = database.Exec("DELETE FROM appeals WHERE id = $1", id)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("appeal id: %v: error deleting the appeal from the database: %v", id, err))
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		// check picture is an image
		// using http.DetectContentType
		// if it is not an image, return an error
		if http.DetectContentType(picture)[:5] != "image" {
			// log the error
			toolkit.LogError(r, fmt.Errorf("appeal id: %v: the 'picture' field must be an image", id))
			// delete the appeal by id from the database
			_, err = database.Exec("DELETE FROM appeals WHERE id = $1", id)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("appeal id: %v: error deleting the appeal from the database: %v", id, err))
			}
			response.Res(w, "error", http.StatusBadRequest, "The 'picture' field must be an image.")
			return
		}

		// insert the picture into the database by id: picture column is of type bytea
		_, err = database.Exec("UPDATE appeals SET picture = $1 WHERE id = $2", picture, id)
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("appeal id: %v: error inserting the picture into the database: %v", id, err))
			// delete the appeal by id from the database
			_, err = database.Exec("DELETE FROM appeals WHERE id = $1", id)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("appeal id: %v: error deleting the appeal from the database: %v", id, err))
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// insert the video into the database
	videoFile, _, err := r.FormFile("video")
	if err != nil {
		if err != http.ErrMissingFile {
			toolkit.LogError(r, fmt.Errorf("appeal id: %v: error reading the video file: %v", id, err))
			// delete the appeal by id from the database
			_, err = database.Exec("DELETE FROM appeals WHERE id = $1", id)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("appeal id: %v: error deleting the appeal from the database: %v", id, err))
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	} else {
		defer videoFile.Close()

		// read the video file
		video, err := io.ReadAll(videoFile)
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("appeal id: %v: error reading the video file: %v", id, err))
			// delete the appeal by id from the database
			_, err = database.Exec("DELETE FROM appeals WHERE id = $1", id)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("appeal id: %v: error deleting the appeal from the database: %v", id, err))
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		// check video is a video
		// using http.DetectContentType
		// if it is not a video, return an error
		if http.DetectContentType(video)[:5] != "video" {
			// log the error
			toolkit.LogError(r, fmt.Errorf("appeal id: %v: the 'video' field must be a video", id))
			// delete the appeal by id from the database
			_, err = database.Exec("DELETE FROM appeals WHERE id = $1", id)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("appeal id: %v: error deleting the appeal from the database: %v", id, err))
			}
			response.Res(w, "error", http.StatusBadRequest, "The 'video' field must be a video.")
			return
		}

		// insert the video into the database by id: video column is of type bytea
		_, err = database.Exec("UPDATE appeals SET video = $1 WHERE id = $2", video, id)
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("appeal id: %v: error inserting the video into the database: %v", id, err))
			// delete the appeal by id from the database
			_, err = database.Exec("DELETE FROM appeals WHERE id = $1", id)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("appeal id: %v: error deleting the appeal from the database: %v", id, err))
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusCreated, "The appeal form has been submitted successfully.")

	go sendToTBot(r, id)
}

func sendToTBot(r *http.Request, id int) {
	// create a new database connection
	database, err := db.DB()
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("sendToTBot appeal id: %v: error creating a new database connection: %v", id, err))
		return
	}
	// defer the close of the database connection
	defer database.Close()

	// get the appeal by id from the appeals table
	var appeal Appeal
	// select name, surname, phone_number, message
	err = database.QueryRow("SELECT name, surname, phone_number, message FROM appeals WHERE id = $1", id).Scan(&appeal.Name, &appeal.Surname, &appeal.PhoneNumber, &appeal.Message)
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("sendToTBot appeal id: %v: error getting the appeal from the database: %v", id, err))
	}

	// Create a new Telegram bot
	bot, err := telegramBot.TBot()
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("sendToTBot appeal id: %v: error creating a new Telegram bot: %v", id, err))
	}

	// send message to the telegram bot
	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("sendToTBot appeal id: %v: error parsing the chat id: %v", id, err))
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Янги мурожаат келди:\nМурожаатчининг исми: %s\nФамилияси: %s\nТелефон рақами: %s\nХабар: %s", appeal.Name, appeal.Surname, appeal.PhoneNumber, appeal.Message))

	_, err = bot.Send(msg)
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("sendToTBot appeal id: %v: error sending the message to the Telegram bot: %v", id, err))
	}

	const fileSizeLimit int = 50 * 1024 * 1024 // 50MB

	// get the picture of the appeal by id from the database
	err = database.QueryRow("SELECT picture from appeals WHERE id = $1", id).Scan(&appeal.Picture)
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("sendToTBot appeal id: %v: error getting the picture from the database: %v", id, err))
	}

	if len(appeal.Picture) > 0 {
		if len(appeal.Picture) > fileSizeLimit {
			message := fmt.Sprintf("[%s %s %s] дан мурожаатда телеграм бот 50МБ ҳажм чегарасидан ошган расм келди. Уни телеграм бот юклай олмайди. Илтимос уни вебсайт администратор панелида коʻринг.", appeal.Name, appeal.Surname, appeal.PhoneNumber)
			msg := tgbotapi.NewMessage(chatID, message)
			_, err = bot.Send(msg)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("sendToTBot len(appeal.Picture) > fileSizeLimit appeal id: %v: error sending the message to the Telegram bot: %v", id, err))
			}
		} else {
			// Send the picture
			pic := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Name: "picture", Bytes: appeal.Picture})
			pic.Caption = fmt.Sprintf("Мурожаатчининг исми: %s\nФамилияси: %s\nТелефон рақами: %s", appeal.Name, appeal.Surname, appeal.PhoneNumber)
			_, err = bot.Send(pic)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("sendToTBot appeal id: %v: error sending the picture to the Telegram bot: %v", id, err))
			}
		}
	}

	// get the video of the appeal by id from the database
	err = database.QueryRow("SELECT video from appeals WHERE id = $1", id).Scan(&appeal.Video)
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("sendToTBot appeal id: %v: error getting the video from the database: %v", id, err))
	}

	if len(appeal.Video) > 0 {
		if len(appeal.Video) > fileSizeLimit {
			message := fmt.Sprintf("[%s %s %s] дан мурожаатда телеграм бот 50МБ ҳажм чегарасидан ошган видео келди. Уни телеграм бот юклай олмайди. Илтимос уни вебсайт администратор панелида коʻринг.", appeal.Name, appeal.Surname, appeal.PhoneNumber)
			msg := tgbotapi.NewMessage(chatID, message)
			_, err = bot.Send(msg)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("sendToTBot len(appeal.Video) > fileSizeLimit appeal id: %v: error sending the message to the Telegram bot: %v", id, err))
			}
		} else {
			// Send the video
			vid := tgbotapi.NewVideoUpload(chatID, tgbotapi.FileBytes{Name: "video", Bytes: appeal.Video})
			vid.Caption = fmt.Sprintf("Мурожаатчининг исми: %s\nФамилияси: %s\nТелефон рақами: %s", appeal.Name, appeal.Surname, appeal.PhoneNumber)
			_, err = bot.Send(vid)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("sendToTBot appeal id: %v: error sending the video to the Telegram bot: %v", id, err))
			}
		}
	}
}

type Appeal struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	CreatedAt   string `json:"created_at"`
	Picture     []byte `json:"picture"`
	Video       []byte `json:"video"`
}
