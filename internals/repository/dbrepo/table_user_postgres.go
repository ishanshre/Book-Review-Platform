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
		SELECT id, username, access_level, created_at
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
			&user.Username,
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
		SELECT id, username, email
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
		&u.ID,
		&u.Username,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
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
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE (access_level = $1 AND id= $1)
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	var u *models.User
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
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
		SELECT id, username, email, created_at, updated_at
		FROM users where id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	u := &models.User{}
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
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
		SET email = $2, access_level = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.ID,
		u.Email,
		u.AccessLevel,
		u.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("cannot update the user with id %d : %s", u.ID, err)
	}
	return nil
}

// InsertUser insert new user into database.
// This method is used for new user sign up
func (m *postgresDBRepo) InsertUser(u *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	db, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			db.Rollback()
			panic(p)
		} else if err != nil {
			db.Rollback()
		} else {
			err = db.Commit()
		}
	}()
	stmt := `
		INSERT INTO users (email, username, password, created_at, updated_at, last_login)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	res := db.QueryRowContext(
		ctx,
		stmt,
		u.Email,
		u.Username,
		u.Password,
		time.Now(),
		time.Now(),
		time.Time{},
	)
	if err != nil {
		return fmt.Errorf("could not create new user: %s", err)
	}
	var id int
	if err := res.Scan(&id); err != nil {
		return err
	}
	kycquery := `
		INSERT INTO kycs (user_id, first_name, last_name, gender, address, phone, profile_pic, dob, document_number, document_front, document_back, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	kycres, err := db.ExecContext(
		ctx,
		kycquery,
		id,
		"",
		"",
		"Unknown",
		"",
		"",
		"",
		time.Now().Format(time.DateOnly),
		"",
		"",
		"",
		time.Now(),
	)
	if err != nil {
		return err
	}
	rows_affected, _ := kycres.RowsAffected()
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
		INSERT INTO users (email, username, password, created_at, updated_at, last_login)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	res, err := m.DB.ExecContext(
		ctx,
		stmt,
		u.Email,
		u.Username,
		u.Password,
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
func (m *postgresDBRepo) Authenticate(username, testPassword string) (int, int, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var access_level int
	var hashedPassword string
	var is_validated bool
	query := `
		SELECT u.id, u.password, u.access_level, k.is_validated 
		FROM users AS u
		JOIN
			kycs AS k ON u.id = k.user_id
		WHERE u.username=$1
	`
	row := m.DB.QueryRowContext(ctx, query, username)
	if err := row.Scan(&id, &hashedPassword, &access_level, &is_validated); err != nil {
		return id, 2, false, err
	}
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, 2, false, fmt.Errorf("incorrect password: %s", bcrypt.ErrMismatchedHashAndPassword)
	} else if err != nil {
		return 0, 2, false, err
	}
	return id, access_level, is_validated, nil
}

// Get information for personal profile page
func (m *postgresDBRepo) GetProfilePersonal(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
		SELECT email, username, dcreated_at, updated_at, last_login
		FROM users
		WHERE id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)
	u := &models.User{}
	if err := row.Scan(
		&u.Email,
		&u.Username,
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

// UserListFilter
func (m *postgresDBRepo) UserListFilter(limit, page int, searchKey, sort string) (*models.AdminUserListApi, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	sql := `
	SELECT 
		u.id, 
		u.username, 
		u.access_level, 
		u.created_at, 
		COALESCE(k.is_validated, false)
	FROM 
		users AS u 
	LEFT JOIN
		kycs AS k ON k.user_id = u.id
	`
	countSql := `
	SELECT 
		COUNT(*)
	FROM 
		users AS u 
	LEFT JOIN
		kycs AS k ON k.user_id = u.id
	`
	if searchKey != "" {
		sql = fmt.Sprintf("%s WHERE u.username LIKE '%%%s%%' OR u.email LIKE '%%%s%%'", sql, searchKey, searchKey)
		countSql = fmt.Sprintf("%s WHERE u.username LIKE '%%%s%%' OR u.email LIKE '%%%s%%'", countSql, searchKey, searchKey)
	}
	if sort != "" {
		sql = fmt.Sprintf("%s ORDER BY u.username %s", sql, sort)
		// countSql = fmt.Sprintf("%s ORDER BY u.username %s", countSql, sort)
	}
	var count int
	if err := m.DB.QueryRowContext(ctx, countSql).Scan(&count); err != nil {
		return nil, err
	}
	sql = fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, limit, offset)
	rows, err := m.DB.QueryContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	users := []*models.AdminUserList{}
	for rows.Next() {
		user := &models.AdminUserList{}
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.AccessLevel,
			&user.CreatedAt,
			&user.IsValidated,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	lastPage := m.CalculateLastPage(limit, count)
	return &models.AdminUserListApi{
		Total:    count,
		Page:     page,
		LastPage: lastPage,
		Users:    users,
	}, nil
}
