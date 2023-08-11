package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// AllReviews fetches all the records from reviews db table.
func (m *postgresDBRepo) AllReviews() ([]*models.Review, error) {
	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the sql statement
	query := `SELECT * FROM reviews`

	// Execute the query using Query Context.
	// If any error occurs, nil and error is returned
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Create a slice of Review model
	reviews := []*models.Review{}

	// Loop through the record.
	for rows.Next() {

		// Store the record in Review instance
		review := new(models.Review)
		if err := rows.Scan(
			&review.ID,
			&review.Rating,
			&review.Body,
			&review.BookID,
			&review.UserID,
			&review.IsActive,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Append the Review instance to the slice of Review
		reviews = append(reviews, review)
	}

	// Return reviews
	return reviews, nil
}

// ReviewExists return true if Review book, review and user  exists else return false.
// It takes Review model instance as parameters
func (m *postgresDBRepo) ReviewExists(u *models.Review) (bool, error) {

	// Creating a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the sql query to check for existing relationship
	query := `
		SELECT COUNT(*) FROM reviews
		WHERE (id=$1 AND book_id=$2 AND user_id=$3)
	`

	// intializing a count variable that stores the no of records
	var count int

	// Executing the query row context and store the total record in count variable.
	// If any error occurs, false and error are returned
	if err := m.DB.QueryRowContext(ctx, query, u.ID, u.BookID, u.UserID).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	// returning exists is true if count > 0 else retuirn false
	return count > 0, nil
}

// InsertReview add new book user review relation table to db
// Takes Review model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertReview(u *models.Review) error {

	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare a insert query statement
	stmt := `
		INSERT INTO reviews (rating, body, book_id, user_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	// Executing the query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.Rating,
		u.Body,
		u.BookID,
		u.UserID,
		u.IsActive,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetReviewByID returns the Review detail from database using id.
// It takes review id as parameters.
// Returns a Review struct instance.
func (m *postgresDBRepo) GetReviewByID(id int) (*models.Review, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM reviews
		WHERE id=$1
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, id)

	// Initializing a Review struct instance
	review := &models.Review{}

	// Scannin the row and storing the result in Review Intance.
	if err := row.Scan(
		&review.ID,
		&review.Rating,
		&review.Body,
		&review.BookID,
		&review.UserID,
		&review.IsActive,
		&review.CreatedAt,
		&review.UpdatedAt,
	); err != nil {
		return nil, err
	}

	// Return a Review Instance and nil
	return review, nil
}

// GetReviewByUserID returns the Review detail from database using user id.
// It takes user id as parameters.
// Returns a Review struct instance.
func (m *postgresDBRepo) GetReviewByUserID(id int) (*models.Review, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM reviews
		WHERE user_id=$1
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, id)

	// Initializing a Review struct instance
	review := &models.Review{}

	// Scannin the row and storing the result in Review Intance.
	if err := row.Scan(
		&review.ID,
		&review.Rating,
		&review.Body,
		&review.BookID,
		&review.UserID,
		&review.IsActive,
		&review.CreatedAt,
		&review.UpdatedAt,
	); err != nil {
		return nil, err
	}

	// Return a Review Instance and nil
	return review, nil
}

// DeleteReview deletes the record of Review table from the db.
// It takes book id and user id as parameter
func (m *postgresDBRepo) DeleteReview(id int) error {

	// Using context with timeout of 3 second
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the delete sql statment
	stmt := `DELETE FROM reviews WHERE (id=$1)`

	// executing the query
	_, err := m.DB.ExecContext(ctx, stmt, id)

	// returns nil if success else returns error
	return err
}

// UpdateReview updates the Review
// Takes update value Review model and id of review to be updated as paramaters
func (m *postgresDBRepo) UpdateReview(u *models.Review) error {

	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the query statement for update follower
	stmt := `
		UPDATE reviews
		SET rating = $2, body = $3, book_id = $4, user_id = $5, is_active = $6, updated_at = $7
		WHERE id = $1
	`

	// Executing the sql query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.Rating,
		u.Body,
		u.BookID,
		u.UserID,
		u.IsActive,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetReviewsByBookID(bookID int) ([]*models.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM reviews WHERE book_id=$1 ORDER BY id"
	rows, err := m.DB.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	reviews := []*models.Review{}
	for rows.Next() {
		review := &models.Review{}
		if err := rows.Scan(
			&review.ID,
			&review.Rating,
			&review.Body,
			&review.BookID,
			&review.UserID,
			&review.IsActive,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}
