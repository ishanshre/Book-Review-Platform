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
