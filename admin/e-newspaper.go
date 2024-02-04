package admin

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
)

func addENewspaper(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 20) // Max memory 20MB
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	title_latin := r.FormValue("title_latin")
	title_cyrillic := r.FormValue("title_cyrillic")

	file_latin, file_latin_header, err := r.FormFile("file_latin")
	var fileLatinForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: file_latin error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		fileLatinForDB = nil
	} else {
		// Check size limits
		if file_latin_header.Size > int64(6<<20) {
			response.Res(w, "error", http.StatusBadRequest, "file_latin exceeds 6MB limit")
			return
		}
		fileLatinForDB, _ = io.ReadAll(file_latin)
		file_latin.Close()
	}

	file_cyrillic, file_cyrillic_header, err := r.FormFile("file_cyrillic")
	var fileCyrillicForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: file_cyrillic error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		fileCyrillicForDB = nil
	} else {
		if file_cyrillic_header.Size > int64(6<<20) {
			response.Res(w, "error", http.StatusBadRequest, "file_cyrillic exceeds 6MB limit")
			return
		}
		fileCyrillicForDB, _ = io.ReadAll(file_cyrillic)
		file_cyrillic.Close()
	}

	cover_image, cover_image_header, err := r.FormFile("cover_image")
	var coverImageForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: cover_image error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		coverImageForDB = nil
	} else {
		if cover_image_header.Size > int64(3<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Cover image exceeds 3MB limit")
			return
		}
		coverImageForDB, _ = io.ReadAll(cover_image)
		cover_image.Close()
	}

	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO e_newspapers (title_latin, title_cyrillic, file_latin, file_cyrillic, cover_image) VALUES ($1, $2, $3, $4, $5)`,
		title_latin, title_cyrillic, fileLatinForDB, fileCyrillicForDB, coverImageForDB)
	if err != nil {
		log.Printf("%v: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "e-newspaper has been added successfully.")
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

func eNewspaperIsArchived(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT archived FROM e_newspapers WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute the SQL statement
	var archived bool
	err = stmt.QueryRow(id).Scan(&archived)
	if err != nil {
		return nil, err
	}

	return &archived, nil
}

func editENewspaper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := eNewspaperExists(id)
	if err != nil {
		log.Printf("%v: edit e-newspaper eNewspaperExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: edit e-newspaper eNewspaperExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit non existent e-newspaper")
		return
	}

	archived, err := eNewspaperIsArchived(id)
	if err != nil {
		log.Printf("%v: edit e-newspaper eNewspaperIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: edit e-newspaper eNewspaperIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit archived e-newspaper")
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(15 << 20)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error while connecting to db: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	title_latin := r.FormValue("title_latin")
	if title_latin != "" {
		sqlStatement := `
			UPDATE e_newspapers
			SET title_latin = $1, edited_at = NOW() 
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, title_latin, id)
		if err != nil {
			log.Printf("%v: writing title_latin into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	title_cyrillic := r.FormValue("title_cyrillic")
	if title_cyrillic != "" {
		sqlStatement := `
			UPDATE e_newspapers
			SET title_cyrillic = $1, edited_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, title_cyrillic, id)
		if err != nil {
			log.Printf("%v: writing title_cyrillic into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	file_latin, file_latin_header, err := r.FormFile("file_latin")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: file_latin error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
	} else {
		// Check size limits
		if file_latin_header.Size > int64(6<<20) {
			response.Res(w, "error", http.StatusBadRequest, "file_latin exceeds 6MB limit")
			return
		}
		fileLatinForDB, _ := io.ReadAll(file_latin)
		file_latin.Close()
		sqlStatement := `
			UPDATE e_newspapers
			SET file_latin = $1, edited_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, fileLatinForDB, id)
		if err != nil {
			log.Printf("%v: writing file_latin into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	file_cyrillic, file_cyrillic_header, err := r.FormFile("file_cyrillic")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: file_cyrillic error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
	} else {
		if file_cyrillic_header.Size > int64(6<<20) {
			response.Res(w, "error", http.StatusBadRequest, "file_cyrillic exceeds 6MB limit")
			return
		}
		fileCyrillicForDB, _ := io.ReadAll(file_cyrillic)
		file_cyrillic.Close()
		sqlStatement := `
			UPDATE e_newspapers
			SET file_cyrillic = $1, edited_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, fileCyrillicForDB, id)
		if err != nil {
			log.Printf("%v: writing file_cyrillic into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	cover_image, cover_image_header, err := r.FormFile("cover_image")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: cover_image error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
	} else {
		if cover_image_header.Size > int64(3<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Cover image exceeds 3MB limit")
			return
		}
		coverImageForDB, _ := io.ReadAll(cover_image)
		cover_image.Close()
		sqlStatement := `
			UPDATE e_newspapers
			SET cover_image = $1, edited_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, coverImageForDB, id)
		if err != nil {
			log.Printf("%v: writing cover_image into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusOK, "E-newspaper edited")
}

func deleteENewspaper(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	exists, err := eNewspaperExists(id)
	if err != nil {
		log.Printf("%v: delete e-newspaper eNewspaperExists(id) error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: delete e-newspaper eNewspaperExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete non existent e-newspaper")
		return
	}

	archived, err := eNewspaperIsArchived(id)
	if err != nil {
		log.Printf("%v: delete e-newspaper eNewspaperIsArchived(id) error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: delete e-newspaper eNewspaperIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete archived e-newspaper")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: db error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM e_newspapers WHERE id=$1")
	if err != nil {
		log.Printf("%v: db error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Printf("%v: db statement execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "deleted")
}

func archiveENewspaper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := eNewspaperExists(id)
	if err != nil {
		log.Printf("%v: archive e-newspaper eNewspaperExists(id) error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: archive e-newspaper eNewspaperExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive non existent e-newspaper")
		return
	}

	archived, err := eNewspaperIsArchived(id)
	if err != nil {
		log.Printf("%v: archive e-newspaper eNewspaperIsArchived(id) error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: archive e-newspaper eNewspaperIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive already archived e-newspaper")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: db error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE e_newspapers SET archived = true WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: db error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "archived")
}

func unArchiveENewspaper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := eNewspaperExists(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: unarchive e-newspaper eNewspaperExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot unarchive non existent e-newspaper")
		return
	}

	archived, err := eNewspaperIsArchived(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*archived {
		log.Printf("%v: unarchive e-newspaper eNewspaperIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot unarchive not archived e-newspaper")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: db error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE e_newspapers SET archived = false WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: db error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "unarchive done")
}

type ENewspaperCount struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

func getENewspaperCountAll(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM e_newspapers").Scan(&count)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, ENewspaperCount{Period: "all", Count: count})
}

func getENewspaperCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	period := vars["period"]

	if period != "week" && period != "month" && period != "year" {
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
	query := fmt.Sprintf("SELECT COUNT(*) FROM e_newspapers WHERE created_at > current_date - interval '1 %s'", period)
	err = database.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, ENewspaperCount{Period: period, Count: count})
}

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

	rows, err := database.Query("SELECT * FROM e_newspapers ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var eNewspapers []ENewspaper
	for rows.Next() {
		var eNewspaper ENewspaper
		err := rows.Scan(&eNewspaper.ID, &eNewspaper.TitleLatin, &eNewspaper.TitleCyrillic, &eNewspaper.FileLatin, &eNewspaper.FileCyrillic, &eNewspaper.CoverImage, &eNewspaper.CreatedAt, &eNewspaper.UpdatedAt, &eNewspaper.Archived)
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
	FileLatin     []byte `json:"file_latin"`
	FileCyrillic  []byte `json:"file_cyrillic"`
	CoverImage    []byte `json:"cover_image"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Archived      bool   `json:"archived"`
}
