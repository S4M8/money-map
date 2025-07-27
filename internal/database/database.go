
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	var err error
	dsn := os.Getenv("DB_SOURCE")
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("failed to ping database:", err)
	}

	fmt.Println("Successfully connected to database")
}
