package dbaccess

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func Run() {
	if err := connect(); err != nil {
		log.Fatal(err.Error())
	}
}

// Connect  function will make a connection to the database only once.
func connect() error {
	// Example local: postgres://postgres:password@localhost/DB_1?sslmode=disable
	// Example postgresql://<user>:<pass>@xyz-4782.7tc.cockroachlabs.cloud:26257/horstdb?sslmode=verify-full
	// if you get errors, dry prefixing the database name with the first part of the host (e.g. xyz-4782.horstdb)
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=verify-full",
		os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"),
		os.Getenv("PGHOST"), os.Getenv("PGPORT"),
		os.Getenv("PGDATABASE"),
	)
	fmt.Printf("Connected to %s\n", strings.ReplaceAll(connStr, os.Getenv("PGPASSWORD"), "*****"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// called when function exits
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("The database has been successfully disconnected, goodbye!")
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}

	// this will be printed in the terminal, confirming the connection to the database
	fmt.Println("The database is connected")

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS quotes(quote text,tags text[])"); err != nil {
		return errors.Wrap(err, "create db")
	}
	table := "quotes"
	fmt.Printf("Querying table %s\n", table)
	rows, err := db.Query(`SELECT "quote","tags" FROM ` + table)
	if err != nil {
		return errors.Wrap(err, "query db")
	}
	defer func(rows *sql.Rows) { _ = rows.Close() }(rows)
	var count int
	for rows.Next() {
		count++
		var quote string
		var tags []string
		// https://www.opsdash.com/blog/postgres-arrays-golang.html
		if err = rows.Scan(&quote, pq.Array(&tags)); err != nil {
			return err
		}
		fmt.Printf("Quote #%d: %s tags: %v\n", count, quote, tags)
	}
	if count < 1 {
		if _, err := db.Exec(`INSERT INTO quotes(quote, tags) VALUES ('Be yourself; everyone else is already taken.', '{classic}')`); err != nil {
			return errors.Wrap(err, "insert into")
		}
	}

	return nil
}
