package dbrepo

import (
	"context"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

func (m *postgresDBRepo) InsertRequestedBook(i *models.RequestedBook) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO request_books (book_title, author, requested_email, requested_date)
		VALUES ($1, $2, $3, $4)
	`
	_, err := m.DB.ExecContext(
		ctx,
		query,
		i.BookTitle,
		i.Author,
		i.RequestedEmail,
		i.RequestedDate,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) AllRequestBooks() ([]*models.RequestedBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	query := `
		SELECT id, book_title, requested_date
		FROM request_books
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	request_books := []*models.RequestedBook{}
	for rows.Next() {
		request_book := &models.RequestedBook{}
		if err := rows.Scan(
			&request_book.ID,
			&request_book.BookTitle,
			&request_book.RequestedDate,
		); err != nil {
			return nil, err
		}
		request_books = append(request_books, request_book)
	}
	return request_books, nil
}

func (m *postgresDBRepo) DeleteRequestBooks(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	query := `DELETE FROM request_books WHERE id = $1`
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) GetRequestBookById(id int) (*models.RequestedBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	query := `SELECT * FROM request_books WHERE id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	request_book := &models.RequestedBook{}
	if err := row.Scan(
		&request_book.ID,
		&request_book.BookTitle,
		&request_book.Author,
		&request_book.RequestedEmail,
		&request_book.RequestedDate,
	); err != nil {
		return nil, err
	}
	return request_book, nil
}
