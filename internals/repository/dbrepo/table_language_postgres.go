package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// Language interface implementation

// AllLanguage fetches all languages from database
func (m *postgresDBRepo) AllLanguage() ([]*models.Language, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM languages`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	languages := []*models.Language{}
	for rows.Next() {
		language := new(models.Language)
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

// InsertLanguage add new author to db
func (m *postgresDBRepo) InsertLanguage(u *models.Language) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO languages (language)
		VALUES ($1);
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.Language,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateLanguage updates the existing Language in db
func (m *postgresDBRepo) UpdateLanguage(u *models.Language) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE languages
		SET language=$2
		WHERE id=$1; 
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.Language,
	)
	return err
}

// DeleteLanguage deletes the Language from the db
func (m *postgresDBRepo) DeleteLanguage(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `DELETE FROM languages WHERE id=$1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	return err
}

// GetLanguageByID fetches the Language detail from the database
func (m *postgresDBRepo) GetLanguageByID(id int) (*models.Language, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT * FROM languages
		WHERE id=$1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	language := &models.Language{}
	if err := row.Scan(
		&language.ID,
		&language.Language,
	); err != nil {
		return nil, err
	}
	return language, nil
}

// LanguageExists return false if does not else true
func (m *postgresDBRepo) LanguageExists(language string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT COUNT(*) FROM languages
		WHERE language=$1
	`
	var count int
	row := m.DB.QueryRowContext(ctx, query, language)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute the query :%w", err)
	}
	return count > 0, nil
}
