package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
	"golang.org/x/crypto/bcrypt"
)

// AllUsers returns all the users with all access level
func (m *postgresDBRepo) AllUsers(limit, offset int) ([]*models.User, error) {
	// creating database transcation atomic with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query stores the sql query statement
	query := `
		SELECT id, first_name, last_name, username, email, access_level, created_at
		FROM users
		LIMIT=$1 OFFSET=$2
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
			&user.Email,
			&user.AccessLevel,
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
	query := `
		SELECT * 
		FROM users
		WHERE id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	var u *models.User
	if err := row.Scan(
		&u.ID,
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
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLogin,
	); err != nil {
		return nil, fmt.Errorf("cound not fetch id %d from database: %s", id, err)
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
		SET first_name = $2, last_name = $3, email = $4, gender = $5, address = $6, phone = $7, profile_pic = $8, updated_at = $9
		WHERE id = $1
	`
	res, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Username,
		u.Gender,
		u.Address,
		u.Phone,
		u.ProfilePic,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("cannot update the user with id %d : %s", u.ID, err)
	}
	rows_affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows_affected == 0 {
		return fmt.Errorf("could not update user: %s", err)
	}
	return nil
}

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
// Returns id, hashed password and error
func (m *postgresDBRepo) Authenticate(username, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string
	query := `SELECT id, password FROM users WHERE username=$1`
	row := m.DB.QueryRowContext(ctx, query, username)
	if err := row.Scan(&id, &hashedPassword); err != nil {
		return id, "", err
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", fmt.Errorf("incorrect password: %s", bcrypt.ErrMismatchedHashAndPassword)
	} else if err != nil {
		return 0, "", err
	}
	return id, "auth", nil
}
