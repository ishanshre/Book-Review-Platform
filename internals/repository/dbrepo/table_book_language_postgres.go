package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// AllBookLanguage retrieves all book language relationships from the PostgreSQL database.
// It returns a slice of BookLanguage struct and error
func (m *postgresDBRepo) AllBookLanguage() ([]*models.BookLanguage, error) {

	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the sql statement to select all book language relationship
	query := `SELECT * FROM book_languages`

	// Exectue the query and get the result row
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Initialize an empty slice of BookLanguage model
	bookLanguages := []*models.BookLanguage{}

	// Iterate through the rows
	for rows.Next() {

		// create a new instance of BookLanguage
		bookLanguage := new(models.BookLanguage)

		// Scan the value from current row and store in bookLanguage instance
		if err := rows.Scan(
			&bookLanguage.BookID,
			&bookLanguage.LanguageID,
		); err != nil {
			return nil, err
		}

		// append the bookLanguage to the slice
		bookLanguages = append(bookLanguages, bookLanguage)
	}

	// Return the retrieved book languages with no error
	return bookLanguages, nil
}

// DeleteBookLanguage deletes the record of Book Language table from the db.
// It takes book id and language id as parameter
func (m *postgresDBRepo) DeleteBookLanguage(book_id, language_id int) error {

	// Using context with timeout of 3 second
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the delete sql statment
	stmt := `DELETE FROM book_languages WHERE (book_id=$1 AND language_id=$2)`

	// executing the query
	_, err := m.DB.ExecContext(ctx, stmt, book_id, language_id)

	// returns nil if success else returns error
	return err
}

// GetBookLanguageByID returns the book from database using id.
// It takes book id and language id as parameters.
// Returns a BookLanguage struct instance.
func (m *postgresDBRepo) GetBookLanguageByID(book_id, language_id int) (*models.BookLanguage, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM book_languages
		WHERE (book_id=$1 AND language_id=$2)
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, book_id, language_id)

	// Initializing a BookLanguage struct instance
	bookLanguage := &models.BookLanguage{}

	// Scannin the row and storing the result in BookLanguage Intance.
	if err := row.Scan(
		&bookLanguage.BookID,
		&bookLanguage.LanguageID,
	); err != nil {
		return nil, err
	}

	// Return a BookLanguage Instance and nil
	return bookLanguage, nil
}

// BookLanguageExists return true if book Language relation exists else return false.
// It takes book id and language id as parameters
func (m *postgresDBRepo) BookLanguageExists(book_id, language_id int) (bool, error) {

	// Creating a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the sql query to check for existing relationship
	query := `
		SELECT COUNT(*) FROM book_languages
		WHERE (book_id=$1 AND language_id=$2)
	`

	// intializing a count variable that stores the no of records
	var count int

	// Executing the query row context and store the total record in count variable.
	// If any error occurs, false and error are returned
	if err := m.DB.QueryRowContext(ctx, query, book_id, language_id).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	// returning exists is true if count > 0 else retuirn false
	return count > 0, nil
}

// UpdateBookLanguage updates the book Language
// Takes update value BookLanguage model and previous book_id , language_id
func (m *postgresDBRepo) UpdateBookLanguage(u *models.BookLanguage, book_id, language_id int) error {

	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the query statement for update book language relationship
	stmt := `
		UPDATE book_languages
		SET book_id = $3, language_id = $4
		WHERE (book_id = $1 AND language_id = $2)
	`

	// Executing the sql query
	_, err := m.DB.ExecContext(ctx, stmt, book_id, language_id, u.BookID, u.LanguageID)
	if err != nil {
		return err
	}
	return nil
}

// InsertBookLanguage add new book Language to db
// Takes BookLanguage model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertBookLanguage(u *models.BookLanguage) error {

	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare a insert query statement
	stmt := `
		INSERT INTO book_languages (book_id, language_id)
		VALUES ($1, $2);
	`

	// Executing the query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.BookID,
		u.LanguageID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetLanguagesFromBookID(book_id int) ([]*models.Language, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT 
		COALESCE(l.id, 0), COALESCE(l.language, '')
		FROM
			book_languages AS bl
		LEFT JOIN
			languages AS l ON l.id = bl.language_id
		WHERE 
			book_id = $1
	`
	languages := []*models.Language{}
	rows, err := m.DB.QueryContext(ctx, query, book_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		language := &models.Language{}
		if err := rows.Scan(
			&language.ID,
			&language.Language,
		); err != nil {
			return nil, err
		}
		languages = append(languages, language)
	}
	return languages, nil
}

func (m *postgresDBRepo) GetAllBooksByLanguage(limit, page int, searchKey, sort, language string) (*models.BookApiFilter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := `
		SELECT 
			COALESCE(b.id, 0) AS b_id,
			COALESCE(b.title, '') AS b_title,
			COALESCE(b.isbn, 0) AS b_isbn,
			COALESCE(b.cover, '') AS b_cover
		FROM 
			book_languages AS bl
		LEFT JOIN 
			books AS b ON b.id = bl.book_id
		LEFT JOIN
			languages AS l ON l.id = bl.language_id
		WHERE
			l.language = $1
	`
	countQuery := `
		SELECT 
			COUNT(*)
		FROM 
			book_languages AS bl
		LEFT JOIN 
			books AS b ON b.id = bl.book_id
		LEFT JOIN
			languages AS l ON l.id = bl.language_id
		WHERE
			l.language = $1
	`
	if searchKey != "" {
		query = fmt.Sprintf("%s AND (b.title LIKE '%%%s%%' OR CAST(b.isbn AS TEXT) LIKE '%%%s%%')", query, searchKey, searchKey)
		countQuery = fmt.Sprintf("%s AND (b.title LIKE '%%%s%%' OR CAST(b.isbn AS TEXT) LIKE '%%%s%%')", countQuery, searchKey, searchKey)
	}
	if sort != "" {
		query = fmt.Sprintf("%s ORDER BY b.title %s", query, sort)
	}

	var count int
	if err := m.DB.QueryRowContext(ctx, countQuery, language).Scan(&count); err != nil {
		return nil, err
	}

	query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
	rows, err := m.DB.QueryContext(ctx, query, language)
	if err != nil {
		return nil, err
	}
	books := []*models.Book{}
	for rows.Next() {
		book := &models.Book{}
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Isbn,
			&book.Cover,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	last_page := m.CalculateLastPage(limit, count)
	return &models.BookApiFilter{
		Total:    count,
		LastPage: last_page,
		Page:     page,
		Books:    books,
	}, nil
}

func (m *postgresDBRepo) TotalLanguageCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT COUNT(*) FROM languages;`
	var count int
	if err := m.DB.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
