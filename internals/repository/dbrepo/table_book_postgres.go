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

func (m *postgresDBRepo) AllBookData(limit, page int) ([]*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	skip := (page - 1) * limit
	query := `SELECT * FROM books ORDER BY title ASC LIMIT $1 OFFSET $2`
	rows, err := m.DB.QueryContext(ctx, query, limit, skip)
	if err != nil {
		return nil, err
	}
	books := []*models.Book{}
	for rows.Next() {
		book := new(models.Book)
		if err := rows.Scan(
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
		books = append(books, book)
	}
	return books, nil
}

func (m *postgresDBRepo) AllBookDataRandom() ([]*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM books ORDER BY RANDOM()`
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
		books = append(books, book)
	}
	return books, nil
}

// AllBookPage returns slice of books of length limit
func (m *postgresDBRepo) AllBookRandomPage(limit, page int) ([]*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if limit == 0 || limit < 0 {
		limit = 10
	}
	if page == 0 || page < 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := "SELECT * FROM books ORDER BY RANDOM() LIMIT $1 OFFSET $2"
	rows, err := m.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	books := []*models.Book{}
	for rows.Next() {
		book := &models.Book{}
		if err := rows.Scan(
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

// GetBookByISBN returns the book from database using id
func (m *postgresDBRepo) GetBookByISBN(isbn int64) (*models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT * FROM books
		WHERE isbn=$1
	`
	row := m.DB.QueryRowContext(ctx, query, isbn)
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
	query := `SELECT id, title FROM books WHERE id=$1`
	book := &models.Book{}
	row := m.DB.QueryRowContext(ctx, query, id)
	if err := row.Scan(&book.ID, &book.Title); err != nil {
		return nil, err
	}
	return book, nil
}

func (m *postgresDBRepo) TotalBooks() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT COUNT(*) FROM books"
	var count int
	if err := m.DB.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (m *postgresDBRepo) AllBooksFilter(limit, page int, searchKey, sort string) (*models.BookApiFilter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	sql := "SELECT id, title, description, cover, isbn, published_date FROM books"
	if searchKey != "" {
		sql = fmt.Sprintf("%s WHERE title LIKE '%%%s%%' OR description LIKE '%%%s%%'", sql, searchKey, searchKey)
	}

	if sort != "" {
		sql = fmt.Sprintf("%s ORDER BY title %s", sql, sort)
	}
	sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, limit, offset)

	res, err := m.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	books := []*models.Book{}
	for res.Next() {
		book := &models.Book{}
		if err := res.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.Cover,
			&book.Isbn,
			&book.PublishedDate,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	count, _ := m.TotalBooks()
	lastPage := m.CalculateLastPage(limit, count)
	return &models.BookApiFilter{
		Total:    count,
		Page:     page,
		LastPage: lastPage,
		Books:    books,
	}, nil
}

func (m *postgresDBRepo) CalculateLastPage(limit, total int) int {
	if limit <= 0 {
		return 1
	}
	lastPage := (total + limit - 1) / limit
	if lastPage <= 0 {
		return 1
	}
	return lastPage
}
