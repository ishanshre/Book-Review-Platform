package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// Genre interface implementations

// AllGenre returns all the genre in db
func (m *postgresDBRepo) AllGenre() ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM genres`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	genres := []*models.Genre{}
	for rows.Next() {
		genre := new(models.Genre)
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

// InsertGenre add new genre to db
func (m *postgresDBRepo) InsertGenre(u *models.Genre) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO genres (title)
		VALUES ($1);
	`
	_, err := m.DB.ExecContext(ctx, stmt, u.Title)
	if err != nil {
		return err
	}
	return nil
}

// UpdateGenre updates the existing genre in db
func (m *postgresDBRepo) UpdateGenre(u *models.Genre) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE genres
		SET title = $2
		where id= $1
	`
	_, err := m.DB.ExecContext(ctx, stmt, u.ID, u.Title)
	if err != nil {
		return err
	}
	return nil
}

// DeleteGerre deletes the existing genre from db
func (m *postgresDBRepo) DeleteGenre(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		DELETE FROM genres
		WHERE id=$1;
	`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}

// GetGenreByID return genre using id
func (m *postgresDBRepo) GetGenreByID(id int) (*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM genres WHERE id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)
	u := &models.Genre{}
	if err := row.Scan(&u.ID, &u.Title); err != nil {
		return nil, err
	}
	return u, nil
}

// GenreExists return false if does not else true
func (m *postgresDBRepo) GenreExists(title string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT COUNT(*) FROM genres
		WHERE title=$1
	`
	var count int
	row := m.DB.QueryRowContext(ctx, query, title)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return count > 0, nil
}
