package DB

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./bot_database.db")
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
