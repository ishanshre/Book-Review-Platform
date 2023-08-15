package dbrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

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

func (m *postgresDBRepo) BuyListCount(user_id int) (int, error) {
	// create a timeout of 3 second with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	query := `SELECT COUNT(*) FROM buy_lists WHERE user_id = $1`
	if err := m.DB.QueryRowContext(ctx, query, user_id).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (m *postgresDBRepo) GetAllBooksFromBuyListByUserId(limit, page, user_id int, searchKey, sort string) (*models.BookApiFilter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	query := `
		SELECT 
			COALESCE(b.id, 0) AS b_id,
			COALESCE(b.title, '') AS b_title,
			COALESCE(b.isbn, 0) AS b_isbn,
			COALESCE(b.cover, '') AS b_cover
		FROM 
			buy_lists AS bl
		LEFT JOIN
			books AS b ON b.id = bl.book_id
		where bl.user_id = $1
	`
	countQuery := `
			SELECT 
				COUNT(*)
			FROM 
				buy_lists AS bl
			LEFT JOIN
				books AS b ON b.id = bl.book_id
			where bl.user_id = $1
	`
	if searchKey != "" {
		query = fmt.Sprintf("%s AND (b.title LIKE '%%%s%%' OR CAST(b.isbn AS TEXT) LIKE '%%%s%%')", query, searchKey, searchKey)
		countQuery = fmt.Sprintf("%s AND (b.title LIKE '%%%s%%' OR CAST(b.isbn AS TEXT) LIKE '%%%s%%')", countQuery, searchKey, searchKey)
	}
	if sort != "" {
		query = fmt.Sprintf("%s ORDER BY b.title %s", query, sort)
	}

	var count int
	if err := m.DB.QueryRowContext(ctx, countQuery, user_id).Scan(&count); err != nil {
		return nil, err
	}

	query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)
	rows, err := m.DB.QueryContext(ctx, query, user_id)
	if err != nil {
		return nil, err
	}
	books := []*models.Book{}
	for rows.Next() {
		book := &models.Book{}
		if err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Isbn,
			&book.Cover,
		); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	last_page := m.CalculateLastPage(limit, count)
	return &models.BookApiFilter{
		Total:    count,
		LastPage: last_page,
		Page:     page,
		Books:    books,
	}, nil
}

func (m *postgresDBRepo) BuyListFilter(limit, page int, searchKey, sort string) (*models.BuyListFilterApi, error) {
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
		SELECT u.id, u.username, b.id, b.title, bl.created_at
		FROM buy_lists as bl
		JOIN
			users AS u ON u.id = bl.user_id
		JOIN
			books AS b ON b.id = bl.book_id
	`
	countSql := `
	SELECT 
		COUNT(*)
		FROM buy_lists as bl
		JOIN
			users AS u ON u.id = bl.user_id
		JOIN
			books AS b ON b.id = bl.book_id
	`
	if searchKey != "" {
		sql = fmt.Sprintf("%s WHERE b.title LIKE '%%%s%%' OR u.username LIKE '%%%s%%'", sql, searchKey, searchKey)
		countSql = fmt.Sprintf("%s WHERE b.title LIKE '%%%s%%' OR u.username LIKE '%%%s%%'", countSql, searchKey, searchKey)

	}
	if sort != "" {
		sql = fmt.Sprintf("%s ORDER BY bl.created_at %s", sql, sort)
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
	buyListFilters := []*models.BuyListFilter{}
	for rows.Next() {
		buyListFilter := &models.BuyListFilter{}
		if err := rows.Scan(
			&buyListFilter.UserID,
			&buyListFilter.Username,
			&buyListFilter.BookID,
			&buyListFilter.BookTitle,
			&buyListFilter.CreatedAt,
		); err != nil {
			return nil, err
		}
		buyListFilters = append(buyListFilters, buyListFilter)
	}
	lastPage := m.CalculateLastPage(limit, count)
	return &models.BuyListFilterApi{
		Total:          count,
		Page:           page,
		LastPage:       lastPage,
		BuyListFilters: buyListFilters,
	}, nil
}
