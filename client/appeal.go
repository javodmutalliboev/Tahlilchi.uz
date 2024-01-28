package client

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/telegramBot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jackc/pgx/v4"
)

func Appeal(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // Max memory 10MB
	if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
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

	db, err := db.DB()
	if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	sqlStatement := `INSERT INTO appeals (name, surname, phone_number, message) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err = db.QueryRow(sqlStatement, name, surname, phoneNumber, message).Scan(&id)
	if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	connString := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME"))
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer conn.Close(context.Background())

	// For optional fields, check if they are empty
	pictureFile, pictureHeader, err := r.FormFile("picture")
	var pictureOID uint32
	if err == http.ErrMissingFile {

	} else if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
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

		tx, err := conn.Begin(context.Background())
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		defer tx.Rollback(context.Background())

		lob := tx.LargeObjects()
		pictureOID, err = lob.Create(context.Background(), 0)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		obj, err := lob.Open(context.Background(), pictureOID, pgx.LargeObjectModeWrite)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		defer obj.Close()

		_, err = io.Copy(obj, pictureFile)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		_, err = tx.Exec(context.Background(), "UPDATE appeals SET picture = $1 WHERE id = $2", pictureOID, id)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		err = tx.Commit(context.Background())
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		pictureFile.Close()
	}

	videoFile, videoHeader, err := r.FormFile("video")
	var videoOID uint32
	if err == http.ErrMissingFile {

	} else if err != nil {
		fmt.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
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

		tx, err := conn.Begin(context.Background())
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		defer tx.Rollback(context.Background())

		lob := tx.LargeObjects()
		videoOID, err = lob.Create(context.Background(), 0)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		obj, err := lob.Open(context.Background(), videoOID, pgx.LargeObjectModeWrite)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		defer obj.Close()

		_, err = io.Copy(obj, videoFile)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		_, err = tx.Exec(context.Background(), "UPDATE appeals SET video = $1 WHERE id = $2", videoOID, id)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		err = tx.Commit(context.Background())
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		videoFile.Close()
	}

	response.Res(w, "success", http.StatusCreated, "The appeal form has been submitted successfully.")

	err = sendToTBot(db, id)

	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
	}
}

func sendToTBot(db *sql.DB, id int64) error {
	// Query the database
	row := db.QueryRow("SELECT name, surname, phone_number, message, picture, video FROM appeals WHERE id = $1", id)

	connString := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", os.Getenv("DBUSER"), os.Getenv("DBPASSWORD"), os.Getenv("DBHOST"), os.Getenv("DBPORT"), os.Getenv("DBNAME"))
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	var name, surname, phoneNumber, message string
	var picture, video sql.NullInt64
	err = row.Scan(&name, &surname, &phoneNumber, &message, &picture, &video)
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
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Янги мурожаат келди:\nМурожаатчининг исми: %s\nФамилияси: %s\nТелефон рақами: %s\nХабар: %s", name, surname, phoneNumber, message))

	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	const fileSizeLimit int = 50 * 1024 * 1024 // 50MB

	// Send the picture to the Telegram Bot, if it exists
	if picture.Valid {
		tx, err := conn.Begin(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback(context.Background())

		lob := tx.LargeObjects()
		obj, err := lob.Open(context.Background(), uint32(picture.Int64), pgx.LargeObjectModeRead)
		if err != nil {
			return err
		}
		defer obj.Close()

		picBytes, err := io.ReadAll(obj)
		if err != nil {
			return err
		}

		if len(picBytes) > fileSizeLimit {
			message := fmt.Sprintf("[%s %s %s] дан мурожаатда телеграм бот 50МБ ҳажм чегарасидан ошган расм келди. Уни телеграм бот юклай олмайди. Илтимос уни вебсайт администратор панелида коʻринг.", name, surname, phoneNumber)
			msg := tgbotapi.NewMessage(chatID, message)
			_, err = bot.Send(msg)
			if err != nil {
				return err
			}
		} else {
			// Send the picture
			pic := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Name: "picture.jpg", Bytes: picBytes})
			pic.Caption = fmt.Sprintf("Мурожаатчининг исми: %s\nФамилияси: %s\nТелефон рақами: %s", name, surname, phoneNumber)
			_, err = bot.Send(pic)
			if err != nil {
				return err
			}
		}
	}

	// Send the video to the Telegram Bot, if it exists
	if video.Valid {
		tx, err := conn.Begin(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		defer tx.Rollback(context.Background())

		lob := tx.LargeObjects()
		obj, err := lob.Open(context.Background(), uint32(video.Int64), pgx.LargeObjectModeRead)
		if err != nil {
			return err
		}
		defer obj.Close()

		vidBytes, err := io.ReadAll(obj)
		if err != nil {
			return err
		}

		if len(vidBytes) > fileSizeLimit {
			message := fmt.Sprintf("[%s %s %s] дан мурожаатда телеграм бот 50МБ ҳажм чегарасидан ошган видео келди. Уни телеграм бот юклай олмайди. Илтимос уни вебсайт администратор панелида коʻринг.", name, surname, phoneNumber)
			msg := tgbotapi.NewMessage(chatID, message)
			_, err = bot.Send(msg)
			if err != nil {
				return err
			}
		} else {
			// Send the video
			vid := tgbotapi.NewVideoUpload(chatID, tgbotapi.FileBytes{Name: "video.mp4", Bytes: vidBytes})
			vid.Caption = fmt.Sprintf("Мурожаатчининг исми: %s\nФамилияси: %s\nТелефон рақами: %s", name, surname, phoneNumber)
			_, err = bot.Send(vid)
			if err != nil {
				return err
			}
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
