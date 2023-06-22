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
		u.CitizenshipNumber,
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
		SELECT COUNT(*) FROM users
		WHERE username=$1
	`
	var count int
	if err := m.DB.QueryRowContext(ctx, query, username).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return count > 0, nil
}

// EmailExists return true if email exists else return false
func (m *postgresDBRepo) EmailExists(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT COUNT(*) FROM users
		WHERE email=$1
	`
	var count int
	row := m.DB.QueryRowContext(ctx, query, email)
	if err := row.Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return count > 0, nil
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
	query := `SELECT title FROM books WHERE id=$1`
	book := &models.Book{}
	row := m.DB.QueryRowContext(ctx, query, id)
	if err := row.Scan(&book.Title); err != nil {
		return nil, err
	}
	return book, nil
}

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

// Book Genre db method implementation

// AllBookGenre fetches all record of Book Genre table from database
func (m *postgresDBRepo) AllBookGenre() ([]*models.BookGenre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT * FROM book_genres`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	bookGenres := []*models.BookGenre{}
	for rows.Next() {
		bookGenre := new(models.BookGenre)
		if err := rows.Scan(
			&bookGenre.BookID,
			&bookGenre.GenreID,
		); err != nil {
			return nil, err
		}
		bookGenres = append(bookGenres, bookGenre)
	}
	return bookGenres, nil
}

// DeleteBookGenre deletes the record of Book genre table from the db
func (m *postgresDBRepo) DeleteBookGenre(book_id, genre_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `DELETE FROM book_genres WHERE (book_id=$1 AND genre_id=$2)`
	_, err := m.DB.ExecContext(ctx, stmt, book_id, genre_id)
	return err
}

// GetBookGenreByID returns the book from database using id
func (m *postgresDBRepo) GetBookGenreByID(book_id, genre_id int) (*models.BookGenre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT * FROM book_genres
		WHERE (book_id=$1 AND genre_id=$2)
	`
	row := m.DB.QueryRowContext(ctx, query, book_id, genre_id)
	bookGenre := &models.BookGenre{}
	if err := row.Scan(
		&bookGenre.BookID,
		&bookGenre.GenreID,
	); err != nil {
		return nil, err
	}
	return bookGenre, nil
}

// BookGenreExists return true if book genre relation exists else return false
func (m *postgresDBRepo) BookGenreExists(book_id, genre_id int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT COUNT(*) FROM book_genres
		WHERE (book_id=$1 AND genre_id=$2)
	`
	var count int
	if err := m.DB.QueryRowContext(ctx, query, book_id, genre_id).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	return count > 0, nil
}

// UpdateBookGenre updates the book genre
// Takes update value BookGenre model and previous book_id , genre_id
func (m *postgresDBRepo) UpdateBookGenre(u *models.BookGenre, book_id, genre_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		UPDATE book_genres
		SET book_id = $3, genre_id = $4
		WHERE (book_id = $1 AND genre_id = $2)
	`
	_, err := m.DB.ExecContext(ctx, stmt, book_id, genre_id, u.BookID, u.GenreID)
	if err != nil {
		return err
	}
	return nil
}

// InsertBookGenre add new book genre to db
// Takes BookGenre model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertBookGenre(u *models.BookGenre) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO book_genres (book_id, genre_id)
		VALUES ($1, $2);
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.BookID,
		u.GenreID,
	)
	if err != nil {
		return err
	}
	return nil
}

// AllBookLanguage retrieves all book language relationships from the PostgreSQL database.
// It returns a slice of BookLanguage struct and error
func (m *postgresDBRepo) AllBookLanguage() ([]*models.BookLanguage, error) {

	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare the sql statement to select all book language relationship
	query := `SELECT * FROM book_languages`

	// Exectue the query and get the result row
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Initialize an empty slice of BookLanguage model
	bookLanguages := []*models.BookLanguage{}

	// Iterate through the rows
	for rows.Next() {

		// create a new instance of BookLanguage
		bookLanguage := new(models.BookLanguage)

		// Scan the value from current row and store in bookLanguage instance
		if err := rows.Scan(
			&bookLanguage.BookID,
			&bookLanguage.LanguageID,
		); err != nil {
			return nil, err
		}

		// append the bookLanguage to the slice
		bookLanguages = append(bookLanguages, bookLanguage)
	}

	// Return the retrieved book languages with no error
	return bookLanguages, nil
}

// DeleteBookLanguage deletes the record of Book Language table from the db.
// It takes book id and language id as parameter
func (m *postgresDBRepo) DeleteBookLanguage(book_id, language_id int) error {

	// Using context with timeout of 3 second
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the delete sql statment
	stmt := `DELETE FROM book_languages WHERE (book_id=$1 AND language_id=$2)`

	// executing the query
	_, err := m.DB.ExecContext(ctx, stmt, book_id, language_id)

	// returns nil if success else returns error
	return err
}

// GetBookLanguageByID returns the book from database using id.
// It takes book id and language id as parameters.
// Returns a BookLanguage struct instance.
func (m *postgresDBRepo) GetBookLanguageByID(book_id, language_id int) (*models.BookLanguage, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM book_languages
		WHERE (book_id=$1 AND language_id=$2)
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, book_id, language_id)

	// Initializing a BookLanguage struct instance
	bookLanguage := &models.BookLanguage{}

	// Scannin the row and storing the result in BookLanguage Intance.
	if err := row.Scan(
		&bookLanguage.BookID,
		&bookLanguage.LanguageID,
	); err != nil {
		return nil, err
	}

	// Return a BookLanguage Instance and nil
	return bookLanguage, nil
}

// BookLanguageExists return true if book Language relation exists else return false.
// It takes book id and language id as parameters
func (m *postgresDBRepo) BookLanguageExists(book_id, language_id int) (bool, error) {

	// Creating a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the sql query to check for existing relationship
	query := `
		SELECT COUNT(*) FROM book_languages
		WHERE (book_id=$1 AND language_id=$2)
	`

	// intializing a count variable that stores the no of records
	var count int

	// Executing the query row context and store the total record in count variable.
	// If any error occurs, false and error are returned
	if err := m.DB.QueryRowContext(ctx, query, book_id, language_id).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	// returning exists is true if count > 0 else retuirn false
	return count > 0, nil
}

// UpdateBookLanguage updates the book Language
// Takes update value BookLanguage model and previous book_id , language_id
func (m *postgresDBRepo) UpdateBookLanguage(u *models.BookLanguage, book_id, language_id int) error {

	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the query statement for update book language relationship
	stmt := `
		UPDATE book_languages
		SET book_id = $3, language_id = $4
		WHERE (book_id = $1 AND language_id = $2)
	`

	// Executing the sql query
	_, err := m.DB.ExecContext(ctx, stmt, book_id, language_id, u.BookID, u.LanguageID)
	if err != nil {
		return err
	}
	return nil
}

// InsertBookLanguage add new book Language to db
// Takes BookLanguage model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertBookLanguage(u *models.BookLanguage) error {

	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare a insert query statement
	stmt := `
		INSERT INTO book_languages (book_id, language_id)
		VALUES ($1, $2);
	`

	// Executing the query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.BookID,
		u.LanguageID,
	)
	if err != nil {
		return err
	}
	return nil
}

// AllReadList fetches all the records from readLists db table.
func (m *postgresDBRepo) AllReadList() ([]*models.ReadList, error) {
	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the sql statement
	query := `SELECT * FROM read_lists`

	// Execute the query using Query Context.
	// If any error occurs, nil and error is returned
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Create a slice of ReadList model
	readLists := []*models.ReadList{}

	// Loop through the record.
	for rows.Next() {

		// Store the record in ReadList instance
		readList := new(models.ReadList)
		if err := rows.Scan(&readList.UserID, &readList.BookID, &readList.CreatedAt); err != nil {
			return nil, err
		}

		// Append the ReadList instance to the slice of ReadList
		readLists = append(readLists, readList)
	}

	// Return readLists
	return readLists, nil
}

// ReadListExists return true if ReadList book and user relation exists else return false.
// It takes book id and language id as parameters
func (m *postgresDBRepo) ReadListExists(user_id, book_id int) (bool, error) {

	// Creating a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the sql query to check for existing relationship
	query := `
		SELECT COUNT(*) FROM read_lists
		WHERE (user_id=$1 AND book_id=$2)
	`

	// intializing a count variable that stores the no of records
	var count int

	// Executing the query row context and store the total record in count variable.
	// If any error occurs, false and error are returned
	if err := m.DB.QueryRowContext(ctx, query, user_id, book_id).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	// returning exists is true if count > 0 else retuirn false
	return count > 0, nil
}

// InsertReadList add new book user readlist relation to db
// Takes ReadList model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertReadList(u *models.ReadList) error {

	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare a insert query statement
	stmt := `
		INSERT INTO read_lists (user_id, book_id, created_at)
		VALUES ($1, $2, $3);
	`

	// Executing the query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.UserID,
		u.BookID,
		u.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetReadListByID returns the readlist detail from database using id.
// It takes book id and user id as parameters.
// Returns a ReadList struct instance.
func (m *postgresDBRepo) GetReadListByID(user_id, book_id int) (*models.ReadList, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM read_lists
		WHERE (user_id=$1 AND book_id=$2)
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, user_id, book_id)

	// Initializing a ReadList struct instance
	readList := &models.ReadList{}

	// Scannin the row and storing the result in ReadList Intance.
	if err := row.Scan(
		&readList.UserID,
		&readList.BookID,
		&readList.CreatedAt,
	); err != nil {
		return nil, err
	}

	// Return a ReadList Instance and nil
	return readList, nil
}

// DeleteReadList deletes the record of read_lists table from the db.
// It takes book id and user id as parameter
func (m *postgresDBRepo) DeleteReadList(user_id, book_id int) error {

	// Using context with timeout of 3 second
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the delete sql statment
	stmt := `DELETE FROM read_lists WHERE (user_id=$1 AND book_id=$2)`

	// executing the query
	_, err := m.DB.ExecContext(ctx, stmt, user_id, book_id)

	// returns nil if success else returns error
	return err
}

// UpdateReadList updates the read_lists
// Takes update value ReadList model and previous book_id , user
func (m *postgresDBRepo) UpdateReadList(u *models.ReadList, book_id, user_id int) error {

	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the query statement for update readList
	stmt := `
		UPDATE read_lists
		SET user_id = $3, book_id = $4
		WHERE (user_id = $1 AND book_id = $2)
	`

	// Executing the sql query
	_, err := m.DB.ExecContext(ctx, stmt, user_id, book_id, u.UserID, u.BookID)
	if err != nil {
		return err
	}
	return nil
}

// AllBuyList fetches all the records from BuyLists db table.
func (m *postgresDBRepo) AllBuyList() ([]*models.BuyList, error) {
	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the sql statement
	query := `SELECT * FROM buy_lists`

	// Execute the query using Query Context.
	// If any error occurs, nil and error is returned
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Create a slice of buyList model
	buyLists := []*models.BuyList{}

	// Loop through the record.
	for rows.Next() {

		// Store the record in ReadList instance
		buyList := new(models.BuyList)
		if err := rows.Scan(&buyList.UserID, &buyList.BookID, &buyList.CreatedAt); err != nil {
			return nil, err
		}

		// Append the BuyList instance to the slice of BuyList
		buyLists = append(buyLists, buyList)
	}

	// Return readLists
	return buyLists, nil
}

// BuyListExists return true if BuyList book and user relation exists else return false.
// It takes book id and language id as parameters
func (m *postgresDBRepo) BuyListExists(user_id, book_id int) (bool, error) {

	// Creating a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the sql query to check for existing relationship
	query := `
		SELECT COUNT(*) FROM buy_lists
		WHERE (user_id=$1 AND book_id=$2)
	`

	// intializing a count variable that stores the no of records
	var count int

	// Executing the query row context and store the total record in count variable.
	// If any error occurs, false and error are returned
	if err := m.DB.QueryRowContext(ctx, query, user_id, book_id).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	// returning exists is true if count > 0 else retuirn false
	return count > 0, nil
}

// InsertBuyList add new book user buy_lists relation to db
// Takes BuyList model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertBuyList(u *models.BuyList) error {

	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare a insert query statement
	stmt := `
		INSERT INTO buy_lists (user_id, book_id, created_at)
		VALUES ($1, $2, $3);
	`

	// Executing the query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.UserID,
		u.BookID,
		u.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetBuyListByID returns the Buylist detail from database using id.
// It takes book id and user id as parameters.
// Returns a BuyList struct instance.
func (m *postgresDBRepo) GetBuyListByID(user_id, book_id int) (*models.BuyList, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM buy_lists
		WHERE (user_id=$1 AND book_id=$2)
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, user_id, book_id)

	// Initializing a ReadList struct instance
	buyList := &models.BuyList{}

	// Scannin the row and storing the result in BuyList Intance.
	if err := row.Scan(
		&buyList.UserID,
		&buyList.BookID,
		&buyList.CreatedAt,
	); err != nil {
		return nil, err
	}

	// Return a BuyList Instance and nil
	return buyList, nil
}

// DeleteBuyList deletes the record of Buy_lists table from the db.
// It takes book id and user id as parameter
func (m *postgresDBRepo) DeleteBuyList(user_id, book_id int) error {

	// Using context with timeout of 3 second
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the delete sql statment
	stmt := `DELETE FROM buy_lists WHERE (user_id=$1 AND book_id=$2)`

	// executing the query
	_, err := m.DB.ExecContext(ctx, stmt, user_id, book_id)

	// returns nil if success else returns error
	return err
}

// UpdateBuyList updates the Buy_lists
// Takes update value BuyList model and previous book_id , user
func (m *postgresDBRepo) UpdateBuyList(u *models.BuyList, book_id, user_id int) error {

	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the query statement for update readList
	stmt := `
		UPDATE buy_lists
		SET user_id = $3, book_id = $4
		WHERE (user_id = $1 AND book_id = $2)
	`

	// Executing the sql query
	_, err := m.DB.ExecContext(ctx, stmt, user_id, book_id, u.UserID, u.BookID)
	if err != nil {
		return err
	}
	return nil
}

// AllFollowers fetches all the records from followers db table.
func (m *postgresDBRepo) AllFollowers() ([]*models.Follower, error) {
	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the sql statement
	query := `SELECT * FROM followers`

	// Execute the query using Query Context.
	// If any error occurs, nil and error is returned
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Create a slice of Follower model
	followers := []*models.Follower{}

	// Loop through the record.
	for rows.Next() {

		// Store the record in Follower instance
		follower := new(models.Follower)
		if err := rows.Scan(&follower.UserID, &follower.AuthorID, &follower.FollowedAt); err != nil {
			return nil, err
		}

		// Append the ReadList instance to the slice of ReadList
		followers = append(followers, follower)
	}

	// Return followers
	return followers, nil
}

// FollowerExists return true if Follower book and user relation exists else return false.
// It takes pointer to Follower model instance as parameters
func (m *postgresDBRepo) FollowerExists(u *models.Follower) (bool, error) {

	// Creating a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the sql query to check for existing relationship
	query := `
		SELECT COUNT(*) FROM followers
		WHERE (user_id=$1 AND author_id=$2)
	`

	// intializing a count variable that stores the no of records
	var count int

	// Executing the query row context and store the total record in count variable.
	// If any error occurs, false and error are returned
	if err := m.DB.QueryRowContext(ctx, query, u.UserID, u.AuthorID).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	// returning exists is true if count > 0 else retuirn false
	return count > 0, nil
}

// InsertFollower add new book user follower relation to db
// Takes Follower model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertFollower(u *models.Follower) error {

	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare a insert query statement
	stmt := `
		INSERT INTO followers (user_id, author_id, followed_at)
		VALUES ($1, $2, $3);
	`

	// Executing the query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.UserID,
		u.AuthorID,
		u.FollowedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetFollowerByID returns the Follower detail from database using id.
// It takes book id and user id as parameters.
// Returns a Follower struct instance.
func (m *postgresDBRepo) GetFollowerByID(user_id, author_id int) (*models.Follower, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM followers
		WHERE (user_id=$1 AND author_id=$2)
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, user_id, author_id)

	// Initializing a ReadList struct instance
	follower := &models.Follower{}

	// Scannin the row and storing the result in BuyList Intance.
	if err := row.Scan(
		&follower.UserID,
		&follower.AuthorID,
		&follower.FollowedAt,
	); err != nil {
		return nil, err
	}

	// Return a Follower Instance and nil
	return follower, nil
}

// DeleteFollower deletes the record of Follower table from the db.
// It takes book id and user id as parameter
func (m *postgresDBRepo) DeleteFollower(user_id, author_id int) error {

	// Using context with timeout of 3 second
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the delete sql statment
	stmt := `DELETE FROM followers WHERE (user_id=$1 AND author_id=$2)`

	// executing the query
	_, err := m.DB.ExecContext(ctx, stmt, user_id, author_id)

	// returns nil if success else returns error
	return err
}

// UpdateFollower updates the Follower
// Takes update value Follower model and previous user id and author id
func (m *postgresDBRepo) UpdateFollower(u *models.Follower, user_id, author_id int) error {

	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the query statement for update follower
	stmt := `
		UPDATE followers
		SET user_id = $3, author_id = $4
		WHERE (user_id = $1 AND author_id = $2)
	`

	// Executing the sql query
	_, err := m.DB.ExecContext(ctx, stmt, user_id, author_id, u.UserID, u.AuthorID)
	if err != nil {
		return err
	}
	return nil
}

// AllReviews fetches all the records from reviews db table.
func (m *postgresDBRepo) AllReviews() ([]*models.Review, error) {
	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the sql statement
	query := `SELECT * FROM reviews`

	// Execute the query using Query Context.
	// If any error occurs, nil and error is returned
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Create a slice of Review model
	reviews := []*models.Review{}

	// Loop through the record.
	for rows.Next() {

		// Store the record in Review instance
		review := new(models.Review)
		if err := rows.Scan(
			&review.ID,
			&review.Rating,
			&review.Body,
			&review.BookID,
			&review.UserID,
			&review.IsActive,
			&review.CreatedAt,
			&review.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// Append the Review instance to the slice of Review
		reviews = append(reviews, review)
	}

	// Return reviews
	return reviews, nil
}

// ReviewExists return true if Review book, review and user  exists else return false.
// It takes Review model instance as parameters
func (m *postgresDBRepo) ReviewExists(u *models.Review) (bool, error) {

	// Creating a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the sql query to check for existing relationship
	query := `
		SELECT COUNT(*) FROM reviews
		WHERE (id=$1 AND book_id=$2 AND user_id=$3)
	`

	// intializing a count variable that stores the no of records
	var count int

	// Executing the query row context and store the total record in count variable.
	// If any error occurs, false and error are returned
	if err := m.DB.QueryRowContext(ctx, query, u.ID, u.BookID, u.UserID).Scan(&count); err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}

	// returning exists is true if count > 0 else retuirn false
	return count > 0, nil
}

// InsertReview add new book user review relation table to db
// Takes Review model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertReview(u *models.Review) error {

	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare a insert query statement
	stmt := `
		INSERT INTO reviews (rating, body, book_id, user_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
	`

	// Executing the query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.Rating,
		u.Body,
		u.BookID,
		u.UserID,
		u.IsActive,
		u.CreatedAt,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// GetReviewByID returns the Review detail from database using id.
// It takes review id as parameters.
// Returns a Review struct instance.
func (m *postgresDBRepo) GetReviewByID(id int) (*models.Review, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM reviews
		WHERE id=$1
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, id)

	// Initializing a Review struct instance
	review := &models.Review{}

	// Scannin the row and storing the result in Review Intance.
	if err := row.Scan(
		&review.ID,
		&review.Rating,
		&review.Body,
		&review.BookID,
		&review.UserID,
		&review.IsActive,
		&review.CreatedAt,
		&review.UpdatedAt,
	); err != nil {
		return nil, err
	}

	// Return a Review Instance and nil
	return review, nil
}

// GetReviewByBookID returns the Review detail from database using book id.
// It takes book id as parameters.
// Returns a Review struct instance.
func (m *postgresDBRepo) GetReviewByBookID(id int) (*models.Review, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM reviews
		WHERE book_id=$1
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, id)

	// Initializing a Review struct instance
	review := &models.Review{}

	// Scannin the row and storing the result in Review Intance.
	if err := row.Scan(
		&review.ID,
		&review.Rating,
		&review.Body,
		&review.BookID,
		&review.UserID,
		&review.IsActive,
		&review.CreatedAt,
		&review.UpdatedAt,
	); err != nil {
		return nil, err
	}

	// Return a Review Instance and nil
	return review, nil
}

// GetReviewByUserID returns the Review detail from database using user id.
// It takes user id as parameters.
// Returns a Review struct instance.
func (m *postgresDBRepo) GetReviewByUserID(id int) (*models.Review, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM reviews
		WHERE user_id=$1
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, id)

	// Initializing a Review struct instance
	review := &models.Review{}

	// Scannin the row and storing the result in Review Intance.
	if err := row.Scan(
		&review.ID,
		&review.Rating,
		&review.Body,
		&review.BookID,
		&review.UserID,
		&review.IsActive,
		&review.CreatedAt,
		&review.UpdatedAt,
	); err != nil {
		return nil, err
	}

	// Return a Review Instance and nil
	return review, nil
}

// DeleteReview deletes the record of Review table from the db.
// It takes book id and user id as parameter
func (m *postgresDBRepo) DeleteReview(id int) error {

	// Using context with timeout of 3 second
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the delete sql statment
	stmt := `DELETE FROM reviews WHERE (id=$1)`

	// executing the query
	_, err := m.DB.ExecContext(ctx, stmt, id)

	// returns nil if success else returns error
	return err
}

// UpdateReview updates the Review
// Takes update value Review model and id of review to be updated as paramaters
func (m *postgresDBRepo) UpdateReview(u *models.Review) error {

	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the query statement for update follower
	stmt := `
		UPDATE reviews
		SET rating = $2, body = $3, book_id = $4, user_id = $5, is_active = $6, updated_at = $7
		WHERE id = $1
	`

	// Executing the sql query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.Rating,
		u.Body,
		u.BookID,
		u.UserID,
		u.IsActive,
		u.UpdatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

// AllContacts fetches all the records from contacts db table.
func (m *postgresDBRepo) AllContacts() ([]*models.Contact, error) {
	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// prepare the sql statement
	query := `SELECT * FROM contacts`

	// Execute the query using Query Context.
	// If any error occurs, nil and error is returned
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	// Create a slice of Review model
	contacts := []*models.Contact{}

	// Loop through the record.
	for rows.Next() {

		// Store the record in Contact instance
		contact := new(models.Contact)
		if err := rows.Scan(
			&contact.ID,
			&contact.FirstName,
			&contact.LastName,
			&contact.Email,
			&contact.Phone,
			&contact.Subject,
			&contact.Message,
			&contact.SubmittedAt,
			&contact.IpAddress,
			&contact.BrowserInfo,
			&contact.ReferringPage,
		); err != nil {
			return nil, err
		}

		// Append the Contact instance to the slice of Contact
		contacts = append(contacts, contact)
	}

	// Return contacts
	return contacts, nil
}

// GetContactByID returns the Contact detail from database using contact id.
// It takes contact id as parameters.
// Returns a Contact struct instance.
func (m *postgresDBRepo) GetContactByID(id int) (*models.Contact, error) {

	// Create timeout of 3 secod with context.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the query statement
	query := `
		SELECT * FROM contacts
		WHERE id=$1
	`

	// Execting the query using row context and returns a row
	row := m.DB.QueryRowContext(ctx, query, id)

	// Initializing a Review struct instance
	contact := &models.Contact{}

	// Scannin the row and storing the result in Review Intance.
	if err := row.Scan(
		&contact.ID,
		&contact.FirstName,
		&contact.LastName,
		&contact.Email,
		&contact.Phone,
		&contact.Subject,
		&contact.Message,
		&contact.Subject,
		&contact.IpAddress,
		&contact.BrowserInfo,
		&contact.ReferringPage,
	); err != nil {
		return nil, err
	}

	// Return a contact Instance and nil
	return contact, nil
}

// DeleteContact deletes the record of Contact table from the db.
// It takes contact id as parameter
func (m *postgresDBRepo) DeleteContact(id int) error {

	// Using context with timeout of 3 second
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Preparing the delete sql statment
	stmt := `DELETE FROM contacts WHERE (id=$1)`

	// executing the query
	_, err := m.DB.ExecContext(ctx, stmt, id)

	// returns nil if success else returns error
	return err
}

// InsertContact add new contact to contacts table to db
// Takes Contact model as a parameter
// Returns an error if something goes wrong
func (m *postgresDBRepo) InsertContact(u *models.Contact) error {

	// Create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Prepare a insert query statement
	stmt := `
		INSERT INTO contacts (first_name, last_name, email, phone, subject, message, submitted_at, ip_address, browser_info, referring_page)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	// Executing the query
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Phone,
		u.Subject,
		u.Message,
		u.SubmittedAt,
		u.IpAddress,
		u.BrowserInfo,
		u.ReferringPage,
	)
	if err != nil {
		return err
	}
	return nil
}
