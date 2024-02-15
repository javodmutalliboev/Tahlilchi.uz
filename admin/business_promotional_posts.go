package admin

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/model"
	"Tahlilchi.uz/response"
	"Tahlilchi.uz/toolkit"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func addBusinessPromotionalPost(w http.ResponseWriter, r *http.Request) {
	// parse multipart form
	// maxMemory is 100MB
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// declare a new businessPromotionalPost
	var businessPromotionalPost model.BusinessPromotionalPost

	// title_latin
	titleLatin := r.FormValue("title_latin")
	if titleLatin == "" {
		toolkit.LogError(r, fmt.Errorf("title_latin is empty"))
		response.Res(w, "error", http.StatusBadRequest, "title_latin is empty")
		return
	}
	businessPromotionalPost.TitleLatin = titleLatin

	// description_latin
	descriptionLatin := r.FormValue("description_latin")
	if descriptionLatin != "" {
		businessPromotionalPost.DescriptionLatin = descriptionLatin
	}

	// title_cyrillic
	titleCyrillic := r.FormValue("title_cyrillic")
	if titleCyrillic == "" {
		toolkit.LogError(r, fmt.Errorf("title_cyrillic is empty"))
		response.Res(w, "error", http.StatusBadRequest, "title_cyrillic is empty")
		return
	}
	businessPromotionalPost.TitleCyrillic = titleCyrillic

	// description_cyrillic
	descriptionCyrillic := r.FormValue("description_cyrillic")
	if descriptionCyrillic != "" {
		businessPromotionalPost.DescriptionCyrillic = descriptionCyrillic
	}

	// photos
	photoArray := r.MultipartForm.File["photo"]
	for _, fh := range photoArray {
		if fh.Size > 10<<20 {
			err := fmt.Errorf("photo %v size exceeds 10MB limit: %v", fh.Filename, fh.Size)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		// check if file is image using http.DetectContentType
		contentType := fh.Header.Get("Content-Type")
		if contentType[:5] != "image" {
			err := fmt.Errorf("photo %v is not an image: %v", fh.Filename, contentType)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		file, err := fh.Open()
		if err != nil {
			err := fmt.Errorf("error opening photo %v: %v", fh.Filename, err)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		defer file.Close()
		photo, err := io.ReadAll(file)
		if err != nil {
			err := fmt.Errorf("error reading photo %v: %v", fh.Filename, err)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		businessPromotionalPost.Photos = append(businessPromotionalPost.Photos, model.BusinessPromotionalPostPhoto{FileName: fh.Filename, File: photo})
	}

	// videoArray
	// each video is string
	videos := r.MultipartForm.Value["video"]
	businessPromotionalPost.Videos = videos

	// cover_image
	_, coverImageHeader, err := r.FormFile("cover_image")
	if err != nil && err != http.ErrMissingFile {
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	} else if err == http.ErrMissingFile {
	} else {
		if coverImageHeader.Size > 10<<20 {
			err := fmt.Errorf("cover_image %v size exceeds 10MB limit: %v", coverImageHeader.Filename, coverImageHeader.Size)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		contentType := coverImageHeader.Header.Get("Content-Type")
		if contentType[:5] != "image" {
			err := fmt.Errorf("cover_image %v is not an image: %v", coverImageHeader.Filename, contentType)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		file, err := coverImageHeader.Open()
		if err != nil {
			err := fmt.Errorf("error opening cover_image %v: %v", coverImageHeader.Filename, err)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		defer file.Close()
		coverImage, err := io.ReadAll(file)
		if err != nil {
			err := fmt.Errorf("error reading cover_image %v: %v", coverImageHeader.Filename, err)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		businessPromotionalPost.CoverImage = coverImage
	}

	// expiration
	expiration := r.FormValue("expiration")
	expirationValid := checkExpiration(expiration)
	if !expirationValid {
		err := fmt.Errorf("expiration value %v is invalid", expiration)
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}
	switch expiration {
	case "1 day":
		businessPromotionalPost.Expiration = time.Now().UTC().Add(24 * time.Hour)
	case "1 week":
		businessPromotionalPost.Expiration = time.Now().UTC().Add(7 * 24 * time.Hour)
	case "1 month":
		businessPromotionalPost.Expiration = time.Now().UTC().AddDate(0, 1, 0)
	}

	// partner
	partner := r.FormValue("partner")
	if partner == "" {
		err := fmt.Errorf("partner is empty")
		toolkit.LogError(r, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}
	businessPromotionalPost.Partner = partner

	// Open a connection to the database
	database, err := db.DB()
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("database connection error: %v", err))
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// Start a new transaction
	tx, err := database.Begin()
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("start a new transaction: %v", err))
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// Prepare the SQL statement: add into business_promotional_posts: title_latin, description_latin, title_cyrillic, description_cyrillic, videos, cover_image, expiration, partner returning id
	stmt, err := tx.Prepare("INSERT INTO business_promotional_posts (title_latin, description_latin, title_cyrillic, description_cyrillic, videos, cover_image, expiration, partner) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("prepare the SQL statement: %v", err))
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// Execute the SQL statement
	err = stmt.QueryRow(businessPromotionalPost.TitleLatin, businessPromotionalPost.DescriptionLatin, businessPromotionalPost.TitleCyrillic, businessPromotionalPost.DescriptionCyrillic, pq.Array(businessPromotionalPost.Videos), businessPromotionalPost.CoverImage, businessPromotionalPost.Expiration, businessPromotionalPost.Partner).Scan(&businessPromotionalPost.ID)
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("execute the SQL statement: %v", err))
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		toolkit.LogError(r, fmt.Errorf("commit the transaction: %v", err))
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// add business promotional post photos into bpp_photos
	for _, photo := range businessPromotionalPost.Photos {
		// Start a new transaction
		tx, err := database.Begin()
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("start a new transaction: %v", err))
			// delete business promotional post from business_promotional_posts
			_, err = database.Exec("DELETE FROM business_promotional_posts WHERE id = $1", businessPromotionalPost.ID)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("delete business promotional post %v from business_promotional_posts: %v", businessPromotionalPost.ID, err))
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		// Prepare the SQL statement: add into bpp_photos: bpp, file_name, file
		stmt, err := tx.Prepare("INSERT INTO bpp_photos (bpp, file_name, file) VALUES ($1, $2, $3)")
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("prepare the SQL statement: %v", err))
			// delete business promotional post from business_promotional_posts
			_, err = database.Exec("DELETE FROM business_promotional_posts WHERE id = $1", businessPromotionalPost.ID)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("delete business promotional post %v from business_promotional_posts: %v", businessPromotionalPost.ID, err))
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		// Execute the SQL statement
		_, err = stmt.Exec(businessPromotionalPost.ID, photo.FileName, photo.File)
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("execute the SQL statement: %v", err))
			// delete business promotional post from business_promotional_posts
			_, err = database.Exec("DELETE FROM business_promotional_posts WHERE id = $1", businessPromotionalPost.ID)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("delete business promotional post %v from business_promotional_posts: %v", businessPromotionalPost.ID, err))
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("commit the transaction: %v", err))
			// delete business promotional post from business_promotional_posts
			_, err = database.Exec("DELETE FROM business_promotional_posts WHERE id = $1", businessPromotionalPost.ID)
			if err != nil {
				toolkit.LogError(r, fmt.Errorf("delete business promotional post %v from business_promotional_posts: %v", businessPromotionalPost.ID, err))
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusCreated, "business promotional post added")
}

func checkExpiration(expiration string) bool {
	pattern := regexp.MustCompile(`^1 (day|week|month)$`)
	return pattern.MatchString(expiration)
}

func bpPostExists(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM business_promotional_posts WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

func bpPostIsArchived(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT archived FROM business_promotional_posts WHERE id = $1")
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

func editBusinessPromotionalPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: edit business promotional post bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: edit business promotional post bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: edit business promotional post bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: edit business promotional post bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit archived business promotional post")
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(100 << 20) // maxMemory is 100MB
	if err != nil {
		log.Printf("%v: edit business promotional post: %v", r.URL, err)
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

	title_latin := r.FormValue("title_latin")
	if title_latin != "" {
		sqlStatement := `
			UPDATE business_promotional_posts
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
			UPDATE business_promotional_posts
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
			UPDATE business_promotional_posts
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
			UPDATE business_promotional_posts
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

	// videos
	videos := r.Form["video"]
	if len(videos) > 0 {
		sqlStatement := `
			UPDATE business_promotional_posts
			SET videos = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, pq.Array(videos), id)
		if err != nil {
			log.Printf("%v: writing videos into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	_, coverImageHeader, err := r.FormFile("cover_image")
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: FormFile(\"cover_image\"): %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	} else if err == http.ErrMissingFile {
	} else {
		if coverImageHeader.Size > 10<<20 {
			err := fmt.Errorf("cover_image %v size exceeds 10MB limit: %v", coverImageHeader.Filename, coverImageHeader.Size)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		contentType := coverImageHeader.Header.Get("Content-Type")
		if contentType[:5] != "image" {
			err := fmt.Errorf("cover_image %v is not an image: %v", coverImageHeader.Filename, contentType)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		file, err := coverImageHeader.Open()
		if err != nil {
			err := fmt.Errorf("error opening cover_image %v: %v", coverImageHeader.Filename, err)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		defer file.Close()
		coverImage, err := io.ReadAll(file)
		if err != nil {
			err := fmt.Errorf("error reading cover_image %v: %v", coverImageHeader.Filename, err)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		sqlStatement := `
			UPDATE business_promotional_posts
			SET cover_image = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, coverImage, id)
		if err != nil {
			log.Printf("%v: writing cover_image into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	expiration := r.FormValue("expiration")
	if expiration != "" {
		expirationValid := checkExpiration(expiration)
		if !expirationValid {
			log.Printf("%v: FormValue(\"expiration\") valid: %v", r.URL, expirationValid)
			response.Res(w, "error", http.StatusBadRequest, "expiration value is invalid")
			return
		}
		var expirationForDB time.Time
		switch expiration {
		case "1 day":
			expirationForDB = time.Now().UTC().Add(24 * time.Hour)
		case "1 week":
			expirationForDB = time.Now().UTC().Add(7 * 24 * time.Hour)
		case "1 month":
			expirationForDB = time.Now().UTC().AddDate(0, 1, 0)
		}
		sqlStatement := `
			UPDATE business_promotional_posts
			SET expiration = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, expirationForDB, id)
		if err != nil {
			log.Printf("%v: writing expiration into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	partner := r.FormValue("partner")
	if partner != "" {
		sqlStatement := `
			UPDATE business_promotional_posts
			SET partner = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, partner, id)
		if err != nil {
			log.Printf("%v: writing partner into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusOK, "business promotional post updated")
}

func deleteBPPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: delete business promotional post bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: delete business promotional post bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: delete business promotional post bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: delete business promotional post bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete archived business promotional post")
		return
	}

	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// first delete photos from bpp_photos
	_, err = database.Exec("DELETE FROM bpp_photos WHERE bpp = $1", id)
	if err != nil {
		log.Printf("%v: delete business promotional post photos: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// Prepare the SQL statement
	stmt, err := database.Prepare("DELETE FROM business_promotional_posts WHERE id=$1")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Printf("%v: delete business promotional post: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "business promotional post deleted")
}

func archiveBPPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: archive business promotional post bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: archive business promotional post bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: archive business promotional post bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: archive business promotional post bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive already archived business promotional post")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE business_promotional_posts SET archived = true WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "archived")
}

func unArchiveBPPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: archive business promotional post bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot unarchive non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*archived {
		log.Printf("%v: unarchive business promotional post bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot unarchive not archived business promotional post")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE business_promotional_posts SET archived = false WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "unarchive done")
}

func CheckAndArchiveExpiredBPPosts() {
	db, err := db.DB()
	if err != nil {
		log.Printf("checkAndArchiveExpiredBPPosts(): error: %v", err)
		return
	}
	defer db.Close()

	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("checkAndArchiveExpiredBPPosts(): Start a new transaction: error: %v", err)
		return
	}

	// Prepare the SQL statement
	stmt, err := tx.Prepare("UPDATE business_promotional_posts SET archived = true WHERE expiration < $1")
	if err != nil {
		log.Printf("checkAndArchiveExpiredBPPosts(): Prepare the SQL statement: error: %v", err)
		return
	}

	// Execute the SQL statement
	_, err = stmt.Exec(time.Now().UTC())
	if err != nil {
		log.Printf("checkAndArchiveExpiredBPPosts(): Execute the SQL statement: error: %v", err)
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Printf("checkAndArchiveExpiredBPPosts(): Commit the transaction: error: %v", err)
		return
	}
}

type businessPromotionalPostCount struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

func getBusinessPromotionalPostCount(w http.ResponseWriter, r *http.Request) {
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
	query := fmt.Sprintf("SELECT COUNT(*) FROM business_promotional_posts WHERE created_at > current_date - interval '1 %s'", period)
	err = database.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, businessPromotionalPostCount{Period: period, Count: count})
}

func getBusinessPromotionalPosts(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, expiration, created_at, updated_at, archived, partner, completed FROM business_promotional_posts ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var bpPosts []bpPost
	for rows.Next() {
		var bpp bpPost
		err := rows.Scan(&bpp.ID, &bpp.TitleLatin, &bpp.DescriptionLatin, &bpp.TitleCyrillic, &bpp.DescriptionCyrillic, pq.Array(&bpp.Videos), &bpp.Expiration, &bpp.CreatedAt, &bpp.UpdatedAt, &bpp.Archived, &bpp.Partner, &bpp.Completed)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		bpPosts = append(bpPosts, bpp)
	}

	response.Res(w, "success", http.StatusOK, bpPosts)
}

type bpPost struct {
	ID                  int      `json:"id"`
	TitleLatin          string   `json:"title_latin"`
	DescriptionLatin    string   `json:"description_latin"`
	TitleCyrillic       string   `json:"title_cyrillic"`
	DescriptionCyrillic string   `json:"description_cyrillic"`
	Videos              []string `json:"videos"`
	Expiration          string   `json:"expiration"`
	CreatedAt           string   `json:"created_at"`
	UpdatedAt           string   `json:"updated_at"`
	Archived            bool     `json:"archived"`
	Partner             string   `json:"partner"`
	Completed           bool     `json:"completed"`
}

// businessPromotionalPostCompleted is a handler to make business promotional post completed field true/false
func businessPromotionalPostCompleted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: businessPromotionalPostCompleted bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: businessPromotionalPostCompleted bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot update completed field of non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: businessPromotionalPostCompleted bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: businessPromotionalPostCompleted bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot update completed field of archived business promotional post")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: businessPromotionalPostCompleted db connection error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE business_promotional_posts SET completed = NOT completed WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: businessPromotionalPostCompleted db execution error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "completed field updated")
}

// addBusinessPromotionalPostPhoto is a handler to add photo to business promotional post
func addBusinessPromotionalPostPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: addBusinessPromotionalPostPhoto bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: addBusinessPromotionalPostPhoto bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot add photo to non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: addBusinessPromotionalPostPhoto bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: addBusinessPromotionalPostPhoto bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot add photo to archived business promotional post")
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(100 << 20) // maxMemory is 100MB
	if err != nil {
		log.Printf("%v: addBusinessPromotionalPostPhoto: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		log.Printf("%v: addBusinessPromotionalPostPhoto: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	photoArray := r.MultipartForm.File["photo"]
	for _, fh := range photoArray {
		if fh.Size > 10<<20 {
			err := fmt.Errorf("photo %v size exceeds 10MB limit: %v", fh.Filename, fh.Size)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		// check if file is image using http.DetectContentType
		contentType := fh.Header.Get("Content-Type")
		if contentType[:5] != "image" {
			err := fmt.Errorf("photo %v is not an image: %v", fh.Filename, contentType)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		file, err := fh.Open()
		if err != nil {
			err := fmt.Errorf("error opening photo %v: %v", fh.Filename, err)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}
		defer file.Close()
		photo, err := io.ReadAll(file)
		if err != nil {
			err := fmt.Errorf("error reading photo %v: %v", fh.Filename, err)
			toolkit.LogError(r, err)
			response.Res(w, "error", http.StatusBadRequest, err.Error())
			return
		}

		// Start a new transaction
		tx, err := db.Begin()
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("start a new transaction: %v", err))
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		// Prepare the SQL statement: add into bpp_photos: bpp, file_name, file
		stmt, err := tx.Prepare("INSERT INTO bpp_photos (bpp, file_name, file) VALUES ($1, $2, $3)")
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("prepare the SQL statement: %v", err))
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		// Execute the SQL statement
		_, err = stmt.Exec(id, fh.Filename, photo)
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("execute the SQL statement: %v", err))
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("commit the transaction: %v", err))
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}

		// update business promotional post updated_at
		_, err = db.Exec("UPDATE business_promotional_posts SET updated_at = NOW() WHERE id = $1", id)
		if err != nil {
			toolkit.LogError(r, fmt.Errorf("update business promotional post updated_at: %v", err))
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusCreated, "photos added")
}

// getBusinessPromotionalPostPhotoList is a handler to get list of photos of business promotional post
func getBusinessPromotionalPostPhotoList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: getBusinessPromotionalPostPhotoList bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: getBusinessPromotionalPostPhotoList bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot get photo list of non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: getBusinessPromotionalPostPhotoList bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: getBusinessPromotionalPostPhotoList bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot get photo list of archived business promotional post")
		return
	}

	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		log.Printf("%v: getBusinessPromotionalPostPhotoList: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	// select id, bpp, file_name, created_at order by id
	rows, err := db.Query("SELECT id, bpp, file_name, created_at FROM bpp_photos WHERE bpp = $1 ORDER BY id", id)
	if err != nil {
		log.Printf("%v: getBusinessPromotionalPostPhotoList: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var photos []model.BusinessPromotionalPostPhoto
	for rows.Next() {
		var photo model.BusinessPromotionalPostPhoto
		err := rows.Scan(&photo.ID, &photo.BPP, &photo.FileName, &photo.CreatedAt)
		if err != nil {
			log.Printf("%v: getBusinessPromotionalPostPhotoList: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		photos = append(photos, photo)
	}

	response.Res(w, "success", http.StatusOK, photos)
}

// getBusinessPromotionalPostPhoto is a handler to get photo of business promotional post
func getBusinessPromotionalPostPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: getBusinessPromotionalPostPhoto bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: getBusinessPromotionalPostPhoto bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot get photo of non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: getBusinessPromotionalPostPhoto bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: getBusinessPromotionalPostPhoto bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot get photo of archived business promotional post")
		return
	}

	// Open a connection to the database
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: getBusinessPromotionalPostPhoto: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	photo_id := vars["photo_id"]
	var photo model.BusinessPromotionalPostPhoto
	err = database.QueryRow("SELECT id, bpp, file_name, file, created_at FROM bpp_photos WHERE id = $1 AND bpp = $2", photo_id, id).Scan(&photo.ID, &photo.BPP, &photo.FileName, &photo.File, &photo.CreatedAt)
	if err != nil {
		log.Printf("%v: getBusinessPromotionalPostPhoto: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	imageType := http.DetectContentType(photo.File)
	w.Header().Set("Content-Type", imageType)
	w.Header().Set("Content-Length", strconv.Itoa(len(photo.File)))
	w.Write(photo.File)
}

// deleteBusinessPromotionalPostPhoto is a handler to delete photo of business promotional post
func deleteBusinessPromotionalPostPhoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := bpPostExists(id)
	if err != nil {
		log.Printf("%v: deleteBusinessPromotionalPostPhoto bpPostExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: deleteBusinessPromotionalPostPhoto bpPostExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete photo of non existent business promotional post")
		return
	}

	archived, err := bpPostIsArchived(id)
	if err != nil {
		log.Printf("%v: deleteBusinessPromotionalPostPhoto bpPostIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: deleteBusinessPromotionalPostPhoto bpPostIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete photo of archived business promotional post")
		return
	}

	photo_id := vars["photo_id"]

	// Open a connection to the database
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: deleteBusinessPromotionalPostPhoto: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	stmt, err := database.Prepare("DELETE FROM bpp_photos WHERE id=$1 AND bpp=$2")
	if err != nil {
		log.Printf("%v: deleteBusinessPromotionalPostPhoto: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	_, err = stmt.Exec(photo_id, id)
	if err != nil {
		log.Printf("%v: deleteBusinessPromotionalPostPhoto: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "deleted")
}
