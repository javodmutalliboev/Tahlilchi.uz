package client

import (
	"log"
	"net/http"
	"strings"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
)

type AdminContactGET struct {
	ID          int      `json:"id"`
	Address     string   `json:"address"`
	SocMedAcs   []string `json:"soc_med_acs"`
	PhoneNumber string   `json:"phone_number"`
	Email       string   `json:"email"`
}

func getAdminContact(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	result, err := database.Query("SELECT id, address, soc_med_acs, phone_number, email FROM admin_contact LIMIT 1")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer result.Close()

	var contact AdminContactGET
	for result.Next() {
		var socMedAcs []uint8
		err := result.Scan(&contact.ID, &contact.Address, &socMedAcs, &contact.PhoneNumber, &contact.Email)
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
