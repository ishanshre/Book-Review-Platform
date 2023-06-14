package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"golang.org/x/crypto/bcrypt"
)

// AllUsers returns list of all the users with all access level
func (m *postgresDBRepo) AllUsers(limit, offset int) ([]*models.User, error) {
	// creating database transcation atomic with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query stores the sql query statement
	query := `
		SELECT id, first_name, last_name, username, access_level, is_validated, created_at
		FROM users
		LIMIT $1 OFFSET $2
	`
	// QueryContext is used to execute query with database with context included
	rows, err := m.DB.QueryContext(
		ctx,
		query,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("could not fetch all users: %s", err)
	}
	users := []*models.User{}
	for rows.Next() {
		user := new(models.User)
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.AccessLevel,
			&user.IsValidated,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// AllReader fetch the list of all users whose access level is 2
func (m *postgresDBRepo) AllReaders(limit, offset int) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query stores sql query statment that retrives list of all users with access level 2
	query := `
		SELECT id, first_name, last_name, username, email
		FROM users
		WHERE access_level=2;
		LIMIT=$1 OFFSET=$2
	`
	rows, err := m.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not fetch record from the database: %s", err)
	}
	users := []*models.User{}
	for rows.Next() {
		user := new(models.User)
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Username,
			&user.Email,
		); err != nil {
			return nil, fmt.Errorf("could not fetch record from the database: %s", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// GetUserByID fetch data by id and only for authenticated user and owner
func (m *postgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	u := &models.User{}
	query := `
		SELECT * FROM users WHERE id=$1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	if err := row.Scan(
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.Gender,
		&u.Address,
		&u.Phone,
		&u.ProfilePic,
		&u.CitizenshipNumber,
		&u.CitizenshipFront,
		&u.CitizenshipBack,
		&u.AccessLevel,
		&u.IsValidated,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLogin,
	); err != nil {
		return nil, err
	}
	return u, nil
}

// GetGlobalUserByID return user by id
func (m *postgresDBRepo) GetGlobalUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT id, first_name, last_name, gender, address, profile_pic, created_at
		FROM users
		WHERE access_level = 2
	`
	row := m.DB.QueryRowContext(ctx, query)
	var u *models.User
	if err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Gender,
		&u.Address,
		&u.ProfilePic,
		&u.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("could not fetch id %d from database: %s", id, err)
	}
	return u, nil
}

// DeleteUser deletes the user from database
func (m *postgresDBRepo) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		DELETE FROM users
		WHERE id = $1
	`
	res, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf("failed to delete user from database: %s", err)
	}
	rows_affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected == 0 {
		return fmt.Errorf("error in delete user from database: %s", err)
	}
	return nil
}

// Update user updates user information by id.
// Update Fields :- First Name, Last Name, Email, Gender, Address, Phone and ProfilePic
func (m *postgresDBRepo) UpdateUser(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE users
		SET first_name = $2, last_name = $3, gender = $4, address = $5, phone = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.FirstName,
		u.LastName,
		u.Gender,
		u.Address,
		u.Phone,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("cannot update the user with id %d : %s", u.ID, err)
	}
	return nil
}

// UpdateProfilePic updates user profile pic
func (m *postgresDBRepo) UpdateProfilePic(path string, id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `UPDATE users SET profile_pic=$2 WHERE id=$1`
	_, err := m.DB.ExecContext(ctx, stmt, id, path)
	if err != nil {
		return err
	}
	return nil
}

// InsertUser insert new user into database.
// This method is used for new user sign up
func (m *postgresDBRepo) InsertUser(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO users (first_name, last_name, email, username, password,address, gender, phone, profile_pic, citizenship_number, citizenship_front, citizenship_back, created_at, updated_at, last_login)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	res, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Username,
		u.Password,
		"",
		u.Gender,
		"",
		"",
		u.CitizenshipNumber,
		u.CitizenshipFront,
		u.CitizenshipBack,
		time.Now(),
		time.Time{},
		time.Time{},
	)
	if err != nil {
		return fmt.Errorf("could not create new user: %s", err)
	}
	rows_affected, _ := res.RowsAffected()
	if rows_affected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

// AdminInsertsUser insert user to db by admin
func (m *postgresDBRepo) AdminInsertUser(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO users (first_name, last_name, email, username, password,address, gender, phone, profile_pic, citizenship_number, citizenship_front, citizenship_back, created_at, updated_at, last_login)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	res, err := m.DB.ExecContext(
		ctx,
		stmt,
		"",
		"",
		u.Email,
		u.Username,
		u.Password,
		"",
		"Male",
		"",
		"",
		"",
		"",
		"",
		time.Now(),
		time.Now(),
		time.Time{},
	)
	if err != nil {
		return fmt.Errorf("could not create new user: %s", err)
	}
	rows_affected, _ := res.RowsAffected()
	if rows_affected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

// UpdateLastLogin updates the last login date of the user
func (m *postgresDBRepo) UpdateLastLogin(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE users
		SET last_login = $2
		WHERE id = $1
	`
	_, err := m.DB.ExecContext(ctx, stmt, id, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Authenticate retrives password and id using username.
// It compares the hash of retrived and password provided.
// Returns id, hashed password and error.
func (m *postgresDBRepo) Authenticate(username, testPassword string) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var access_level int
	var hashedPassword string
	query := `SELECT id, password, access_level FROM users WHERE username=$1`
	row := m.DB.QueryRowContext(ctx, query, username)
	if err := row.Scan(&id, &hashedPassword, &access_level); err != nil {
		return id, 2, err
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, 2, fmt.Errorf("incorrect password: %s", bcrypt.ErrMismatchedHashAndPassword)
	} else if err != nil {
		return 0, 2, err
	}
	return id, access_level, nil
}

// Get information for personal profile page
func (m *postgresDBRepo) GetProfilePersonal(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT first_name, last_name, email, username, gender, address, phone, profile_pic, citizenship_number, citizenship_front, citizenship_back, created_at, updated_at, last_login
		FROM users
		WHERE id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	u := &models.User{}
	if err := row.Scan(
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Username,
		&u.Gender,
		&u.Address,
		&u.Phone,
		&u.ProfilePic,
		&u.CitizenshipNumber,
		&u.CitizenshipFront,
		&u.CitizenshipBack,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLogin,
	); err != nil {
		return nil, err
	}
	return u, nil
}

// UsernameExists checks if username already exists in database.
// It returns true if username exists else return false
func (m *postgresDBRepo) UsernameExists(username string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT id FROM users
		WHERE username=$1
	`
	var id int
	row := m.DB.QueryRowContext(ctx, query, username)
	if err := row.Scan(&id); err != nil {
		return false, err

	}
	return true, nil
}

// EmailExists return true if email exists else return false
func (m *postgresDBRepo) EmailExists(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT id FROM users
		WHERE email=$1
	`
	var id int
	row := m.DB.QueryRowContext(ctx, query, email)
	if err := row.Scan(&id); err != nil {
		return false, err
	}
	return true, nil
}

// ChangePassword chnage the password using email
func (m *postgresDBRepo) ChangePassword(password, email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE users
		SET password = $2, updated_at = $3
		WHERE email = $1
	`
	res, err := m.DB.ExecContext(ctx, stmt, email, password, time.Now())
	if err != nil {
		return err
	}
	rows_affected, _ := res.RowsAffected()
	if rows_affected == 0 {
		return fmt.Errorf("error in changing the password")
	}
	return nil
}

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
		SELECT id FROM genres
		WHERE title=$1
	`
	var id int
	row := m.DB.QueryRowContext(ctx, query, title)
	if err := row.Scan(&id); err != nil {
		return false, err
	}
	return true, nil
}

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
		SELECT id FROM publishers
		WHERE name=$1
	`
	var id int
	row := m.DB.QueryRowContext(ctx, query, name)
	if err := row.Scan(&id); err != nil {
		return false, err
	}
	return true, nil
}
