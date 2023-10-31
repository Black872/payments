package db

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"payments/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/zeebo/errs"
)

var dbErr = errs.Class("database")

func NewPostgresDB(cfg string) (*sql.DB, error) {
	log.Print("connecting database... ")
	db, err := sql.Open("pgx", cfg)
	if err != nil {
		return nil, dbErr.Wrap(err)
	}

	// Test connection.
	if err = db.Ping(); err != nil {
		return nil, dbErr.Wrap(err)
	}
	log.Println("OK")

	// If database is empty create schema.
	empty, err := isEmpty(db)
	if err != nil {
		return nil, dbErr.Wrap(err)
	}
	if empty {
		if err = CreateSchema(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func DatabaseURL() (url string) {
	url = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.DBUser(),
		config.DBPassword(),
		config.DBAddress(),
		config.DBPort(),
		config.DBName(),
	)
	return url
}

func CreateSchema(db *sql.DB) error {
	log.Println("creating database schema... ")
	schemaErr := errs.Class("schema error")

	// Open file with database schema.
	file, err := os.Open(config.DBSchemaPath())
	if err != nil {
		return schemaErr.Wrap(err)
	}
	defer file.Close()

	// Scan the query from the file.
	var createSchemaQuery string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		createSchemaQuery += scanner.Text() + "\n"
	}
	if err = scanner.Err(); err != nil {
		return schemaErr.Wrap(err)
	}

	// Execute the query.
	if _, err := db.Exec(createSchemaQuery); err != nil {
		return schemaErr.Wrap(err)
	}
	log.Println("OK")

	return nil
}

// Check whether the database is empty.
func isEmpty(db *sql.DB) (bool, error) {
	testQuery := `SELECT COUNT(*) 
		FROM information_schema.tables 
		WHERE table_schema = 'public';`

	// Read the number of tables from the database.
	var count int
	if err := db.QueryRow(testQuery).Scan(&count); err != nil {
		return false, err
	}

	// If the number of tables is 0 then db is empty.
	return count == 0, nil
}
