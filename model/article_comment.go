package model

import "Tahlilchi.uz/db"

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
func (ac *ArticleComment) GetArticleCommentList(id, page, limit int) (ArticleCommentListResponse, error) {
	// create a new database connection
	database, err := db.DB()
	if err != nil {
		return ArticleCommentListResponse{}, err
	}
	// defer the close of the database connection
	defer database.Close()

	// get ArticleCommentList from the database
	rows, err := database.Query("SELECT id, text, created_at FROM article_comments WHERE article = $1 AND approved = true ORDER BY id DESC LIMIT $2 OFFSET $3", id, limit, (page-1)*limit)
	if err != nil {
		return ArticleCommentListResponse{}, err
	}
	// defer the close of the rows
	defer rows.Close()

	// create a new ArticleCommentListResponse
	acr := ArticleCommentListResponse{}
	// create a new ArticleCommentList
	acs := []ArticleComment{}
	// iterate over the rows
	for rows.Next() {
		// create a new ArticleComment
		ac := ArticleComment{}
		// scan the rows to the ArticleComment
		err = rows.Scan(&ac.ID, &ac.Text, &ac.CreatedAt)
		if err != nil {
			return ArticleCommentListResponse{}, err
		}
		// append the ArticleComment to the ArticleCommentList
		acs = append(acs, ac)
	}
	// check if there is an error
	if err := rows.Err(); err != nil {
		return ArticleCommentListResponse{}, err
	}

	// get the count of the ArticleCommentList
	var count int
	err = database.QueryRow("SELECT COUNT(*) FROM article_comments WHERE article = $1 AND approved = true", id).Scan(&count)
	if err != nil {
		return ArticleCommentListResponse{}, err
	}

	// check if the page is greater than 1
	if page > 1 {
		acr.Previous = true
	}
	// check if the count is greater than the page*limit
	if count > page*limit {
		acr.Next = true
	}
	// set the ArticleCommentList to the ArticleCommentListResponse
	acr.ArticleCommentList = acs

	// return the ArticleCommentListResponse
	return acr, nil
}
