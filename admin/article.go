package admin

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"Tahlilchi.uz/db"
	"Tahlilchi.uz/response"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

// ArticleCategory is a struct to hold article category data
type ArticleCategory struct {
	ID                  int    `json:"id"`
	TitleLatin          string `json:"title_latin"`
	DescriptionLatin    string `json:"description_latin"`
	TitleCyrillic       string `json:"title_cyrillic"`
	DescriptionCyrillic string `json:"description_cyrillic"`
}

// getArticleCategory is a handler function to get article categories from the database
func getArticleCategory(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	rows, err := database.Query("select * from article_category order by id")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	var articleCategoryList []ArticleCategory
	for rows.Next() {
		var a ArticleCategory
		err := rows.Scan(&a.ID, &a.TitleLatin, &a.DescriptionLatin, &a.TitleCyrillic, &a.DescriptionCyrillic)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		articleCategoryList = append(articleCategoryList, a)
	}

	response.Res(w, "success", http.StatusOK, articleCategoryList)
}

// updateArticleCategory is a handler function to update an article category in the database
func updateArticleCategory(w http.ResponseWriter, r *http.Request) {
	// get id from url params
	vars := mux.Vars(r)
	id := vars["id"]

	// check if the category exists
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM article_category WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !exists {
		response.Res(w, "error", http.StatusBadRequest, "Category does not exist")
		return
	}

	// Parse the form
	err = r.ParseForm()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, "form parse error")
		return
	}

	// title_latin
	titleLatin := r.FormValue("title_latin")
	if titleLatin != "" {
		sqlStatement := `
			UPDATE article_category
			SET title_latin = $1
			WHERE id = $2;
		`
		_, err = database.Exec(sqlStatement, titleLatin, id)
		if err != nil {
			log.Printf("%v: writing title_latin into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// description_latin
	descriptionLatin := r.FormValue("description_latin")
	if descriptionLatin != "" {
		sqlStatement := `
			UPDATE article_category
			SET description_latin = $1
			WHERE id = $2;
		`
		_, err = database.Exec(sqlStatement, descriptionLatin, id)
		if err != nil {
			log.Printf("%v: writing description_latin into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// title_cyrillic
	titleCyrillic := r.FormValue("title_cyrillic")
	if titleCyrillic != "" {
		sqlStatement := `
			UPDATE article_category
			SET title_cyrillic = $1
			WHERE id = $2;
		`
		_, err = database.Exec(sqlStatement, titleCyrillic, id)
		if err != nil {
			log.Printf("%v: writing title_cyrillic into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// description_cyrillic
	descriptionCyrillic := r.FormValue("description_cyrillic")
	if descriptionCyrillic != "" {
		sqlStatement := `
			UPDATE article_category
			SET description_cyrillic = $1
			WHERE id = $2;
		`
		_, err = database.Exec(sqlStatement, descriptionCyrillic, id)
		if err != nil {
			log.Printf("%v: writing description_cyrillic into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	response.Res(w, "success", http.StatusOK, "Category updated")
}

// deleteArticleCategory is a handler function to delete an article category from the database
func deleteArticleCategory(w http.ResponseWriter, r *http.Request) {
	// get id from url params
	vars := mux.Vars(r)
	id := vars["id"]

	// check if the category exists
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var exists bool
	err = database.QueryRow("SELECT EXISTS(SELECT 1 FROM article_category WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !exists {
		response.Res(w, "error", http.StatusBadRequest, "Category does not exist")
		return
	}

	// Prepare the SQL statement: delete from article_category where id = $1
	stmt, err := database.Prepare("DELETE FROM article_category WHERE id = $1")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "Category deleted")
}

// addArticleCategory is a handler function to add a new article category to the database
func addArticleCategory(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.Exec("INSERT INTO article_category(title_latin, description_latin, title_cyrillic, description_cyrillic) VALUES($1, $2, $3, $4)", c.TitleLatin, c.DescriptionLatin, c.TitleCyrillic, c.DescriptionCyrillic)
	if err != nil {
		log.Println(err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusCreated, "Category Added")
}

// ArticlePhoto is a struct to hold article photo data
type ArticlePhoto struct {
	ID        int    `json:"id"`
	Article   int    `json:"article"`
	FileName  string `json:"file_name"`
	File      []byte `json:"file"`
	CreatedAt string `json:"created_at"`
}

// addArticle is a handler function to add a new article to the database
func addArticle(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(200 << 20)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, "file too large")
		return
	}

	// title_latin
	titleLatin := r.FormValue("title_latin")
	if titleLatin == "" {
		response.Res(w, "error", http.StatusBadRequest, "title_latin is required")
		return
	}

	// description_latin
	descriptionLatin := r.FormValue("description_latin")

	// title_cyrillic
	titleCyrillic := r.FormValue("title_cyrillic")
	if titleCyrillic == "" {
		response.Res(w, "error", http.StatusBadRequest, "title_cyrillic is required")
		return
	}

	// description_cyrillic
	descriptionCyrillic := r.FormValue("description_cyrillic")

	// photos
	photoFiles := r.MultipartForm.File["photo"]
	var photos []ArticlePhoto
	if len(photoFiles) > 0 {
		for _, fh := range photoFiles {
			// check whether the file is an image by checking the content type
			if fh.Header.Get("Content-Type")[:5] != "image" {
				log.Printf("%v: photo is not an image: %v", r.URL, fh.Header.Get("Content-Type"))
				response.Res(w, "error", http.StatusBadRequest, "photo is not an image")
				return
			}
			// Check if the file size is greater than 10MB
			if fh.Size > 10<<20 {
				log.Printf("%v: photo size exceeds 10MB limit: %v", r.URL, fh.Size)
				response.Res(w, "error", http.StatusBadRequest, "photo size exceeds 10MB limit")
				return
			}
			// Read the file
			file, err := fh.Open()
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
			// create a variable of type ArticlePhoto
			var photo ArticlePhoto
			photo.FileName = fh.Filename
			// Read the file into a byte slice
			photo.File, err = io.ReadAll(file)
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
			// Append the byte slice to the photos slice
			photos = append(photos, photo)
			err = file.Close()
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
		}
	}

	// videos
	videos := r.Form["video"]
	// videosString := "{" + strings.Join(videos, ",") + "}"

	// cover_image
	coverImageFile, coverImageHeader, err := r.FormFile("cover_image")
	// cover_image is []byte
	var coverImage []byte
	if err != nil && err != http.ErrMissingFile {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, "cover_image error")
		return
	} else if err == http.ErrMissingFile {
		coverImage = nil
	} else {
		// check whether the file is an image by checking the content type
		if coverImageHeader.Header.Get("Content-Type")[:5] != "image" {
			log.Printf("%v: cover_image is not an image: %v", r.URL, coverImageHeader.Header.Get("Content-Type"))
			response.Res(w, "error", http.StatusBadRequest, "cover_image is not an image")
			return
		}
		// Check if the file size is greater than 15MB
		if coverImageHeader.Size > 15<<20 {
			log.Printf("%v: cover_image size exceeds 15MB limit: %v", r.URL, coverImageHeader.Size)
			response.Res(w, "error", http.StatusBadRequest, "cover_image size exceeds 15MB limit")
			return
		}
		// Read the file
		coverImage, err = io.ReadAll(coverImageFile)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		// Close the file
		err = coverImageFile.Close()
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	// tags
	tags := r.Form["tag[]"]
	// tagsString := "{" + strings.Join(tags, ",") + "}"

	// category
	categoryStr := r.FormValue("category")
	// category nullable int
	var category *int
	if categoryStr != "" {
		categoryInt, err := strconv.Atoi(categoryStr)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		category = &categoryInt
	}

	// related
	relatedStr := r.FormValue("related")
	// related nullable int
	var related *int
	if relatedStr != "" {
		relatedInt, err := strconv.Atoi(relatedStr)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		related = &relatedInt
	}

	// Open a connection to the database
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// Prepare the SQL statement: insert title_latin, description_latin, title_cyrillic, description_cyrillic, videos, cover_image, tags, category, related into articles return id
	stmt, err := database.Prepare("INSERT INTO articles(title_latin, description_latin, title_cyrillic, description_cyrillic, videos, cover_image, tags, category, related) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	// id is bigint
	var id int64
	err = stmt.QueryRow(titleLatin, descriptionLatin, titleCyrillic, descriptionCyrillic, pq.Array(videos), coverImage, pq.Array(tags), category, related).Scan(&id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// check if len(photos) > 0
	if len(photos) > 0 {
		// Prepare the SQL statement: insert article, file_name, file into article_photos
		stmt, err = database.Prepare("INSERT INTO article_photos(article, file_name, file) VALUES($1, $2, $3)")
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			// delete the article from the articles table
			_, err = database.Exec("DELETE FROM articles WHERE id = $1", id)
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
			}
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		defer stmt.Close()

		// Execute the SQL statement
		for _, photo := range photos {
			_, err = stmt.Exec(id, photo.FileName, photo.File)
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				// delete the article from the articles table
				_, err = database.Exec("DELETE FROM articles WHERE id = $1", id)
				if err != nil {
					log.Printf("%v: error: %v", r.URL, err)
				}
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
		}
	}

	response.Res(w, "success", http.StatusCreated, "Article added")
}

func articleExists(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM articles WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &exists, nil
}

func articleIsArchived(id string) (*bool, error) {
	// Open a connection to the database
	db, err := db.DB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Prepare the SQL statement
	stmt, err := db.Prepare("SELECT archived FROM articles WHERE id = $1")
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

func editArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: edit article articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: edit article articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit non existent article")
		return
	}

	archived, err := articleIsArchived(id)
	if err != nil {
		log.Printf("%v: edit article articleIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: edit article articleIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot edit archived article")
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(200 << 20)
	if err != nil {
		log.Printf("%v: edit article: %v", r.URL, err)
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
			UPDATE articles
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
			UPDATE articles
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
			UPDATE articles
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
			UPDATE articles
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

	if videos, ok := r.Form["videos"]; ok {
		videosString := "{" + strings.Join(videos, ",") + "}"
		sqlStatement := `
			UPDATE articles
			SET videos = $1, updated_at = NOW()
			WHERE id = $2;
		`
		_, err = db.Exec(sqlStatement, videosString, id)
		if err != nil {
			log.Printf("%v: writing videos into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	coverImageFile, coverImageHeader, err := r.FormFile("cover_image")
	if err != nil {
		if err == http.ErrMissingFile {
			coverImageFile = nil
		} else {
			log.Printf("%v: cover_image error: %v", r.URL, err)
			response.Res(w, "error", http.StatusBadRequest, "cover_image error")
			return
		}
	}

	var coverImage []byte = nil

	if coverImageFile != nil {
		// check whether the file is an image by checking the content type
		if coverImageHeader.Header.Get("Content-Type")[:5] != "image" {
			log.Printf("%v: cover_image is not an image: %v", r.URL, coverImageHeader.Header.Get("Content-Type"))
			response.Res(w, "error", http.StatusBadRequest, "cover_image is not an image")
			return
		}
		// Check if the file size is greater than 15MB
		if coverImageHeader.Size > 15<<20 {
			log.Printf("%v: cover_image size exceeds 15MB limit: %v", r.URL, coverImageHeader.Size)
			response.Res(w, "error", http.StatusBadRequest, "cover_image size exceeds 15MB limit")
			return
		}
		// Read the file
		coverImage, err = io.ReadAll(coverImageFile)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		// Close the file
		coverImageFile.Close()
		// Prepare the SQL statement
		sqlStatement := `
			UPDATE articles
			SET cover_image = $1, updated_at = NOW()
			WHERE id = $2;
		`
		// Execute the SQL statement
		_, err = db.Exec(sqlStatement, coverImage, id)
		if err != nil {
			log.Printf("%v: writing cover_image into db: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
	}

	if tags, ok := r.Form["tags"]; ok {
		tagsString := "{" + strings.Join(tags, ",") + "}"
		sqlStatement := `
			UPDATE articles
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
			UPDATE articles
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

	if related := r.FormValue("related"); related != "" {
		relatedInt, err := strconv.ParseInt(related, 10, 64)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		sqlStatement := `
			UPDATE articles
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

	response.Res(w, "success", http.StatusOK, "Article edited")
}

// addArticlePhotos is a handler function to add photos to an article
func addArticlePhotos(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(200 << 20)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusBadRequest, "file too large")
		return
	}

	// Parse the article id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if the article exists
	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: add article photos articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: add article photos articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot add photos to non existent article")
		return
	}

	archived, err := articleIsArchived(id)
	if err != nil {
		log.Printf("%v: add article photos articleIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: add article photos articleIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot add photos to archived article")
		return
	}

	// Get the photos from the multipart form
	photoFiles := r.MultipartForm.File["photo"]
	var photos []ArticlePhoto
	if len(photoFiles) > 0 {
		for _, fh := range photoFiles {
			// check whether the file is an image by checking the content type
			if fh.Header.Get("Content-Type")[:5] != "image" {
				log.Printf("%v: photo is not an image: %v", r.URL, fh.Header.Get("Content-Type"))
				response.Res(w, "error", http.StatusBadRequest, "photo is not an image")
				return
			}
			// Check if the file size is greater than 10MB
			if fh.Size > 10<<20 {
				log.Printf("%v: photo size exceeds 10MB limit: %v", r.URL, fh.Size)
				response.Res(w, "error", http.StatusBadRequest, "photo size exceeds 10MB limit")
				return
			}
			// Read the file
			file, err := fh.Open()
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
			// create a variable of type ArticlePhoto
			var photo ArticlePhoto
			photo.FileName = fh.Filename
			// Read the file into a byte slice
			photo.File, err = io.ReadAll(file)
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
			// close the file
			err = file.Close()
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
			// Append the byte slice to the photos slice
			photos = append(photos, photo)
		}
		// Open a connection to the database
		database, err := db.DB()
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		defer database.Close()
		// Prepare the SQL statement: insert article, file_name, file into article_photos
		stmt, err := database.Prepare("INSERT INTO article_photos(article, file_name, file) VALUES($1, $2, $3)")
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		defer stmt.Close()
		// Execute the SQL statement
		for _, photo := range photos {
			_, err = stmt.Exec(id, photo.FileName, photo.File)
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
		}
	}
	response.Res(w, "success", http.StatusCreated, "Photos added")
}

// getArticlePhotos is a handler function to get photos of an article
func getArticlePhotos(w http.ResponseWriter, r *http.Request) {
	// Parse the article id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if the article exists
	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: get article photos articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: get article photos articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot get photos of non existent article")
		return
	}

	// Open a connection to the database
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// Prepare the SQL statement: select id, article, file_name, created_at from article_photos where article = $1
	rows, err := database.Query("SELECT id, article, file_name, created_at FROM article_photos WHERE article = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	// Create a slice of type ArticlePhoto
	var photos []ArticlePhoto
	// Iterate over the rows
	for rows.Next() {
		// Create a variable of type ArticlePhoto
		var p ArticlePhoto
		// Scan the rows into the variable
		err := rows.Scan(&p.ID, &p.Article, &p.FileName, &p.CreatedAt)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		// Append the variable to the slice
		photos = append(photos, p)
	}
	// Send the response
	response.Res(w, "success", http.StatusOK, photos)
}

// getArticlePhoto is a handler function to get a photo of an article
func getArticlePhoto(w http.ResponseWriter, r *http.Request) {
	// Parse the article id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if the article exists
	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: get article photo articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: get article photo articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot get photo of non existent article")
		return
	}

	// Parse the photo id from the URL
	photoID := vars["photo_id"]

	// Open a connection to the database
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// Prepare the SQL statement: select file from article_photos where article = $1 and id = $2
	stmt, err := database.Prepare("SELECT file FROM article_photos WHERE article = $1 AND id = $2")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	var photo []byte
	err = stmt.QueryRow(id, photoID).Scan(&photo)
	if err != nil {
		// check if the error is no rows in result set
		if err == sql.ErrNoRows {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusNotFound, "photo not found")
			return
		}
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	// Send the response as photo
	contentType := http.DetectContentType(photo)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(photo)))
	w.Write(photo)
}

// deleteArticlePhoto is a handler function to delete a photo of an article
func deleteArticlePhoto(w http.ResponseWriter, r *http.Request) {
	// Parse the article id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if the article exists
	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: delete article photo articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: delete article photo articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete photo of non existent article")
		return
	}

	// Parse the photo id from the URL
	photoID := vars["photo_id"]

	// Open a connection to the database
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// Prepare the SQL statement: delete from article_photos where id = $1
	stmt, err := database.Prepare("DELETE FROM article_photos WHERE id = $1 AND article = $2")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(photoID, id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// update the updated_at column of the articles table
	_, err = database.Exec("UPDATE articles SET updated_at = NOW() WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		// no response is sent to the client
		// do not return
	}

	// Send the response
	response.Res(w, "success", http.StatusOK, "Photo deleted")
}

// getArticleCoverImage is a handler function to get the cover image of an article
func getArticleCoverImage(w http.ResponseWriter, r *http.Request) {
	// Parse the article id from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if the article exists
	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: get article cover image articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: get article cover image articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot get cover image of non existent article")
		return
	}

	// Open a connection to the database
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	// Prepare the SQL statement: select cover_image from articles where id = $1
	stmt, err := database.Prepare("SELECT cover_image FROM articles WHERE id = $1")
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	var coverImage []byte
	err = stmt.QueryRow(id).Scan(&coverImage)
	if err != nil {
		// check if the error is no rows in result set
		if err == sql.ErrNoRows {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusNotFound, "cover image not found")
			return
		}
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	// Send the response as cover image
	contentType := http.DetectContentType(coverImage)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(coverImage)))
	w.Write(coverImage)
}

// deleteArticle is a handler function to delete an article
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: delete article articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: delete article articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete non existent article")
		return
	}

	archived, err := articleIsArchived(id)
	if err != nil {
		log.Printf("%v: delete article articleIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: delete article articleIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot delete archived article")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	// First delete the photos of the article
	_, err = db.Exec("DELETE FROM article_photos WHERE article = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	// Prepare the SQL statement
	stmt, err := db.Prepare("DELETE FROM articles WHERE id=$1")
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

func archiveArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: archive article articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: archive article articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive non existent article")
		return
	}

	archived, err := articleIsArchived(id)
	if err != nil {
		log.Printf("%v: archive article articleIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: archive article articleIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot archive already archived article")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE articles SET archived = true WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "archived")
}

func unArchiveArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: unarchive article articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot unarchive non existent article")
		return
	}

	archived, err := articleIsArchived(id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*archived {
		log.Printf("%v: unarchive article articleIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot unarchive not archived article")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE articles SET archived = false WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "unarchive done")
}

type ArticleCount struct {
	Period string `json:"period"`
	Count  int    `json:"count"`
}

func getArticleCountAll(w http.ResponseWriter, r *http.Request) {
	database, err := db.DB()
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer database.Close()

	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM articles").Scan(&count)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, ArticleCount{Period: "all", Count: count})
}

func getArticleCount(w http.ResponseWriter, r *http.Request) {
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
	query := fmt.Sprintf("SELECT COUNT(*) FROM articles WHERE created_at > current_date - interval '1 %s'", period)
	err = database.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, ArticleCount{Period: period, Count: count})
}

type Article struct {
	ID                  int            `json:"id"`
	TitleLatin          string         `json:"title_latin"`
	DescriptionLatin    string         `json:"description_latin"`
	TitleCyrillic       string         `json:"title_cyrillic"`
	DescriptionCyrillic string         `json:"description_cyrillic"`
	Photos              []ArticlePhoto `json:"photos"`
	Videos              []string       `json:"videos"`
	Tags                []string       `json:"tags"`
	Archived            bool           `json:"archived"`
	CreatedAt           string         `json:"created_at"`
	UpdatedAt           string         `json:"updated_at"`
	Category            *int           `json:"category"`
	Related             *int           `json:"related"`
	Completed           bool           `json:"completed"`
}

func getArticles(w http.ResponseWriter, r *http.Request) {
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

	rows, err := database.Query("SELECT id, title_latin, description_latin, title_cyrillic, description_cyrillic, videos, tags, archived, created_at, updated_at, category, related, completed FROM articles ORDER BY id DESC LIMIT $1 OFFSET $2", limit, start)
	if err != nil {
		log.Printf("%v: error: %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var a Article
		err := rows.Scan(&a.ID, &a.TitleLatin, &a.DescriptionLatin, &a.TitleCyrillic, &a.DescriptionCyrillic, pq.Array(&a.Videos), pq.Array(&a.Tags), &a.Archived, &a.CreatedAt, &a.UpdatedAt, &a.Category, &a.Related, &a.Completed)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		articles = append(articles, a)
	}

	// attach each article its photos
	for i, a := range articles {
		// Prepare the SQL statement: select id, article, file_name, created_at from article_photos where article = $1
		rows, err := database.Query("SELECT id, article, file_name, created_at FROM article_photos WHERE article = $1", a.ID)
		if err != nil {
			log.Printf("%v: error: %v", r.URL, err)
			response.Res(w, "error", http.StatusInternalServerError, "server error")
			return
		}
		defer rows.Close()

		// Create a slice of type ArticlePhoto
		var photos []ArticlePhoto
		// Iterate over the rows
		for rows.Next() {
			// Create a variable of type ArticlePhoto
			var p ArticlePhoto
			// Scan the rows into the variable
			err := rows.Scan(&p.ID, &p.Article, &p.FileName, &p.CreatedAt)
			if err != nil {
				log.Printf("%v: error: %v", r.URL, err)
				response.Res(w, "error", http.StatusInternalServerError, "server error")
				return
			}
			// Append the variable to the slice
			photos = append(photos, p)
		}
		// attach the photos to the article
		articles[i].Photos = photos
	}

	response.Res(w, "success", http.StatusOK, articles)
}

// articleCompleted is a handler to update the completed field of an article
func articleCompleted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	exists, err := articleExists(id)
	if err != nil {
		log.Printf("%v: articleCompleted articleExists(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if !*exists {
		log.Printf("%v: articleCompleted articleExists(id): %v", r.URL, *exists)
		response.Res(w, "error", http.StatusBadRequest, "Cannot update completed field of non existent article")
		return
	}

	archived, err := articleIsArchived(id)
	if err != nil {
		log.Printf("%v: articleCompleted articleIsArchived(id): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	if *archived {
		log.Printf("%v: articleCompleted articleIsArchived(id): %v", r.URL, *archived)
		response.Res(w, "error", http.StatusBadRequest, "Cannot update completed field of archived article")
		return
	}

	db, err := db.DB()
	if err != nil {
		log.Printf("%v: articleCompleted db.DB(): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE articles SET completed = NOT completed WHERE id = $1", id)
	if err != nil {
		log.Printf("%v: articleCompleted db.Exec(): %v", r.URL, err)
		response.Res(w, "error", http.StatusInternalServerError, "server error")
		return
	}

	response.Res(w, "success", http.StatusOK, "completed field updated")
}
