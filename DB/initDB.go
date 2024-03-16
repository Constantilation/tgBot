package DB

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Импорт драйвера MySQL
	"log"
	"os"
)

func InitDB() *sql.DB {
	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	if dbDriver == "" || dbUser == "" || dbPass == "" || dbName == "" {
		log.Fatal("os env param error", dbDriver, dbUser, dbPass, dbName)
	}

	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPass, dbHost, dbName))
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
