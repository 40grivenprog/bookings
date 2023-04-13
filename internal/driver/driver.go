package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/jackc/pgx/v4"
)

// DB holds db connection
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}


const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

// Creates db pool for psql
func ConnectSql(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.SQL = db

	err = testDB(db)

	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// Tries to ping db
func testDB(d *sql.DB) error {
	err := d.Ping()

	if err != nil {
		return err
	}

	return nil
}

// Creates new DB for app
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	
	if err != nil {
		return nil, err
	}

	if err = testDB(db); err != nil {
		return nil, err
	}

	return db, nil
}
