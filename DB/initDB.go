package DB

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // Импорт драйвера MySQL
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

	db, err := sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(%s:3306)/", dbUser, dbPass, dbHost))
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Подключение к БД успешно установлено!")

	// Создание новой базы данных
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("База данных успешно создана или уже существует!")

	db.Close()

	// Шаг 2: Подключение к новой базе данных
	db, err = sql.Open(dbDriver, fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPass, dbHost, dbName))

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
