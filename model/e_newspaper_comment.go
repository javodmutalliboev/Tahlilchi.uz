package model

import "Tahlilchi.uz/db"

// ENewspaperComment is a struct to represent an e-newspaper comment
// It has the following fields:
// ID: an int representing the id of the e-newspaper comment
// ENewspaper: an int representing the id of the e-newspaper
// Text: a string representing the text of the e-newspaper comment
// Contact: a string representing the contact of the e-newspaper comment
// CreatedAt: a timestamp representing the time the e-newspaper comment was created
// Approved: a boolean representing if the e-newspaper comment is approved
type ENewspaperComment struct {
	ID         int    `json:"id"`
	ENewspaper int    `json:"e_newspaper" validate:"required"`
	Text       string `json:"text" validate:"required"`
	Contact    string `json:"contact"`
	CreatedAt  string `json:"created_at"`
	Approved   bool   `json:"approved"`
}

// AddENewspaperComment is a method to add an e-newspaper comment to the database
// It takes no parameters
// It returns an error if any
func (enc *ENewspaperComment) AddENewspaperComment() error {
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
	stmt, err := tx.Prepare("INSERT INTO e_newspaper_comments (e_newspaper, text) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return err
	}
	// defer the close of the statement
	defer stmt.Close()

	// execute the insert statement
	err = stmt.QueryRow(enc.ENewspaper, enc.Text).Scan(&enc.ID)
	if err != nil {
		return err
	}

	// check if the contact is not empty
	if enc.Contact != "" {
		// update the contact of the e-newspaper comment
		_, err = tx.Exec("UPDATE e_newspaper_comments SET contact = $1 WHERE id = $2", enc.Contact, enc.ID)
		if err != nil {
			return err
		}
	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	// return nil
	return nil
}
