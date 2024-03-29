package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

func (m *postgresDBRepo) InsertRequestedBook(i *models.RequestedBook) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO request_books (book_title, author, requested_by, requested_date)
		VALUES ($1, $2, $3, $4)
	`
	_, err := m.DB.ExecContext(
		ctx,
		query,
		i.BookTitle,
		i.Author,
		i.RequestedBy,
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
		SELECT id, book_title, requested_by, requested_date, is_added
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
			&request_book.RequestedBy,
			&request_book.RequestedDate,
			&request_book.IsAdded,
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
		&request_book.RequestedBy,
		&request_book.RequestedDate,
		&request_book.IsAdded,
	); err != nil {
		return nil, err
	}
	return request_book, nil
}

func (m *postgresDBRepo) RequestedBooksListFilter(limit, page int, searchKey, sort string) (*models.RequestedBookFilterApi, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	sql := `
		SELECT 
			rb.id, rb.book_title, rb.author, u.id as user_id, u.username, u.email, rb.requested_date, rb.is_added
		FROM 
			request_books as rb
		LEFT JOIN 
			users AS u on 
				rb.requested_by=u.id
	`
	countSql := `
	SELECT 
		COUNT(*)
		FROM request_books
	`
	if searchKey != "" {
		sql = fmt.Sprintf("%s WHERE book_title LIKE '%%%s%%' OR author LIKE '%%%s%%'", sql, searchKey, searchKey)
		countSql = fmt.Sprintf("%s WHERE book_title LIKE '%%%s%%' OR author LIKE '%%%s%%'", countSql, searchKey, searchKey)
	}
	if sort != "" {
		sql = fmt.Sprintf("%s ORDER BY book_title %s", sql, sort)
		// countSql = fmt.Sprintf("%s ORDER BY u.username %s", countSql, sort)
	}
	var count int
	if err := m.DB.QueryRowContext(ctx, countSql).Scan(&count); err != nil {
		return nil, err
	}
	sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, limit, offset)
	rows, err := m.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	requestedBooks := []*models.RequestedBookUser{}
	for rows.Next() {
		requestedBook := &models.RequestedBookUser{}
		user := &models.User{}
		if err := rows.Scan(
			&requestedBook.ID,
			&requestedBook.BookTitle,
			&requestedBook.Author,
			&user.ID,
			&user.Username,
			&user.Email,
			&requestedBook.RequestedDate,
			&requestedBook.IsAdded,
		); err != nil {
			return nil, err
		}
		requestedBook.RequestedBy = user
		requestedBooks = append(requestedBooks, requestedBook)
	}
	lastPage := m.CalculateLastPage(limit, count)
	return &models.RequestedBookFilterApi{
		Total:          count,
		Page:           page,
		LastPage:       lastPage,
		RequestedBooks: requestedBooks,
	}, nil
}

func (m *postgresDBRepo) UpdateBookRequestStatus(request_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE request_books
		SET is_added = true
		WHERE id = $1
	`
	_, err := m.DB.ExecContext(ctx, stmt, request_id)
	if err != nil {
		return err
	}
	return nil
}
