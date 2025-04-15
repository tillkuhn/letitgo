package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

// Run launches example based on
// * https://www.codeproject.com/Articles/5261771/Golang-SQLite-Simple-Example
// * https://antonz.org/json-virtual-columns/
func Run() {
	// os.Remove("sqlite-database.db") // obsolete since Create truncates if the file already exists

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := file.Close(); err != nil {
		return
	}

	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer func(sqliteDatabase *sql.DB) {
		_ = sqliteDatabase.Close()
	}(sqliteDatabase) // Defer Closing the database
	createTable(sqliteDatabase) // Create Database Tables

	// INSERT RECORDS
	insertStudent(sqliteDatabase, "0001", "Liana Kim", "Bachelor")
	insertStudent(sqliteDatabase, "0002", "Glen Rangel", "Bachelor")
	insertStudent(sqliteDatabase, "0003", "Martin Martins", "Master")
	insertStudent(sqliteDatabase, "0004", "Alayna Ar", "PHD")
	insertStudent(sqliteDatabase, "0005", "Marni Benson", "Bachelor")
	insertStudent(sqliteDatabase, "0006", "Derrick Griffiths", "Master")
	insertStudent(sqliteDatabase, "0007", "Leigh Daly", "Bachelor")
	insertStudent(sqliteDatabase, "0008", "Marni Benson", "PHD")
	insertStudent(sqliteDatabase, "0009", "Hase Klaus", "Bachelor")

	// DISPLAY INSERTED RECORDS
	displayStudents(sqliteDatabase)
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE student (
		"idStudent" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"code" TEXT,
		"name" TEXT,
		"program" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create student table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	if _, err = statement.Exec(); err != nil { // Execute SQL Statements
		log.Println(err.Error())
	} else {
		log.Println("student table created")
	}
}

// We are passing db reference connection from main to our method with other parameters
func insertStudent(db *sql.DB, code string, name string, program string) {
	log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO student(code, name, program) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(code, name, program)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayStudents(db *sql.DB) {
	row, err := db.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer func(row *sql.Rows) { _ = row.Close() }(row)
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var code string
		var name string
		var program string
		if err = row.Scan(&id, &code, &name, &program); err != nil {
			log.Println(err.Error())
		} else {
			log.Println("Student: ", code, " ", name, " ", program)
		}
	}
}
