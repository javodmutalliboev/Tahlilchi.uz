package model

import (
	"database/sql"

	"Tahlilchi.uz/db"
)

// ENewspaperCommentListResponse is a struct to map the e-newspaper comment list response
// It has the following fields:
// ENewspaperCommentList: a slice of ENewspaperComment representing the e-newspaper comment list
// Previous: a boolean representing if there is a previous page
// Next: a boolean representing if there is a next page
type ENewspaperCommentListResponse struct {
	ENewspaperCommentList []ENewspaperComment `json:"e_newspaper_comment_list"`
	Previous              bool                `json:"previous"`
	Next                  bool                `json:"next"`
}

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
	ENewspaper int    `json:"e_newspaper"`
	Text       string `json:"text" validate:"required"`
	Contact    string `json:"contact"`
	CreatedAt  string `json:"created_at"`
	Approved   bool   `json:"approved"`
}

// AddENewspaperComment is a method to add an e-newspaper comment to the database
// It takes no parameters
// It returns an error if any
func (enc *ENewspaperComment) AddENewspaperComment(e_newspaper int) error {
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
	err = stmt.QueryRow(e_newspaper, enc.Text).Scan(&enc.ID)
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

// GetENewspaperCommentListResponse is a method to get the e-newspaper comment list from the database
// It takes a boolean representing if the user is an admin, an int representing the id of the e-newspaper, an int representing the page and an int representing the limit
// It returns an ENewspaperCommentListResponse and an error if any
func (enc *ENewspaperComment) GetENewspaperCommentListResponse(admin bool, id, page, limit int) (ENewspaperCommentListResponse, error) {
	// create a new database connection
	database, err := db.DB()
	if err != nil {
		return ENewspaperCommentListResponse{}, err
	}
	// defer the close of the database connection
	defer database.Close()

	// declare rows variable using *sql.Rows
	var rows *sql.Rows
	// check if the admin is true
	if admin {
		// get the e-newspaper comment list from the database
		rows, err = database.Query("SELECT id, e_newspaper, text, contact, created_at, approved FROM e_newspaper_comments WHERE e_newspaper = $1 ORDER BY id DESC LIMIT $2 OFFSET $3", id, limit, (page-1)*limit)
	} else {
		// get the e-newspaper comment list from the database
		rows, err = database.Query("SELECT id, text, created_at FROM e_newspaper_comments WHERE e_newspaper = $1 AND approved = true ORDER BY id DESC LIMIT $2 OFFSET $3", id, limit, (page-1)*limit)
	}
	// check if there is an error
	if err != nil {
		return ENewspaperCommentListResponse{}, err
	}
	// defer the close of the rows
	defer rows.Close()

	// create a new e-newspaper comment slice using the ENewspaperComment struct
	var encList []ENewspaperComment
	// iterate through the rows
	for rows.Next() {
		// create a new e-newspaper comment
		var enc ENewspaperComment
		// check if the admin is true
		if admin {
			// scan the e-newspaper comment into the e-newspaper comment struct
			err = rows.Scan(&enc.ID, &enc.ENewspaper, &enc.Text, &enc.Contact, &enc.CreatedAt, &enc.Approved)
		} else {
			// scan the e-newspaper comment into the e-newspaper comment struct
			err = rows.Scan(&enc.ID, &enc.Text, &enc.CreatedAt)
		}
		// check if there is an error
		if err != nil {
			return ENewspaperCommentListResponse{}, err
		}
		// append the e-newspaper comment to the e-newspaper comment slice
		encList = append(encList, enc)
	}

	// check error from rows
	if err := rows.Err(); err != nil {
		return ENewspaperCommentListResponse{}, err
	}

	// create a new e-newspaper comment list response
	var encListRes ENewspaperCommentListResponse

	// declare a variable to store the count of the e-newspaper comments
	var count int
	// check if the admin is true
	if admin {
		// get the count of the e-newspaper comments
		err = database.QueryRow("SELECT COUNT(id) FROM e_newspaper_comments WHERE e_newspaper = $1", id).Scan(&count)
	} else {
		// get the count of the e-newspaper comments
		err = database.QueryRow("SELECT COUNT(id) FROM e_newspaper_comments WHERE e_newspaper = $1 AND approved = true", id).Scan(&count)
	}
	// check if there is an error
	if err != nil {
		return ENewspaperCommentListResponse{}, err
	}

	// check if the page is greater than 1
	if page > 1 {
		// set the previous to true
		encListRes.Previous = true
	}
	// check if the count is greater than the page multiplied by the limit
	if count > page*limit {
		// set the next to true
		encListRes.Next = true
	}

	// set the e-newspaper comment list to the e-newspaper comment list response
	encListRes.ENewspaperCommentList = encList

	// return the e-newspaper comment list response and nil
	return encListRes, nil
}

// ApproveENewspaperComment is a method to approve/disapprove an e-newspaper comment in the database
// It takes an int representing the id of the e-newspaper and an int representing the id of the e-newspaper comment
// It returns an error if any
func (enc *ENewspaperComment) ApproveENewspaperComment(e_newspaper, commentID int) error {
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
	// toggle approved column
	_, err = tx.Exec("UPDATE e_newspaper_comments SET approved = NOT approved WHERE e_newspaper = $1 AND id = $2", e_newspaper, commentID)
	if err != nil {
		return err
	}
	// commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	// return nil
	return nil
}
