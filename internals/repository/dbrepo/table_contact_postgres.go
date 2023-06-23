package dbrepo

import (
	"context"
	"time"

	"github.com/ishanshre/Book-Review-Platform/internals/models"
)

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
