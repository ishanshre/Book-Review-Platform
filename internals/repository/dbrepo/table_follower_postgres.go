package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

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

func (m *postgresDBRepo) FollowerCount(user_id int) (int, error) {
	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int

	query := `SELECT COUNT(*) FROM followers WHERE user_id = $1`
	if err := m.DB.QueryRowContext(ctx, query, user_id).Scan(&count); err != nil {
		return 0, nil
	}
	return count, nil
}

func (m *postgresDBRepo) GetAllFollowingsByUserId(user_id int) ([]*models.Author, error) {
	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT a.id, a.first_name, a.last_name
		FROM followers as f
		JOIN authors AS a ON a.id = f.author_id
		WHERE f.user_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	authors := []*models.Author{}
	for rows.Next() {
		author := &models.Author{}
		if err := rows.Scan(
			&author.ID,
			&author.FirstName,
			&author.LastName,
		); err != nil {
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}

func (m *postgresDBRepo) FollowerFilter(limit, page int, searchKey, sort string) (*models.FollowerFilterApi, error) {
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
		SELECT u.id, u.username, a.id, a.first_name, a.last_name, f.followed_at
		FROM followers AS f
		JOIN
			users AS u ON u.id = f.user_id
		JOIN
			authors AS a ON a.id = f.author_id
	`
	countSql := `
		SELECT 
			COUNT(*)
		FROM followers AS f
		JOIN
			users AS u ON u.id = f.user_id
		JOIN
			authors AS a ON a.id = f.author_id
	`
	if searchKey != "" {
		sql = fmt.Sprintf("%s WHERE a.first_name LIKE '%%%s%%' OR a.last_name LIKE '%%%s%%' OR u.username LIKE '%%%s%%'", sql, searchKey, searchKey, searchKey)
		countSql = fmt.Sprintf("%s WHERE a.first_name LIKE '%%%s%%' OR a.last_name LIKE '%%%s%%' OR u.username LIKE '%%%s%%'", countSql, searchKey, searchKey, searchKey)

	}
	if sort != "" {
		sql = fmt.Sprintf("%s ORDER BY f.followed_at %s", sql, sort)
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
	followerFilters := []*models.FollowerFilter{}
	for rows.Next() {
		followerFilter := &models.FollowerFilter{}
		if err := rows.Scan(
			&followerFilter.UserID,
			&followerFilter.Username,
			&followerFilter.AuthorID,
			&followerFilter.AuthorFirstName,
			&followerFilter.AuthorLastName,
			&followerFilter.FollowedAt,
		); err != nil {
			return nil, err
		}
		followerFilters = append(followerFilters, followerFilter)
	}
	lastPage := m.CalculateLastPage(limit, count)
	return &models.FollowerFilterApi{
		Total:           count,
		Page:            page,
		LastPage:        lastPage,
		FollowerFilters: followerFilters,
	}, nil
}
