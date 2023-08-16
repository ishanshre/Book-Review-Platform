package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// AllBookAuthor fetches all Book author relation from database
func (m *postgresDBRepo) AllBookAuthor() ([]*models.BookAuthor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM book_authors`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	bookAuthors := []*models.BookAuthor{}
	for rows.Next() {
		bookAuthor := new(models.BookAuthor)
		if err := rows.Scan(
			&bookAuthor.BookID,
			&bookAuthor.AuthorID,
		); err != nil {
			return nil, err
		}
		bookAuthors = append(bookAuthors, bookAuthor)
	}
	return bookAuthors, nil
}

// DeleteBookAuthor deletes the Book author relation from the db
func (m *postgresDBRepo) DeleteBookAuthor(book_id, author_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `DELETE FROM book_authors WHERE (book_id=$1 AND author_id=$2)`
	_, err := m.DB.ExecContext(ctx, stmt, book_id, author_id)
	return err
}

// GetBookAuthorByID returns the book-author relation from database using id
func (m *postgresDBRepo) GetBookAuthorByID(book_id, author_id int) (*models.BookAuthor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT * FROM book_authors
		WHERE (book_id=$1 AND author_id=$2)
	`
	row := m.DB.QueryRowContext(ctx, query, book_id, author_id)
	bookAuthor := &models.BookAuthor{}
	if err := row.Scan(
		&bookAuthor.BookID,
		&bookAuthor.AuthorID,
	); err != nil {
		return nil, err
	}
	return bookAuthor, nil
}

func (m *postgresDBRepo) GetBookAuthorByBookID(book_id int) ([]*models.BookAuthor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT * FROM book_authors
		WHERE book_id=$1
	`
	rows, err := m.DB.QueryContext(ctx, query, book_id)
	if err != nil {
		return nil, err
	}
	bookAuthors := []*models.BookAuthor{}
	for rows.Next() {
		bookAuthor := &models.BookAuthor{}
		if err := rows.Scan(
			&bookAuthor.BookID,
			&bookAuthor.AuthorID,
		); err != nil {
			return nil, err
		}
		bookAuthors = append(bookAuthors, bookAuthor)
	}
	return bookAuthors, nil
}

// BookAuthorExists return true if book author relation exists else return false
func (m *postgresDBRepo) BookAuthorExists(book_id, author_id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT COUNT(*) FROM book_authors
		WHERE (book_id=$1 AND author_id=$2)
	`
	var count int
	if err := m.DB.QueryRowContext(ctx, query, book_id, author_id).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return count > 0, nil
}

// UpdateBookAuthor updates the book author relation
func (m *postgresDBRepo) UpdateBookAuthor(u *models.BookAuthor, book_id, author_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE book_authors
		SET book_id = $3, author_id = $4
		WHERE (book_id = $1 AND author_id = $2)
	`
	_, err := m.DB.ExecContext(ctx, stmt, book_id, author_id, u.BookID, u.AuthorID)
	if err != nil {
		return err
	}
	return nil
}

// InsertBookAuthor add new book-author relation to db
func (m *postgresDBRepo) InsertBookAuthor(u *models.BookAuthor) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO book_authors (book_id, author_id)
		VALUES ($1, $2);
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.BookID,
		u.AuthorID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDBRepo) BookAuthorListFilter(limit, page int, searchKey, sort string) (*models.BookAuthorListApi, error) {
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
		SELECT b.id, b.title, a.id, a.first_name, a.last_name
		FROM book_authors AS ba
		JOIN 
			books AS b ON b.id = ba.book_id
		JOIN 
			authors AS a ON a.id = ba.author_id
	`
	countSql := `
	SELECT 
		COUNT(*)
		FROM book_authors AS ba
		JOIN 
			books AS b ON b.id = ba.book_id
		JOIN 
			authors AS a ON a.id = ba.author_id
	`
	if searchKey != "" {
		sql = fmt.Sprintf("%s WHERE b.title LIKE '%%%s%%' OR a.first_name LIKE '%%%s%%' OR a.last_name LIKE '%%%s%%'", sql, searchKey, searchKey, searchKey)
		countSql = fmt.Sprintf("%s WHERE b.title LIKE '%%%s%%' OR a.first_name LIKE '%%%s%%' OR a.last_name LIKE '%%%s%%'", countSql, searchKey, searchKey, searchKey)

	}
	if sort != "" {
		sql = fmt.Sprintf("%s ORDER BY b.title %s", sql, sort)
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
	book_authors := []*models.BookAuthorList{}
	for rows.Next() {
		book_author := &models.BookAuthorList{}
		if err := rows.Scan(
			&book_author.BookID,
			&book_author.BookTitle,
			&book_author.AuthorID,
			&book_author.AuthorFirstName,
			&book_author.AuthorLastName,
		); err != nil {
			return nil, err
		}
		book_authors = append(book_authors, book_author)
	}
	lastPage := m.CalculateLastPage(limit, count)
	return &models.BookAuthorListApi{
		Total:       count,
		Page:        page,
		LastPage:    lastPage,
		BookAuthors: book_authors,
	}, nil
}
