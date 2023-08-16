package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// Book Genre db method implementation

// AllBookGenre fetches all record of Book Genre table from database
func (m *postgresDBRepo) AllBookGenre() ([]*models.BookGenre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM book_genres`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	bookGenres := []*models.BookGenre{}
	for rows.Next() {
		bookGenre := new(models.BookGenre)
		if err := rows.Scan(
			&bookGenre.BookID,
			&bookGenre.GenreID,
		); err != nil {
			return nil, err
		}
		bookGenres = append(bookGenres, bookGenre)
	}
	return bookGenres, nil
}

// DeleteBookGenre deletes the record of Book genre table from the db
func (m *postgresDBRepo) DeleteBookGenre(book_id, genre_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `DELETE FROM book_genres WHERE (book_id=$1 AND genre_id=$2)`
	_, err := m.DB.ExecContext(ctx, stmt, book_id, genre_id)
	return err
}

// GetBookGenreByID returns the book from database using id
func (m *postgresDBRepo) GetBookGenreByID(book_id, genre_id int) (*models.BookGenre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT * FROM book_genres
		WHERE (book_id=$1 AND genre_id=$2)
	`
	row := m.DB.QueryRowContext(ctx, query, book_id, genre_id)
	bookGenre := &models.BookGenre{}
	if err := row.Scan(
		&bookGenre.BookID,
		&bookGenre.GenreID,
	); err != nil {
		return nil, err
	}
	return bookGenre, nil
}

// BookGenreExists return true if book genre relation exists else return false
func (m *postgresDBRepo) BookGenreExists(book_id, genre_id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT COUNT(*) FROM book_genres
		WHERE (book_id=$1 AND genre_id=$2)
	`
	var count int
	if err := m.DB.QueryRowContext(ctx, query, book_id, genre_id).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return count > 0, nil
}

// UpdateBookGenre updates the book genre
// Takes update value BookGenre model and previous book_id , genre_id
func (m *postgresDBRepo) UpdateBookGenre(u *models.BookGenre, book_id, genre_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE book_genres
		SET book_id = $3, genre_id = $4
		WHERE (book_id = $1 AND genre_id = $2)
	`
	_, err := m.DB.ExecContext(ctx, stmt, book_id, genre_id, u.BookID, u.GenreID)
	if err != nil {
		return err
	}
	return nil
}

// InsertBookGenre add new book genre to db
// Takes BookGenre model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertBookGenre(u *models.BookGenre) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO book_genres (book_id, genre_id)
		VALUES ($1, $2);
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.BookID,
		u.GenreID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetGenresFromBookID(book_id int) ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT 
		COALESCE(g.id, 0), COALESCE(g.title, '')
		FROM
			book_genres AS bg
		LEFT JOIN
			genres AS g ON g.id = bg.genre_id
		WHERE 
			book_id = $1
	`
	genres := []*models.Genre{}
	rows, err := m.DB.QueryContext(ctx, query, book_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		genre := &models.Genre{}
		if err := rows.Scan(
			&genre.ID,
			&genre.Title,
		); err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (m *postgresDBRepo) GetAllBooksByGenre(limit, page int, searchKey, sort, genre string) (*models.BookApiFilter, error) {
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
			book_genres AS bg
		LEFT JOIN 
			books AS b ON b.id = bg.book_id
		LEFT JOIN
			genres AS g ON g.id = bg.genre_id
		WHERE
			g.title = $1
	`
	countQuery := `
		SELECT 
			COUNT(*)
		FROM 
			book_genres AS bg
		LEFT JOIN 
			books AS b ON b.id = bg.book_id
		LEFT JOIN
			genres AS g ON g.id = bg.genre_id
		WHERE
			g.title = $1
	`
	if searchKey != "" {
		query = fmt.Sprintf("%s AND (b.title LIKE '%%%s%%' OR CAST(b.isbn AS TEXT) LIKE '%%%s%%')", query, searchKey, searchKey)
		countQuery = fmt.Sprintf("%s AND (b.title LIKE '%%%s%%' OR CAST(b.isbn AS TEXT) LIKE '%%%s%%')", countQuery, searchKey, searchKey)
	}
	if sort != "" {
		query = fmt.Sprintf("%s ORDER BY b.title %s", query, sort)
	}

	var count int
	if err := m.DB.QueryRowContext(ctx, countQuery, genre).Scan(&count); err != nil {
		return nil, err
	}

	query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
	rows, err := m.DB.QueryContext(ctx, query, genre)
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
