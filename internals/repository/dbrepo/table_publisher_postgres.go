package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

// AllPublishers returns slice of all publishers
func (m *postgresDBRepo) AllPublishers() ([]*models.Publisher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM publishers`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	publishers := []*models.Publisher{}
	for rows.Next() {
		publisher := new(models.Publisher)
		if err := rows.Scan(
			&publisher.ID,
			&publisher.Name,
			&publisher.Description,
			&publisher.Pic,
			&publisher.Address,
			&publisher.Phone,
			&publisher.Email,
			&publisher.Website,
			&publisher.EstablishedDate,
			&publisher.Latitude,
			&publisher.Longitude,
		); err != nil {
			return nil, err
		}
		publishers = append(publishers, publisher)

	}
	return publishers, nil
}

// InsertPublisher add new genre to db
func (m *postgresDBRepo) InsertPublisher(u *models.Publisher) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO publishers (name, description, pic, address, phone, email, website, established_date, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.Name,
		u.Description,
		u.Pic,
		u.Address,
		u.Phone,
		u.Email,
		u.Website,
		u.EstablishedDate,
		u.Latitude,
		u.Longitude,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdatePublisher updates the existing Publisher in db
func (m *postgresDBRepo) UpdatePublisher(u *models.Publisher) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE publishers
		SET name = $2, description = $3, pic = $4, address = $5, phone = $6, email = $7, website = $8, established_date = $9, latitude = $10, longitude = $11
		WHERE id = $1
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.Name,
		u.Description,
		u.Pic,
		u.Address,
		u.Phone,
		u.Email,
		u.Website,
		u.EstablishedDate,
		u.Latitude,
		u.Longitude,
	)
	if err != nil {
		return err
	}
	return nil
}

// DeletePublisher deletes the existing Publisher from db
func (m *postgresDBRepo) DeletePublisher(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		DELETE FROM publishers
		WHERE id=$1;
	`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}
	return nil
}

// GetPublisherByID return Publisher using id
func (m *postgresDBRepo) GetPublisherByID(id int) (*models.Publisher, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM Publishers WHERE id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)
	u := &models.Publisher{}
	if err := row.Scan(
		&u.ID,
		&u.Name,
		&u.Description,
		&u.Pic,
		&u.Address,
		&u.Phone,
		&u.Email,
		&u.Website,
		&u.EstablishedDate,
		&u.Latitude,
		&u.Longitude,
	); err != nil {
		return nil, err
	}
	return u, nil
}

// PublisherExists return false if does not else true
func (m *postgresDBRepo) PublisherExists(name string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT COUNT(*) FROM publishers
		WHERE name=$1
	`
	var count int
	row := m.DB.QueryRowContext(ctx, query, name)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute the query: %w", err)
	}
	return count > 0, nil
}
