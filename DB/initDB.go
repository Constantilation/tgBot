package DB

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func InitDB() *sql.DB {
	dbPath := os.Getenv("DB_PATH")

	if dbPath == "" {
		log.Fatal("DB_PATH не задан", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	err = createWiFiTable(db)

	if err != nil {
		log.Fatalf("Failed to create wifi_credentials table: %v", err)
	}

	err = createUserTable(db)

	if err != nil {
		log.Fatal(err)
	}

	err = createOrderTabled(db)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
