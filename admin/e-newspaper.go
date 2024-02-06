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

// e-newspaper category type
type ENewspaperCategory struct {
	ID            int    `json:"id"`
	TitleLatin    string `json:"title_latin"`
	TitleCyrillic string `json:"title_cyrillic"`
}

// addENewspaperCategory is a handler to add new e-newspaper category
func addENewspaperCategory(w http.ResponseWriter, r *http.Request) {
	var e ENewspaperCategory
	err := r.ParseForm()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	e.TitleLatin = r.FormValue("title_latin")
	e.TitleCyrillic = r.FormValue("title_cyrillic")

	if e.TitleLatin == "" || e.TitleCyrillic == "" {
		log.Printf("%v: title_latin: %v; title_cyrillic: %v", r.URL, e.TitleLatin, e.TitleCyrillic)
		response.Res(w, "error", http.StatusBadRequest, "Both title_latin and title_cyrillic are required")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: db connection error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	_, err = database.Exec("INSERT INTO e_newspaper_category (title_latin, title_cyrillic) VALUES ($1, $2)", e.TitleLatin, e.TitleCyrillic)
	if err != nil {
		log.Printf("%v: db execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "e-newspaper category added")
}

// getENewspaperCategoryList is a handler to get e-newspaper category list
func getENewspaperCategoryList(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: db connection error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT id, title_latin, title_cyrillic FROM e_newspaper_category")
	if err != nil {
		log.Printf("%v: db execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var eNewspaperCategories []ENewspaperCategory
	for rows.Next() {
		var e ENewspaperCategory
		err := rows.Scan(&e.ID, &e.TitleLatin, &e.TitleCyrillic)
		if err != nil {
			log.Printf("%v: db execution error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		eNewspaperCategories = append(eNewspaperCategories, e)
	}

	response.Res(w, "success", http.StatusOK, eNewspaperCategories)
}

// updateENewspaperCategory is a handler to update e-newspaper category
func updateENewspaperCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// check e-newspaper category exists
	var exists bool
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: db connection error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM e_newspaper_category WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		log.Printf("%v: db execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !exists {
		response.Res(w, "error", http.StatusNotFound, "e-newspaper category not found")
		return
	}

	var e ENewspaperCategory
	err = r.ParseForm()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	e.TitleLatin = r.FormValue("title_latin")
	e.TitleCyrillic = r.FormValue("title_cyrillic")

	if e.TitleLatin != "" {
		// update title_latin column in database
		_, err = database.Exec("UPDATE e_newspaper_category SET title_latin = $1 WHERE id = $2", e.TitleLatin, id)
		if err != nil {
			log.Printf("%v: db execution error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if e.TitleCyrillic != "" {
		// update title_cyrillic column in database
		_, err = database.Exec("UPDATE e_newspaper_category SET title_cyrillic = $1 WHERE id = $2", e.TitleCyrillic, id)
		if err != nil {
			log.Printf("%v: db execution error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusOK, "e-newspaper category updated")
}

// deleteENewspaperCategory is a handler to delete e-newspaper category
func deleteENewspaperCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: db connection error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// check e-newspaper category exists
	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM e_newspaper_category WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		log.Printf("%v: db execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !exists {
		response.Res(w, "error", http.StatusNotFound, "e-newspaper category not found")
		return
	}

	_, err = database.Exec("DELETE FROM e_newspaper_category WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: db execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "e-newspaper category deleted")
}

func addENewspaper(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(80 << 20) // Max memory 80MB
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
		if file_latin_header.Size > int64(30<<20) {
			response.Res(w, "error", http.StatusBadRequest, "file_latin exceeds 6MB limit")
			return
		}
		// check whether file is pdf
		if file_latin_header.Header.Get("Content-Type") != "application/pdf" {
			response.Res(w, "error", http.StatusBadRequest, "file_latin is not a pdf")
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
		if file_cyrillic_header.Size > int64(30<<20) {
			response.Res(w, "error", http.StatusBadRequest, "file_cyrillic exceeds 6MB limit")
			return
		}
		// check whether file is pdf
		if file_cyrillic_header.Header.Get("Content-Type") != "application/pdf" {
			response.Res(w, "error", http.StatusBadRequest, "file_cyrillic is not a pdf")
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
		if cover_image_header.Size > int64(15<<20) {
			response.Res(w, "error", http.StatusBadRequest, "Cover image exceeds 3MB limit")
			return
		}
		// check whether file is image
		if cover_image_header.Header.Get("Content-Type") != "image/jpeg" && cover_image_header.Header.Get("Content-Type") != "image/png" {
			response.Res(w, "error", http.StatusBadRequest, "cover_image is not an image")
			return
		}
		coverImageForDB, _ = io.ReadAll(cover_image)
		cover_image.Close()
	}

	// category
	category := r.FormValue("category")
	// convert category to int
	categoryInt, err := strconv.Atoi(category)
	if err != nil {
		log.Printf("%v: category conversion error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, "category conversion error")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO e_newspapers (title_latin, title_cyrillic, file_latin, file_cyrillic, cover_image, category) VALUES ($1, $2, $3, $4, $5, $6)`,
		title_latin, title_cyrillic, fileLatinForDB, fileCyrillicForDB, coverImageForDB, categoryInt)
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

	rows, err := database.Query("SELECT id, title_latin, title_cyrillic, created_at, updated_at, archived, completed FROM e_newspapers ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var eNewspapers []ENewspaper
	for rows.Next() {
		var eNewspaper ENewspaper
		err := rows.Scan(&eNewspaper.ID, &eNewspaper.TitleLatin, &eNewspaper.TitleCyrillic, &eNewspaper.CreatedAt, &eNewspaper.UpdatedAt, &eNewspaper.Archived, &eNewspaper.Completed)
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
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Archived      bool   `json:"archived"`
	Completed     bool   `json:"completed"`
}

// eNewspaperCompleted is a handler to make e-newspaper completed field true/false
func eNewspaperCompleted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := eNewspaperExists(id)
	if err != nil {
		log.Printf("%v: eNewspaperExists(id) error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: eNewspaperExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot update completed field of non existent e-newspaper")
		return
	}

	archived, err := eNewspaperIsArchived(id)
	if err != nil {
		log.Printf("%v: eNewspaperIsArchived(id) error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: eNewspaperIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot update completed field of archived e-newspaper")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: db error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE e_newspapers SET completed = NOT completed WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: db error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "completed field updated")
}

// getENewspaperFile is a handler to get e-newspaper pdf latin or cyrillic file by id
func getENewspaperFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	alphabet := vars["alphabet"]

	if alphabet != "latin" && alphabet != "cyrillic" {
		response.Res(w, "error", http.StatusBadRequest, "invalid alphabet value")
		return
	}

	exists, err := eNewspaperExists(id)
	if err != nil {
		log.Printf("%v: eNewspaperExists(id) error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: eNewspaperExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusNotFound, "e-newspaper not found")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: db error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	var file []byte
	if alphabet == "latin" {
		err = db.QueryRow("SELECT file_latin FROM e_newspapers WHERE id = $1", id).Scan(&file)
		if err != nil {
			log.Printf("%v: db error: %v", r.URL, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
	}

	if alphabet == "cyrillic" {
		err = db.QueryRow("SELECT file_cyrillic FROM e_newspapers WHERE id = $1", id).Scan(&file)
		if err != nil {
			log.Printf("%v: db error: %v", r.URL, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
	}

	w.Header().Set("Content-Disposition", "attachment; filename=e-newspaper.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Length", strconv.Itoa(len(file)))
	w.Write(file)
}
