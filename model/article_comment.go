package model

import (
	"database/sql"

	"Tahlilchi.uz/db"
)

// ArticleCommentListResponse is a struct to map the article comment list response
type ArticleCommentListResponse struct {
	ArticleCommentList []ArticleComment `json:"article_comment_list"`
	Previous           bool             `json:"previous"`
	Next               bool             `json:"next"`
}

// ArticleComment is a struct to map the article comment data
type ArticleComment struct {
	ID        int    `json:"id"`
	Article   int    `json:"article" validate:"required"`
	Text      string `json:"text" validate:"required"`
	Contact   string `json:"contact"`
	CreatedAt string `json:"created_at"`
	Approved  bool   `json:"approved"`
}

// AddArticleComment is a method to add an article comment to the database
func (ac *ArticleComment) AddArticleComment() error {
	// create a new database connection
	database, err := db.DB()
	if err != nil {
		return err
	}
	// defer the close of the database connection
	defer database.Close()

	// create a new transaction
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	// prepare the insert statement it should return the id of the inserted row
	stmt, err := tx.Prepare("INSERT INTO article_comments (article, text) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return err
	}
	// defer the close of the statement
	defer stmt.Close()

	// execute the insert statement
	err = stmt.QueryRow(ac.Article, ac.Text).Scan(&ac.ID)
	if err != nil {
		return err
	}

	// check if the contact is not empty
	if ac.Contact != "" {
		// update the contact of the article comment
		_, err = tx.Exec("UPDATE article_comments SET contact = $1 WHERE id = $2", ac.Contact, ac.ID)
		if err != nil {
			return err
		}
	}
	// commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// return nil
	return nil
}

// GetArticleCommentList is a method to get the article comment list response from the database by article id, page and limit
func (ac *ArticleComment) GetArticleCommentList(admin bool, id, page, limit int) (ArticleCommentListResponse, error) {
	// create a new database connection
	database, err := db.DB()
	if err != nil {
		return ArticleCommentListResponse{}, err
	}
	// defer the close of the database connection
	defer database.Close()

	// declare rows variable using *sql.Rows
	var rows *sql.Rows
	// check if the admin is true
	if admin {
		// get the article comment list from the database
		rows, err = database.Query("SELECT id, article, text, contact, created_at, approved FROM article_comments WHERE article = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", id, limit, (page-1)*limit)
	} else {
		// get the article comment list from the database
		rows, err = database.Query("SELECT id, text, created_at FROM article_comments WHERE article = $1 AND approved = true ORDER BY id DESC LIMIT $2 OFFSET $3", id, limit, (page-1)*limit)
	}
	// check if there is an error
	if err != nil {
		return ArticleCommentListResponse{}, err
	}
	// defer the close of the rows
	defer rows.Close()

	// create a ner article comment slice using the ArticleComment struct
	var acs []ArticleComment
	// iterate through the rows
	for rows.Next() {
		// create a new article comment
		var ac ArticleComment
		// check if the admin is true
		if admin {
			// scan the article comment data from the rows
			err = rows.Scan(&ac.ID, &ac.Article, &ac.Text, &ac.Contact, &ac.CreatedAt, &ac.Approved)
		} else {
			// scan the article comment data from the rows
			err = rows.Scan(&ac.ID, &ac.Text, &ac.CreatedAt)
		}
		// check if there is an error
		if err != nil {
			return ArticleCommentListResponse{}, err
		}
		// append the article comment to the article comment slice
		acs = append(acs, ac)
	}

	// check error from rows
	if err := rows.Err(); err != nil {
		return ArticleCommentListResponse{}, err
	}

	// create a new article comment list response
	var acr ArticleCommentListResponse

	// declare a variable to store the count of the article comments
	var count int
	// check if the admin is true
	if admin {
		// get the count of the article comments from the database
		err = database.QueryRow("SELECT COUNT(*) FROM article_comments WHERE article = $1", id).Scan(&count)
	} else {
		// get the count of the article comments from the database
		err = database.QueryRow("SELECT COUNT(*) FROM article_comments WHERE article = $1 AND approved = true", id).Scan(&count)
	}
	// check if there is an error
	if err != nil {
		return ArticleCommentListResponse{}, err
	}

	// check if the page is greater than 1
	if page > 1 {
		// set the previous to true
		acr.Previous = true
	}
	// check if the count is greater than the page multiplied by the limit
	if count > page*limit {
		// set the next to true
		acr.Next = true
	}

	// set the article comment list to the article comment list response
	acr.ArticleCommentList = acs

	// return the article comment list response
	return acr, nil
}

// ApproveArticleComment is a method to approve an article comment in the database
func (ac *ArticleComment) ApproveArticleComment(id, commentID int) error {
	// create a new database connection
	database, err := db.DB()
	if err != nil {
		return err
	}
	// defer the close of the database connection
	defer database.Close()

	// create a new transaction
	tx, err := database.Begin()
	if err != nil {
		return err
	}
	// update the article comment in the database
	_, err = tx.Exec("UPDATE article_comments SET approved = true WHERE article = $1 AND id = $2", id, commentID)
	if err != nil {
		return err
	}
	// commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	// return nil
	return nil
}
