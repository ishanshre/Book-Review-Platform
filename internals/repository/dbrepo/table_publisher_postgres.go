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

// PublisherExists return false if does not else true
func (m *postgresDBRepo) PublisherExistsID(id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT COUNT(*) FROM publishers
		WHERE id=$1
	`
	var count int
	row := m.DB.QueryRowContext(ctx, query, id)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute the query: %w", err)
	}
	return count > 0, nil
}

func (m *postgresDBRepo) GetPublisherWithBookByID(publisher_id int) (*models.PublisherWithBooksData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT p.id, p.name, p.description, p.pic, p.address, p.phone, p.email, p.website, p.established_date, p.latitude, p.longitude, COALESCE(b.title, ''), COALESCE(b.isbn, 0), COALESCE(b.cover, '')
		FROM publishers AS p
		LEFT JOIN books AS b ON b.publisher_id = p.id
		WHERE p.id = $1;
	`
	publisher := &models.Publisher{}
	books := []*models.Book{}
	rows, err := m.DB.QueryContext(ctx, query, publisher_id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		book := &models.Book{}
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
			&book.Title,
			&book.Isbn,
			&book.Cover,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	publisherWithBooks := &models.PublisherWithBooksData{
		Publisher: publisher,
		Books:     books,
	}
	return publisherWithBooks, nil
}

func (m *postgresDBRepo) AllPublishersFilter(limit, page int, searchKey, sort string) (*models.AdminPublisherListApi, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := `SELECT id, name, established_date FROM publishers`
	countQuery := `SELECT COUNT(*) FROM publishers`
	if searchKey != "" {
		query = fmt.Sprintf("%s WHERE name LIKE '%%%s%%' OR address LIKE '%%%s%%' OR email LIKE '%%%s%%' OR website LIKE '%%%s%%'", query, searchKey, searchKey, searchKey, searchKey)
		countQuery = fmt.Sprintf("%s WHERE name LIKE '%%%s%%' OR address LIKE '%%%s%%' OR email LIKE '%%%s%%' OR website LIKE '%%%s%%'", countQuery, searchKey, searchKey, searchKey, searchKey)
	}
	if sort != "" {
		query = fmt.Sprintf("%s ORDER BY name %s", query, sort)
	}
	var count int
	if err := m.DB.QueryRowContext(ctx, countQuery).Scan(&count); err != nil {
		return nil, err
	}
	query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
	publishers := []*models.AdminPublisherList{}
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		publisher := &models.AdminPublisherList{}
		if err := rows.Scan(
			&publisher.ID,
			&publisher.Name,
			&publisher.EstablishedDate,
		); err != nil {
			return nil, err
		}
		publishers = append(publishers, publisher)
	}
	lastPage := m.CalculateLastPage(limit, count)
	return &models.AdminPublisherListApi{
		Total:      count,
		Page:       page,
		LastPage:   lastPage,
		Publishers: publishers,
	}, nil
}

// Get the total publishers from database
func (m *postgresDBRepo) TotalPulbishersCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT COUNT(*) FROM publishers;`
	var count int
	if err := m.DB.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
