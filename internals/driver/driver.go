package driver

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// DB represents a database connection.
type DB struct {
	SQL *sql.DB
}

// dbConn is the global database connection instance.
var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDBConn = 5
const maxLifeDBTime = 5 * time.Minute

// ConnectSQL connects to the SQL database and returns a DB instance.
func ConnectSQL(database string, dsn string) (*DB, error) {
	d, err := NewDatabase(database, dsn)
	if err != nil {
		panic(err)
	}
	// configure database options
	d.SetMaxOpenConns(maxOpenDBConn)    // maximum number of connection to database
	d.SetMaxIdleConns(maxIdleDBConn)    // maximim number of connection in idle connection pool
	d.SetConnMaxLifetime(maxLifeDBTime) // max time before the connection expires
	dbConn.SQL = d

	// test the database connection
	if err := testDB(d); err != nil {
		return nil, err
	}
	return dbConn, nil
}

// tries to ping the database i.e. test the connection to the database
func testDB(d *sql.DB) error {
	if err := d.Ping(); err != nil {
		return err
	}
	return nil
}

// create a new database
func NewDatabase(database string, dsn string) (*sql.DB, error) {
	// connect to the database
	db, err := sql.Open(database, dsn)
	if err != nil {
		return nil, err
	}
	// test the database connection
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
