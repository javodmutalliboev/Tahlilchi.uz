package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type Category struct {
	TitleLatin          string
	DescriptionLatin    string
	TitleCyrillic       string
	DescriptionCyrillic string
}

type CategoryForGet struct {
	ID                  int                 `json:"id"`
	TitleLatin          string              `json:"title_latin"`
	DescriptionLatin    string              `json:"description_latin"`
	TitleCyrillic       string              `json:"title_cyrillic"`
	DescriptionCyrillic string              `json:"description_cyrillic"`
	Subcategories       []SubcategoryForGet `json:"subcategories"`
}

func getCategoryList(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT * FROM news_category ORDER BY id")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var categories []CategoryForGet
	for rows.Next() {
		var c CategoryForGet
		err := rows.Scan(&c.ID, &c.TitleLatin, &c.DescriptionLatin, &c.TitleCyrillic, &c.DescriptionCyrillic)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		subRows, err := database.Query("SELECT * FROM news_subcategory WHERE category_id = $1", c.ID)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		defer subRows.Close()

		for subRows.Next() {
			var s SubcategoryForGet
			err := subRows.Scan(&s.ID, &s.CategoryID, &s.TitleLatin, &s.DescriptionLatin, &s.TitleCyrillic, &s.DescriptionCyrillic)
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
			c.Subcategories = append(c.Subcategories, s)
		}

		categories = append(categories, c)
	}

	response.Res(w, "success", http.StatusOK, categories)
}

func getSubCategoryList(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT * FROM news_subcategory ORDER BY id")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var categories []SubcategoryForGet
	for rows.Next() {
		var c SubcategoryForGet
		err := rows.Scan(&c.ID, &c.CategoryID, &c.TitleLatin, &c.DescriptionLatin, &c.TitleCyrillic, &c.DescriptionCyrillic)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		categories = append(categories, c)
	}

	response.Res(w, "success", http.StatusOK, categories)
}

// updateSubCategory is a route handler function to update subcategory by category id and subcategory id
func updateSubCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send the error response
		response.Res(w, "error", http.StatusBadRequest, "invalid id")
		return
	}

	id, err := strconv.Atoi(vars["sub_id"])
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send the error response
		response.Res(w, "error", http.StatusBadRequest, "invalid sub_id")
		return
	}

	// parse the request body
	var subcategory Subcategory
	err = json.NewDecoder(r.Body).Decode(&subcategory)
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send the error response
		response.Res(w, "error", http.StatusBadRequest, "invalid request body")
		return
	}

	// open a connection to the database
	database, err := db.DB()
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send the error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	// defer the close of the database connection
	defer database.Close()

	// title_latin
	if subcategory.TitleLatin != "" {
		_, err = database.Exec("UPDATE news_subcategory SET title_latin = $1 WHERE id = $2 AND category_id = $3", subcategory.TitleLatin, id, category_id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send the error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// description_latin
	if subcategory.DescriptionLatin != "" {
		_, err = database.Exec("UPDATE news_subcategory SET description_latin = $1 WHERE id = $2 AND category_id = $3", subcategory.DescriptionLatin, id, category_id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send the error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// title_cyrillic
	if subcategory.TitleCyrillic != "" {
		_, err = database.Exec("UPDATE news_subcategory SET title_cyrillic = $1 WHERE id = $2 AND category_id = $3", subcategory.TitleCyrillic, id, category_id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send the error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// description_cyrillic
	if subcategory.DescriptionCyrillic != "" {
		_, err = database.Exec("UPDATE news_subcategory SET description_cyrillic = $1 WHERE id = $2 AND category_id = $3", subcategory.DescriptionCyrillic, id, category_id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send the error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// send the success response
	response.Res(w, "success", http.StatusOK, "subcategory updated")
}

// deleteSubCategory is a route handler function to delete subcategory by category id and subcategory id
func deleteSubCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	category_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send the error response
		response.Res(w, "error", http.StatusBadRequest, "invalid id")
		return
	}

	id, err := strconv.Atoi(vars["sub_id"])
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send the error response
		response.Res(w, "error", http.StatusBadRequest, "invalid sub_id")
		return
	}

	// open a connection to the database
	database, err := db.DB()
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send the error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	// defer the close of the database connection
	defer database.Close()

	// delete the subcategory
	_, err = database.Exec("DELETE FROM news_subcategory WHERE id = $1 AND category_id = $2", id, category_id)
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send the error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// send the success response
	response.Res(w, "success", http.StatusOK, "subcategory deleted")
}

// getSubCategoryListByCategory returns subcategories by category id
func getSubCategoryListByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("SELECT * FROM news_subcategory WHERE category_id = $1 ORDER BY id", categoryID)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var categories []SubcategoryForGet
	for rows.Next() {
		var c SubcategoryForGet
		err := rows.Scan(&c.ID, &c.CategoryID, &c.TitleLatin, &c.DescriptionLatin, &c.TitleCyrillic, &c.DescriptionCyrillic)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		categories = append(categories, c)
	}

	response.Res(w, "success", http.StatusOK, categories)
}

type SubcategoryForGet struct {
	ID                  int    `json:"id"`
	CategoryID          int    `json:"category_id"`
	TitleLatin          string `json:"title_latin"`
	DescriptionLatin    string `json:"description_latin"`
	TitleCyrillic       string `json:"title_cyrillic"`
	DescriptionCyrillic string `json:"description_cyrillic"`
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	var c Category
	c.TitleLatin = r.FormValue("title_latin")
	c.DescriptionLatin = r.FormValue("description_latin")
	c.TitleCyrillic = r.FormValue("title_cyrillic")
	c.DescriptionCyrillic = r.FormValue("description_cyrillic")

	if c.TitleLatin == "" {
		response.Res(w, "error", http.StatusBadRequest, "Title latin is required")
		return
	}

	if c.TitleCyrillic == "" {
		response.Res(w, "error", http.StatusBadRequest, "Title cyrillic is required")
		return
	}

	_, err = db.Exec("INSERT INTO news_category(title_latin, description_latin, title_cyrillic, description_cyrillic) VALUES($1, $2, $3, $4)", c.TitleLatin, c.DescriptionLatin, c.TitleCyrillic, c.DescriptionCyrillic)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "Category Added")
}

type Subcategory struct {
	ID                  int    `json:"id"`
	CategoryID          int    `json:"category_id"`
	TitleLatin          string `json:"title_latin"`
	DescriptionLatin    string `json:"description_latin"`
	TitleCyrillic       string `json:"title_cyrillic"`
	DescriptionCyrillic string `json:"description_cyrillic"`
}

func addSubcategory(w http.ResponseWriter, r *http.Request) {
	var s Subcategory
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// Check if category_id and titles are provided
	if s.CategoryID == 0 || s.TitleLatin == "" || s.TitleCyrillic == "" {
		response.Res(w, "error", http.StatusBadRequest, "category_id and titles are required")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO news_subcategory (category_id, title_latin, description_latin, title_cyrillic, description_cyrillic) VALUES ($1, $2, $3, $4, $5)", s.CategoryID, s.TitleLatin, s.DescriptionLatin, s.TitleCyrillic, s.DescriptionCyrillic)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "subcategory added")
}

type Region struct {
	ID                  int    `json:"id"`
	NameLatin           string `json:"name_latin"`
	DescriptionLatin    string `json:"description_latin,omitempty"`
	NameCyrillic        string `json:"name_cyrillic"`
	DescriptionCyrillic string `json:"description_cyrillic"`
}

func addRegions(w http.ResponseWriter, r *http.Request) {
	var regions []Region

	err := json.NewDecoder(r.Body).Decode(&regions)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	for _, region := range regions {
		if region.NameLatin == "" || region.NameCyrillic == "" {
			response.Res(w, "error", http.StatusBadRequest, "Name fields are required")
			return
		}

		_, err = db.Exec("INSERT INTO news_regions (name_latin, description_latin, name_cyrillic, description_cyrillic) VALUES ($1, $2, $3, $4)", region.NameLatin, region.DescriptionLatin, region.NameCyrillic, region.DescriptionCyrillic)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusCreated, "Regions added successfully")
}

var re = regexp.MustCompile(`^(true|false)$`)

func addNewsPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1 << 30) // Max memory 1GB
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	title_latin := r.FormValue("title_latin")
	description_latin := r.FormValue("description_latin")
	title_cyrillic := r.FormValue("title_cyrillic")
	description_cyrillic := r.FormValue("description_cyrillic")

	// Check if required fields are not empty
	if title_latin == "" || description_latin == "" || title_cyrillic == "" || description_cyrillic == "" {
		response.Res(w, "error", http.StatusBadRequest, "Required fields are missing")
		return
	}

	photo, _, err := r.FormFile("photo")
	var photoForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("photo: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		photoForDB = nil
	} else {
		photoForDB, _ = io.ReadAll(photo)
		photo.Close()
		// check file type
		filetype := http.DetectContentType(photoForDB)
		if !strings.HasPrefix(filetype, "image/") {
			log.Printf("%v: photo is not an image file", r.URL)
			response.Res(w, "error", http.StatusBadRequest, "photo is not an image file")
			return
		}
	}

	video := r.FormValue("video")
	/*
		var videoForDB []byte
		if err != nil && err != http.ErrMissingFile {
			log.Printf("video: %v", err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
			} else if err == http.ErrMissingFile {
				videoForDB = nil
				} else {
					if video_header.Size > int64(6<<20) {
						response.Res(w, "error", http.StatusBadRequest, "Video exceeds 6MB limit")
						return
					}
					videoForDB, _ = io.ReadAll(video)
					video.Close()
				}
	*/

	audio, _, err := r.FormFile("audio")
	var audioForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("audio: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		audioForDB = nil
	} else {
		audioForDB, _ = io.ReadAll(audio)
		audio.Close()
		// Check file type
		filetype := http.DetectContentType(audioForDB)
		if !strings.HasPrefix(filetype, "audio/") {
			log.Printf("%v: audio is not an audio file", r.URL)
			response.Res(w, "error", http.StatusBadRequest, "audio is not an audio file")
			return
		}
	}

	cover_image, _, err := r.FormFile("cover_image")
	var coverImageForDB []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("cover_image: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
		coverImageForDB = nil
	} else {
		coverImageForDB, _ = io.ReadAll(cover_image)
		cover_image.Close()

		// check file type
		filetype := http.DetectContentType(coverImageForDB)
		if !strings.HasPrefix(filetype, "image/") {
			log.Printf("%v: cover_image is not an image file", r.URL)
			response.Res(w, "error", http.StatusBadRequest, "cover_image is not an image file")
			return
		}
	}

	// Get tags if they exist
	tags, ok := r.Form["tags[]"]
	if !ok {
		// If tags don't exist, use an empty array
		tags = []string{}
	}

	// Convert tags to PostgreSQL array format
	tagsString := "{" + strings.Join(tags, ",") + "}"

	var categoryInt sql.NullInt64
	if category := r.FormValue("category"); category != "" {
		categoryInt.Int64, err = strconv.ParseInt(category, 10, 64)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		categoryInt.Valid = true
	}

	var subcategoryInt sql.NullInt64
	if subcategory := r.FormValue("subcategory"); subcategory != "" {
		subcategoryInt.Int64, err = strconv.ParseInt(subcategory, 10, 64)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		subcategoryInt.Valid = true
	}

	var regionInt sql.NullInt64
	if region := r.FormValue("region"); region != "" {
		regionInt.Int64, err = strconv.ParseInt(region, 10, 64)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		regionInt.Valid = true
	}

	var topBool sql.NullBool
	if top := r.FormValue("top"); top != "" {
		if !re.MatchString(top) {
			response.Res(w, "error", http.StatusBadRequest, "invalid top value")
			return
		}

		topBool.Bool, err = strconv.ParseBool(top)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		topBool.Valid = true
	}

	var latestBool sql.NullBool
	if latest := r.FormValue("latest"); latest != "" {
		if !re.MatchString(latest) {
			response.Res(w, "error", http.StatusBadRequest, "invalid latest value")
			return
		}

		latestBool.Bool, err = strconv.ParseBool(latest)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		latestBool.Valid = true
	}

	var relatedInt sql.NullInt64
	if related := r.FormValue("related"); related != "" {
		relatedInt.Int64, err = strconv.ParseInt(related, 10, 64)
		if err != nil {
			log.Println(err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		relatedInt.Valid = true
	}

	db, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO news_posts (title_latin, description_latin, title_cyrillic, description_cyrillic, photo, video, audio, cover_image, tags, category, subcategory, region, top, latest, related) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14 , $15)`,
		title_latin, description_latin, title_cyrillic, description_cyrillic, photoForDB, video, audioForDB, coverImageForDB, tagsString, categoryInt, subcategoryInt, regionInt, topBool, latestBool, relatedInt)
	if err != nil {
		log.Println(err, categoryInt)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "New post has been created successfully.")
}

func exists(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM news_posts WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

func isArchived(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT archived FROM news_posts WHERE id = $1")
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

func editNewsPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := exists(id)
	if err != nil {
		log.Printf("%v: edit news post exists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: edit news post exists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit non existent news post")
		return
	}

	archived, err := isArchived(id)
	if err != nil {
		log.Printf("%v: edit news post isArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: edit news post isArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit archived news post")
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(15 << 20)
	if err != nil {
		log.Printf("%v: edit news post: %v", r.URL, err)
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
			UPDATE news_posts
			SET title_latin = $1, updated_at = NOW() 
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, title_latin, id)
		if err != nil {
			log.Printf("%v: writing title_latin into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	description_latin := r.FormValue("description_latin")
	if description_latin != "" {
		sqlStatement := `
			UPDATE news_posts
			SET description_latin = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, description_latin, id)
		if err != nil {
			log.Printf("%v: writing description_latin into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	title_cyrillic := r.FormValue("title_cyrillic")
	if title_cyrillic != "" {
		sqlStatement := `
			UPDATE news_posts
			SET title_cyrillic = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, title_cyrillic, id)
		if err != nil {
			log.Printf("%v: writing title_cyrillic into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	description_cyrillic := r.FormValue("description_cyrillic")
	if description_cyrillic != "" {
		sqlStatement := `
			UPDATE news_posts
			SET description_cyrillic = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, description_cyrillic, id)
		if err != nil {
			log.Printf("%v: writing description_cyrillic into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	photo, _, err := r.FormFile("photo")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("photo: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
	} else {
		photoForDB, _ := io.ReadAll(photo)
		photo.Close()

		// check file type
		filetype := http.DetectContentType(photoForDB)
		if !strings.HasPrefix(filetype, "image/") {
			log.Printf("%v: photo is not an image file", r.URL)
			response.Res(w, "error", http.StatusBadRequest, "photo is not an image file")
			return
		}

		sqlStatement := `
			UPDATE news_posts
			SET photo = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, photoForDB, id)
		if err != nil {
			log.Printf("%v: writing photo into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	video := r.FormValue("video")
	if video != "" {
		sqlStatement := `
			UPDATE news_posts
			SET video = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, video, id)
		if err != nil {
			log.Printf("%v: writing video into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}
	/*
		if err != nil && err != http.ErrMissingFile {
			log.Printf("video: %v", err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		} else if err == http.ErrMissingFile {
		} else {
			if video_header.Size > int64(6<<20) {
				response.Res(w, "error", http.StatusBadRequest, "Video exceeds 6MB limit")
				return
			}
			videoForDB, _ := io.ReadAll(video)
			video.Close()
			sqlStatement := `
				UPDATE news_posts
				SET video = $1, updated_at = NOW()
				WHERE id = $2;
			`
			_, err = db.Exec(sqlStatement, videoForDB, id)
			if err != nil {
				log.Printf("%v: writing video into db: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
		}
	*/

	audio, _, err := r.FormFile("audio")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("audio: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
	} else {
		audioForDB, _ := io.ReadAll(audio)
		audio.Close()

		// Check file type
		filetype := http.DetectContentType(audioForDB)
		if !strings.HasPrefix(filetype, "audio/") {
			log.Printf("%v: audio is not an audio file", r.URL)
			response.Res(w, "error", http.StatusBadRequest, "audio is not an audio file")
			return
		}

		sqlStatement := `
			UPDATE news_posts
			SET audio = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, audioForDB, id)
		if err != nil {
			log.Printf("%v: writing audio into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	cover_image, _, err := r.FormFile("cover_image")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("cover_image: %v", err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	} else if err == http.ErrMissingFile {
	} else {
		coverImageForDB, _ := io.ReadAll(cover_image)
		cover_image.Close()

		// check file type
		filetype := http.DetectContentType(coverImageForDB)
		if !strings.HasPrefix(filetype, "image/") {
			log.Printf("%v: cover_image is not an image file", r.URL)
			response.Res(w, "error", http.StatusBadRequest, "cover_image is not an image file")
			return
		}

		sqlStatement := `
			UPDATE news_posts
			SET cover_image = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, coverImageForDB, id)
		if err != nil {
			log.Printf("%v: writing cover_image into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if tags, ok := r.Form["tags[]"]; ok {
		tagsString := "{" + strings.Join(tags, ",") + "}"
		sqlStatement := `
			UPDATE news_posts
			SET tags = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, tagsString, id)
		if err != nil {
			log.Printf("%v: writing tags into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if category := r.FormValue("category"); category != "" {
		categoryInt, err := strconv.ParseInt(category, 10, 64)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		sqlStatement := `
			UPDATE news_posts
			SET category = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, categoryInt, id)
		if err != nil {
			log.Printf("%v: writing category into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if subcategory := r.FormValue("subcategory"); subcategory != "" {
		subcategoryInt, err := strconv.ParseInt(subcategory, 10, 64)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		sqlStatement := `
			UPDATE news_posts
			SET subcategory = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, subcategoryInt, id)
		if err != nil {
			log.Printf("%v: writing subcategory into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if region := r.FormValue("region"); region != "" {
		regionInt, err := strconv.ParseInt(region, 10, 64)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		sqlStatement := `
			UPDATE news_posts
			SET region = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, regionInt, id)
		if err != nil {
			log.Printf("%v: writing region into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if top := r.FormValue("top"); top != "" {
		if !re.MatchString(top) {
			response.Res(w, "error", http.StatusBadRequest, "invalid top value")
			return
		}

		topBool, err := strconv.ParseBool(top)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		sqlStatement := `
			UPDATE news_posts
			SET top = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, topBool, id)
		if err != nil {
			log.Printf("%v: writing top into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if latest := r.FormValue("latest"); latest != "" {
		if !re.MatchString(latest) {
			response.Res(w, "error", http.StatusBadRequest, "invalid latest value")
			return
		}

		latestBool, err := strconv.ParseBool(latest)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		sqlStatement := `
			UPDATE news_posts
			SET latest = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, latestBool, id)
		if err != nil {
			log.Printf("%v: writing latest into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if related := r.FormValue("related"); related != "" {
		relatedInt, err := strconv.ParseInt(related, 10, 64)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		sqlStatement := `
			UPDATE news_posts
			SET related = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, relatedInt, id)
		if err != nil {
			log.Printf("%v: writing related into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusOK, "News post edited")
}

func deleteNewsPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	exists, err := exists(id)
	if err != nil {
		log.Printf("%v: delete news post exists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: delete news post exists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete non existent news post")
		return
	}

	archived, err := isArchived(id)
	if err != nil {
		log.Printf("%v: delete news post isArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: delete news post isArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete archived news post")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM news_posts WHERE id=$1")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	_, err = stmt.Exec(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "deleted")
}

func archiveNewsPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := exists(id)
	if err != nil {
		log.Printf("%v: archive news post exists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: archive news post exists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive non existent news post")
		return
	}

	archived, err := isArchived(id)
	if err != nil {
		log.Printf("%v: archive news post isArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: archive news post isArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive already archived news post")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE news_posts SET archived = true WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "archived")
}

func unArchiveNewsPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := exists(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: unarchive news post exists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot unarchive non existent news post")
		return
	}

	archived, err := isArchived(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*archived {
		log.Printf("%v: unarchive news post isArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot unarchive not archived news post")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE news_posts SET archived = false WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "unarchive done")
}

type NewsPostCount struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

// get count of all news posts
func getNewsPostCountAll(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM news_posts").Scan(&count)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, NewsPostCount{Period: "all", Count: count})
}

func getNewsPostCount(w http.ResponseWriter, r *http.Request) {
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
	query := fmt.Sprintf("SELECT COUNT(*) FROM news_posts WHERE created_at > current_date - interval '1 %s'", period)
	err = database.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, NewsPostCount{Period: period, Count: count})
}

func getNewsPosts(w http.ResponseWriter, r *http.Request) {
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

	// Query the database
	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, video, tags, archived, created_at, updated_at, category, subcategory, region, top, latest, related, completed FROM news_posts ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	// Write the posts to the response
	var posts []NewsPost
	for rows.Next() {
		var p NewsPost
		if err := rows.Scan(&p.ID, &p.TitleLatin, &p.DescriptionLatin, &p.TitleCyrillic, &p.DescriptionCyrillic, &p.Video, pq.Array(&p.Tags), &p.Archived, &p.CreatedAt, &p.UpdatedAt, &p.Category, &p.Subcategory, &p.Region, &p.Top, &p.Latest, &p.Related, &p.Completed); err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		posts = append(posts, p)
	}

	hasPreviousPage := start > 0
	// get total count of news posts
	var total int
	err = database.QueryRow("SELECT COUNT(*) FROM news_posts").Scan(&total)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	hasNextPage := total > (start + limit)

	data := ResponseNewsPostsData{
		NewsPosts:    posts,
		NextPage:     hasNextPage,
		PreviousPage: hasPreviousPage,
	}

	response.Res(w, "success", http.StatusOK, data)
}

type NewsPost struct {
	ID                  int       `json:"id"`
	TitleLatin          string    `json:"title_latin"`
	DescriptionLatin    string    `json:"description_latin"`
	TitleCyrillic       string    `json:"title_cyrillic"`
	DescriptionCyrillic string    `json:"description_cyrillic"`
	Video               string    `json:"video"`
	Tags                []string  `json:"tags"`
	Archived            bool      `json:"archived"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Category            *int      `json:"category"`
	Subcategory         *int      `json:"subcategory"`
	Region              *int      `json:"region"`
	Top                 *bool     `json:"top"`
	Latest              *bool     `json:"latest"`
	Related             *int      `json:"related"`
	Completed           bool      `json:"completed"`
}

type ResponseNewsPostsData struct {
	NewsPosts    []NewsPost `json:"news_posts"`
	NextPage     bool       `json:"next_page"`
	PreviousPage bool       `json:"previous_page"`
}

// newsPostCompleted is a handler to update the completed field of a news post
func newsPostCompleted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := exists(id)
	if err != nil {
		log.Printf("%v: news post completed exists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: news post completed exists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot update completed field of non existent news post")
		return
	}

	archived, err := isArchived(id)
	if err != nil {
		log.Printf("%v: news post completed isArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: news post completed isArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot update completed field of archived news post")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE news_posts SET completed = NOT completed WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "completed field updated")
}

func getNewsPostPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var photo []byte
	err = database.QueryRow("SELECT photo FROM news_posts WHERE id = $1", id).Scan(&photo)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(photo)
}

func getNewsPostAudio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var audio []byte
	err = database.QueryRow("SELECT audio FROM news_posts WHERE id = $1", id).Scan(&audio)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Write(audio)
}

func getNewsPostCoverImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var CoverImage []byte
	err = database.QueryRow("SELECT cover_image FROM news_posts WHERE id = $1", id).Scan(&CoverImage)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(CoverImage)
}

// getRegions is a route handler function to get all news regions
func getRegions(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// create a slice to hold the regions
	var regions []Region

	// query the database
	rows, err := database.Query("SELECT * FROM news_regions order by id")
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	// iterate through the rows
	for rows.Next() {
		// create a region variable
		var region Region
		// scan the row into the region variable
		if err := rows.Scan(&region.ID, &region.NameLatin, &region.DescriptionLatin, &region.NameCyrillic, &region.DescriptionCyrillic); err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send an error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		// append the region to the regions slice
		regions = append(regions, region)
	}

	// send a success response
	response.Res(w, "success", http.StatusOK, regions)
}

// updateRegion is a route handler function to update a news region
func updateRegion(w http.ResponseWriter, r *http.Request) {
	// get the region id from the url parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// parse the request body
	err := r.ParseForm()
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusBadRequest, "invalid request body")
		return
	}

	// open a connection to the database
	database, err := db.DB()
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// check news region existence
	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM news_regions WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !exists {
		response.Res(w, "error", http.StatusBadRequest, "news region does not exist")
		return
	}

	// name_latin
	nameLatin := r.FormValue("name_latin")
	if nameLatin != "" {
		_, err = database.Exec("UPDATE news_regions SET name_latin = $1 WHERE id = $2", nameLatin, id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send an error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// description_latin
	descriptionLatin := r.FormValue("description_latin")
	if descriptionLatin != "" {
		_, err = database.Exec("UPDATE news_regions SET description_latin = $1 WHERE id = $2", descriptionLatin, id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send an error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// name_cyrillic
	nameCyrillic := r.FormValue("name_cyrillic")
	if nameCyrillic != "" {
		_, err = database.Exec("UPDATE news_regions SET name_cyrillic = $1 WHERE id = $2", nameCyrillic, id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send an error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// description_cyrillic
	descriptionCyrillic := r.FormValue("description_cyrillic")
	if descriptionCyrillic != "" {
		_, err = database.Exec("UPDATE news_regions SET description_cyrillic = $1 WHERE id = $2", descriptionCyrillic, id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send an error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// send a success response
	response.Res(w, "success", http.StatusOK, "news region updated")
}

// deleteRegion is a route handler function to delete a news region
func deleteRegion(w http.ResponseWriter, r *http.Request) {
	// get the region id from the url parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// open a connection to the database
	database, err := db.DB()
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// check news region existence
	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM news_regions WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !exists {
		response.Res(w, "error", http.StatusBadRequest, "news region does not exist")
		return
	}

	// delete the news region
	_, err = database.Exec("DELETE FROM news_regions WHERE id = $1", id)
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// check err is constraint violation
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			response.Res(w, "error", http.StatusBadRequest, "news region is in use")
			return
		}
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// send a success response
	response.Res(w, "success", http.StatusOK, "news region deleted")
}

// updateCategory is a route handler function to update a news category
func updateCategory(w http.ResponseWriter, r *http.Request) {
	// get the category id from the url parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// parse the request body
	err := r.ParseForm()
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusBadRequest, "invalid request body")
		return
	}

	// open a connection to the database
	database, err := db.DB()
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// check news category existence
	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM news_category WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !exists {
		response.Res(w, "error", http.StatusBadRequest, "news category does not exist")
		return
	}

	// title_latin
	titleLatin := r.FormValue("title_latin")
	if titleLatin != "" {
		_, err = database.Exec("UPDATE news_category SET title_latin = $1 WHERE id = $2", titleLatin, id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send an error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// description_latin
	descriptionLatin := r.FormValue("description_latin")
	if descriptionLatin != "" {
		_, err = database.Exec("UPDATE news_category SET description_latin = $1 WHERE id = $2", descriptionLatin, id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send an error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// title_cyrillic
	titleCyrillic := r.FormValue("title_cyrillic")
	if titleCyrillic != "" {
		_, err = database.Exec("UPDATE news_category SET title_cyrillic = $1 WHERE id = $2", titleCyrillic, id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send an error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// description_cyrillic
	descriptionCyrillic := r.FormValue("description_cyrillic")
	if descriptionCyrillic != "" {
		_, err = database.Exec("UPDATE news_category SET description_cyrillic = $1 WHERE id = $2", descriptionCyrillic, id)
		if err != nil {
			// log the error
			toolkit.LogError(r, err)
			// send an error response
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// send a success response
	response.Res(w, "success", http.StatusOK, "news category updated")
}

// deleteCategory is a route handler function to delete a news category
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// get the category id from the url parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// open a connection to the database
	database, err := db.DB()
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// check news category existence
	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM news_category WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !exists {
		response.Res(w, "error", http.StatusBadRequest, "news category does not exist")
		return
	}

	// delete the news category
	_, err = database.Exec("DELETE FROM news_category WHERE id = $1", id)
	if err != nil {
		// log the error
		toolkit.LogError(r, err)
		// check err is constraint violation
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			response.Res(w, "error", http.StatusBadRequest, "news category is in use")
			return
		}
		// send an error response
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// send a success response
	response.Res(w, "success", http.StatusOK, "news category deleted")
}
