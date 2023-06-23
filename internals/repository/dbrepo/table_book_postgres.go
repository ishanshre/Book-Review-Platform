package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// Book interface implementation

// AllBook fetches all Books from database
func (m *postgresDBRepo) AllBook() ([]*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT id, title, is_active, added_at FROM books`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	books := []*models.Book{}
	for rows.Next() {
		book := new(models.Book)
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.IsActive,
			&book.AddedAt,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

// DeleteBook deletes the Book from the db
func (m *postgresDBRepo) DeleteBook(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `DELETE FROM books WHERE id=$1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}

// GetBookByID returns the book from database using id
func (m *postgresDBRepo) GetBookByID(id int) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT * FROM books
		WHERE id=$1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	book := &models.Book{}
	if err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Description,
		&book.Cover,
		&book.Isbn,
		&book.PublishedDate,
		&book.Paperback,
		&book.IsActive,
		&book.AddedAt,
		&book.UpdatedAt,
		&book.PublisherID,
	); err != nil {
		return nil, err
	}
	return book, nil
}

// InsertBook add new author to db
func (m *postgresDBRepo) InsertBook(u *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO books (title, description, cover, isbn, published_date, paperback, is_active, added_at, updated_at, publisher_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.Title,
		u.Description,
		u.Cover,
		u.Isbn,
		u.PublishedDate,
		u.Paperback,
		u.IsActive,
		u.AddedAt,
		u.UpdatedAt,
		u.PublisherID,
	)
	if err != nil {
		return err
	}
	return nil
}

// BookIsbnExists return false if does not else true
func (m *postgresDBRepo) BookIsbnExists(isbn int64) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT COUNT(*) FROM books
		WHERE isbn=$1
	`
	var count int
	row := m.DB.QueryRowContext(ctx, query, isbn)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute the query : %w", err)
	}
	return count > 0, nil
}

// UpdateBook updates the existing Book in db
func (m *postgresDBRepo) UpdateBook(u *models.Book) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE books
		SET title=$2, description=$3, isbn=$4, published_date=$5, paperback=$6, is_active=$7, publisher_id=$8, updated_at=$9
		WHERE id=$1; 
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.Title,
		u.Description,
		u.Isbn,
		u.PublishedDate,
		u.Paperback,
		u.IsActive,
		u.PublisherID,
		u.UpdatedAt,
	)
	return err
}

// GetBookTitleByID return title and id of the book
func (m *postgresDBRepo) GetBookTitleByID(id int) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT title FROM books WHERE id=$1`
	book := &models.Book{}
	row := m.DB.QueryRowContext(ctx, query, id)
	if err := row.Scan(&book.Title); err != nil {
		return nil, err
	}
	return book, nil
}
