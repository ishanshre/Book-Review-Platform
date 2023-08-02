package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// Author interface implementation
func (m *postgresDBRepo) AllAuthor() ([]*models.Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM authors`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	authors := []*models.Author{}
	for rows.Next() {
		author := new(models.Author)
		if err := rows.Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
			&author.Bio,
			&author.DateOfBirth,
			&author.Email,
			&author.CountryOfOrigin,
			&author.Avatar,
		); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

// InsertAuthor add new author to db
func (m *postgresDBRepo) InsertAuthor(u *models.Author) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO authors (first_name, last_name, bio, date_of_birth, email, country_of_origin, avatar)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.FirstName,
		u.LastName,
		u.Bio,
		u.DateOfBirth,
		u.Email,
		u.CountryOfOrigin,
		u.Avatar,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAuthor updates the existing author in db
func (m *postgresDBRepo) UpdateAuthor(u *models.Author) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE authors
		SET first_name=$2, last_name=$3, bio=$4, date_of_birth=$5, email=$6, country_of_origin=$7, avatar=$8
		WHERE id=$1; 
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.FirstName,
		u.LastName,
		u.Bio,
		u.DateOfBirth,
		u.Email,
		u.CountryOfOrigin,
		u.Avatar,
	)
	return err
}

// DeleteAuthor deletes the author from the db
func (m *postgresDBRepo) DeleteAuthor(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `DELETE FROM authors WHERE id=$1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}

// GetAuthorByID fetches the author detail from the database
func (m *postgresDBRepo) GetAuthorByID(id int) (*models.Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT * FROM authors
		WHERE id=$1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	author := &models.Author{}
	if err := row.Scan(
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.Bio,
		&author.DateOfBirth,
		&author.Email,
		&author.CountryOfOrigin,
		&author.Avatar,
	); err != nil {
		return nil, err
	}
	return author, nil
}

// GetAuthorFullNameByID return full name of the author
func (m *postgresDBRepo) GetAuthorFullNameByID(id int) (*models.Author, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT first_name, last_name FROM authors WHERE id=$1`
	author := &models.Author{}
	row := m.DB.QueryRowContext(ctx, query, id)
	if err := row.Scan(&author.FirstName, &author.LastName); err != nil {
		return nil, err
	}
	return author, nil
}

func (m *postgresDBRepo) TotalAuthors() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT COUNT(*) FROM authors"
	var count int
	if err := m.DB.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (m *postgresDBRepo) AllAuthorsFilter(limit, page int, search, sort string) (*models.AuthorApiFilter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	sql := "SELECT id, first_name, last_name, avatar FROM authors"
	if search != "" {
		sql = fmt.Sprintf("%s where first_name LIKE '%%%s%%' OR last_name LIKE '%%%s%%'", sql, search, search)
	}
	if sort != "" {
		sql = fmt.Sprintf("%s ORDER BY first_name %s", sql, sort)
	}
	sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, limit, offset)
	res, err := m.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	authors := []*models.Author{}
	for res.Next() {
		author := &models.Author{}
		if err := res.Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
			&author.Avatar,
		); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	count, _ := m.TotalBooks()
	lastPage := m.CalculateLastPage(limit, count)
	return &models.AuthorApiFilter{
		Total:    count,
		Page:     page,
		LastPage: lastPage,
		Authors:  authors,
	}, nil
}
