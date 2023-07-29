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
		&u.DateOfBirth,
		&u.DocumentType,
		&u.DocumentNumber,
		&u.DocumentFront,
		&u.DocumentBack,
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
		WHERE (access_level = $1 AND id= $1)
	`
	row := m.DB.QueryRowContext(ctx, query, id)
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

// GetGlobalUserByID return user by id
func (m *postgresDBRepo) GetGlobalUserByIDAny(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT id, username, first_name, last_name, gender, address, profile_pic, created_at
		FROM users where id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	u := &models.User{}
	if err := row.Scan(
		&u.ID,
		&u.Username,
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
		SET first_name = $2, last_name = $3, email = $4, gender = $5, address = $6, phone = $7, dob = $8, access_level = $9, is_validated = $10, updated_at = $11
		WHERE id = $1
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Gender,
		u.Address,
		u.Phone,
		u.DateOfBirth,
		u.AccessLevel,
		u.IsValidated,
		u.UpdatedAt,
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
		INSERT INTO users (first_name, last_name, email, username, password,address, gender, phone, profile_pic, dob, document_number, document_front, document_back, created_at, updated_at, last_login)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
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
		time.Time{},
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

// AdminInsertsUser insert user to db by admin
func (m *postgresDBRepo) AdminInsertUser(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	stmt := `
		INSERT INTO users (first_name, last_name, email, username, password, address, gender, phone, profile_pic, dob, document_number, document_front, document_back, created_at, updated_at, last_login)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
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
		time.Time{},
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
		SELECT first_name, last_name, email, username, gender, address, phone, profile_pic, dob, document_type, document_number, document_front, document_back, created_at, updated_at, last_login
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
		&u.DateOfBirth,
		&u.DocumentType,
		&u.DocumentNumber,
		&u.DocumentFront,
		&u.DocumentBack,
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
