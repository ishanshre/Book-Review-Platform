package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// AllReadList fetches all the records from readLists db table.
func (m *postgresDBRepo) AllReadList() ([]*models.ReadList, error) {
	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the sql statement
	query := `SELECT * FROM read_lists`

	// Execute the query using Query Context.
	// If any error occurs, nil and error is returned
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Create a slice of ReadList model
	readLists := []*models.ReadList{}

	// Loop through the record.
	for rows.Next() {

		// Store the record in ReadList instance
		readList := new(models.ReadList)
		if err := rows.Scan(&readList.UserID, &readList.BookID, &readList.CreatedAt); err != nil {
			return nil, err
		}

		// Append the ReadList instance to the slice of ReadList
		readLists = append(readLists, readList)
	}

	// Return readLists
	return readLists, nil
}

// ReadListExists return true if ReadList book and user relation exists else return false.
// It takes book id and language id as parameters
func (m *postgresDBRepo) ReadListExists(user_id, book_id int) (bool, error) {

	// Creating a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the sql query to check for existing relationship
	query := `
		SELECT COUNT(*) FROM read_lists
		WHERE (user_id=$1 AND book_id=$2)
	`

	// intializing a count variable that stores the no of records
	var count int

	// Executing the query row context and store the total record in count variable.
	// If any error occurs, false and error are returned
	if err := m.DB.QueryRowContext(ctx, query, user_id, book_id).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	// returning exists is true if count > 0 else retuirn false
	return count > 0, nil
}

// InsertReadList add new book user readlist relation to db
// Takes ReadList model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertReadList(u *models.ReadList) error {

	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare a insert query statement
	stmt := `
		INSERT INTO read_lists (user_id, book_id, created_at)
		VALUES ($1, $2, $3);
	`

	// Executing the query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.UserID,
		u.BookID,
		u.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetReadListByID returns the readlist detail from database using id.
// It takes book id and user id as parameters.
// Returns a ReadList struct instance.
func (m *postgresDBRepo) GetReadListByID(user_id, book_id int) (*models.ReadList, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM read_lists
		WHERE (user_id=$1 AND book_id=$2)
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, user_id, book_id)

	// Initializing a ReadList struct instance
	readList := &models.ReadList{}

	// Scannin the row and storing the result in ReadList Intance.
	if err := row.Scan(
		&readList.UserID,
		&readList.BookID,
		&readList.CreatedAt,
	); err != nil {
		return nil, err
	}

	// Return a ReadList Instance and nil
	return readList, nil
}

// DeleteReadList deletes the record of read_lists table from the db.
// It takes book id and user id as parameter
func (m *postgresDBRepo) DeleteReadList(user_id, book_id int) error {

	// Using context with timeout of 3 second
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the delete sql statment
	stmt := `DELETE FROM read_lists WHERE (user_id=$1 AND book_id=$2)`

	// executing the query
	_, err := m.DB.ExecContext(ctx, stmt, user_id, book_id)

	// returns nil if success else returns error
	return err
}

// UpdateReadList updates the read_lists
// Takes update value ReadList model and previous book_id , user
func (m *postgresDBRepo) UpdateReadList(u *models.ReadList, book_id, user_id int) error {

	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the query statement for update readList
	stmt := `
		UPDATE read_lists
		SET user_id = $3, book_id = $4
		WHERE (user_id = $1 AND book_id = $2)
	`

	// Executing the sql query
	_, err := m.DB.ExecContext(ctx, stmt, user_id, book_id, u.UserID, u.BookID)
	if err != nil {
		return err
	}
	return nil
}
