package client

import (
	"log"
	"net/http"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
)

func getENewspaperList(w http.ResponseWriter, r *http.Request) {
	// Parse the page number from the query parameters
	pageStr, ok := r.URL.Query()["page"]
	if !ok || len(pageStr[0]) < 1 {
		log.Printf("%v: Url Param 'page' is missing. Setting default value to 1.", r.URL)
		pageStr = []string{"1"}
	}
	page, _ := strconv.Atoi(pageStr[0])

	// Parse the limit from the query parameters
	limitStr, ok := r.URL.Query()["limit"]
	if !ok || len(limitStr[0]) < 1 {
		log.Printf("%v: Url Param 'limit' is missing. Setting default value to 10.", r.URL)
		limitStr = []string{"10"}
	}
	limit, _ := strconv.Atoi(limitStr[0])

	// Calculate the starting index
	var start = (page - 1) * limit

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title_latin, title_cyrillic FROM e_newspapers WHERE archived = false AND completed = true ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var eNewspapers []ENewspaper
	for rows.Next() {
		var eNewspaper ENewspaper
		err := rows.Scan(&eNewspaper.ID, &eNewspaper.TitleLatin, &eNewspaper.TitleCyrillic)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		eNewspapers = append(eNewspapers, eNewspaper)
	}

	response.Res(w, "success", http.StatusOK, eNewspapers)
}

type ENewspaper struct {
	ID            int    `json:"id"`
	TitleLatin    string `json:"title_latin"`
	TitleCyrillic string `json:"title_cyrillic"`
}

func eNewspaperExists(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM e_newspapers WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

type ENewspaperByID struct {
	ID            int    `json:"id"`
	TitleLatin    string `json:"title_latin"`
	TitleCyrillic string `json:"title_cyrillic"`
	FileLatin     []byte `json:"file_latin"`
	FileCyrillic  []byte `json:"file_cyrillic"`
}

// getENewspaperFile is a handler function that by id and alphabet returns file_latin or file_cyrillic as pdf with the name of title_latin or tile_cyrillic
func getENewspaperFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	idStr := vars["id"]
	alphabet := vars["alphabet"]

	if alphabet != "latin" && alphabet != "cyrillic" {
		response.Res(w, "error", http.StatusBadRequest, "invalid alphabet value")
		return
	}

	exists, err := eNewspaperExists(idStr)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: eNewspaperExists(idStr): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusNotFound, "e-newspaper not found")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var eNewspaper ENewspaperByID
	if alphabet == "latin" {
		err = database.QueryRow("SELECT id, title_latin, file_latin FROM e_newspapers WHERE id = $1 AND archived = FALSE AND completed = TRUE", idStr).Scan(&eNewspaper.ID, &eNewspaper.TitleLatin, &eNewspaper.FileLatin)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
	}

	if alphabet == "cyrillic" {
		err = database.QueryRow("SELECT id, title_cyrillic, file_cyrillic FROM e_newspapers WHERE id = $1 AND archived = FALSE AND completed = TRUE", idStr).Scan(&eNewspaper.ID, &eNewspaper.TitleCyrillic, &eNewspaper.FileCyrillic)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
	}

	var file []byte
	var fileName string
	if alphabet == "latin" {
		file = eNewspaper.FileLatin
		fileName = eNewspaper.TitleLatin
	}

	if alphabet == "cyrillic" {
		file = eNewspaper.FileCyrillic
		fileName = eNewspaper.TitleCyrillic
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName+".pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.Write(file)
}
