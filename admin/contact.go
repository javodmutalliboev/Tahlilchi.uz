package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

type Appeal struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	CreatedAt   string `json:"created_at"`
	Picture     *int   `json:"picture"`
	Video       *int   `json:"video"`
}

func appealList(w http.ResponseWriter, r *http.Request) {
	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, name, surname, phone_number, message, created_at, picture, video FROM appeals")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var appeals []Appeal
	for rows.Next() {
		var a Appeal
		if err := rows.Scan(&a.ID, &a.Name, &a.Surname, &a.PhoneNumber, &a.Message, &a.CreatedAt, &a.Picture, &a.Video); err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		appeals = append(appeals, a)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, appeals)
}

func appealExists(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM appeals WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

func appealPicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idStr := vars["id"]

	exists, err := appealExists(idStr)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: appealExists(idStr): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusNotFound, "appeal not found")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var oid sql.NullInt64
	err = database.QueryRow("SELECT picture FROM appeals WHERE id = $1", idStr).Scan(&oid)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	conn, err := pgx.Connect(context.Background(), db.ConnString())
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer tx.Rollback(context.Background())

	lob := tx.LargeObjects()
	obj, err := lob.Open(context.Background(), uint32(oid.Int64), pgx.LargeObjectModeRead)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusNotFound, "File not found")
		return
	}
	defer obj.Close()

	_, err = io.Copy(w, obj)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
}

func appealVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idStr := vars["id"]

	exists, err := appealExists(idStr)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: appealExists(idStr): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusNotFound, "appeal not found")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var oid sql.NullInt64
	err = database.QueryRow("SELECT video FROM appeals WHERE id = $1", idStr).Scan(&oid)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	conn, err := pgx.Connect(context.Background(), db.ConnString())
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer tx.Rollback(context.Background())

	lob := tx.LargeObjects()
	obj, err := lob.Open(context.Background(), uint32(oid.Int64), pgx.LargeObjectModeRead)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusNotFound, "File not found")
		return
	}
	defer obj.Close()

	_, err = io.Copy(w, obj)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
}

func adminContactExists() (*bool, error) {
	database, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer database.Close()

	var exists bool
	err = database.QueryRow("SELECT EXISTS (SELECT 1 FROM admin_contact LIMIT 1)").Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

type AdminContact struct {
	Address     string   `json:"address"`
	SocMedAcs   []string `json:"soc_med_acs"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
}

func createAdminContact(w http.ResponseWriter, r *http.Request) {
	exists, err := adminContactExists()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *exists {
		response.Res(w, "error", http.StatusBadRequest, "admin contact already created")
		return
	}

	var adminContact AdminContact
	err = json.NewDecoder(r.Body).Decode(&adminContact)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	if adminContact.Address == "" || len(adminContact.SocMedAcs) == 0 || adminContact.PhoneNumber == "" || adminContact.Email == "" {
		response.Res(w, "error", http.StatusBadRequest, "address, soc_med_acs, phone_number, email are all required")
		return
	}

	// Check if any element in SocMedAcs is empty
	for _, socMedAc := range adminContact.SocMedAcs {
		if socMedAc == "" {
			response.Res(w, "error", http.StatusBadRequest, "Social media account cannot be empty")
			return
		}
	}

	socMedAcsString := toolkit.SliceToString(adminContact.SocMedAcs)

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	_, err = database.Exec("INSERT INTO admin_contact (address, soc_med_acs, phone_number, email) VALUES ($1, $2, $3, $4)",
		adminContact.Address, socMedAcsString, adminContact.PhoneNumber, adminContact.Email)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "admin contact data created")
}

type AdminContactGET struct {
	ID          int      `json:"id"`
	Address     string   `json:"address"`
	SocMedAcs   []string `json:"soc_med_acs"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}

func getAdminContact(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	result, err := database.Query("SELECT * FROM admin_contact LIMIT 1")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer result.Close()

	var contact AdminContactGET
	for result.Next() {
		var socMedAcs []uint8
		err := result.Scan(&contact.ID, &contact.Address, &socMedAcs, &contact.PhoneNumber, &contact.Email, &contact.CreatedAt, &contact.UpdatedAt)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		socMedAcsString := strings.Trim(string(socMedAcs), "{}")
		contact.SocMedAcs = strings.Split(string(socMedAcsString), ",")
	}

	response.Res(w, "success", http.StatusOK, contact)
}

func updateAdminContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var contact AdminContact
	_ = json.NewDecoder(r.Body).Decode(&contact)

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	if contact.Address != "" {
		_, err := database.Exec("UPDATE admin_contact SET address = $1, updated_at = NOW() WHERE id = $2", contact.Address, params["id"])
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if len(contact.SocMedAcs) > 0 {
		var socMedAcs []string
		for _, socMed := range contact.SocMedAcs {
			if socMed != "" {
				socMedAcs = append(socMedAcs, socMed)
			}
		}
		socMedAcsString := toolkit.SliceToString(socMedAcs)

		_, err := database.Exec("UPDATE admin_contact SET soc_med_acs = $1, updated_at = NOW() WHERE id = $2", socMedAcsString, params["id"])
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if contact.PhoneNumber != "" {
		_, err := database.Exec("UPDATE admin_contact SET phone_number = $1, updated_at = NOW() WHERE id = $2", contact.PhoneNumber, params["id"])
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if contact.Email != "" {
		_, err := database.Exec("UPDATE admin_contact SET email = $1, updated_at = NOW() WHERE id = $2", contact.Email, params["id"])
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusCreated, "admin contact updated")
}

type AppealCount struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

func getAppealCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	period := vars["period"]

	if period != "day" && period != "week" && period != "month" && period != "year" {
		response.Res(w, "error", http.StatusBadRequest, "invalid period value")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM appeals WHERE created_at > current_date - interval '1 %s'", period)
	err = database.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, AppealCount{Period: period, Count: count})
}
