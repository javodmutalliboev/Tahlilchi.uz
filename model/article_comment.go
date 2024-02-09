package model

import "Tahlilchi.uz/db"

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
