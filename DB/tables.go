package DB

import (
	"database/sql"
	"fmt"
	"log"
	config "tgBotElgora/Config"
)

var (
	SHALE_AMOUNT = 6
	MINI_AMOUNT  = 3

	SHALE_NAME = "Шале %d"
	MINI_NAME  = "Мини-дом %d"
)

func fillWifiCredentials(db *sql.DB) {
	// Заполняем пароли шале
	for i := 1; i <= SHALE_AMOUNT; i++ {
		loginKey := config.ConfigKey(fmt.Sprintf("SHALE_%d_LOGIN", i))
		passwordKey := config.ConfigKey(fmt.Sprintf("SHALE_%d_PASSWORD", i))

		login := config.GetConfig(loginKey)
		password := config.GetConfig(passwordKey)

		_, err := db.Exec("INSERT INTO wifi_credentials (house_name, login, password) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE login = VALUES(login), password = VALUES(password)", fmt.Sprintf("shale_%d", i), login, password)
		if err != nil {
			log.Fatalf("Failed to insert WIFI credentials for house %d: %v", i, err)
		}
	}

	// Заполняем пароли минидомиков
	for i := 1; i <= MINI_AMOUNT; i++ {
		loginKey := config.ConfigKey(fmt.Sprintf("MINI_HOUSE_%d_LOGIN", i))
		passwordKey := config.ConfigKey(fmt.Sprintf("MINI_HOUSE_%d_PASSWORD", i))

		login := config.GetConfig(loginKey)
		password := config.GetConfig(passwordKey)

		_, err := db.Exec("INSERT INTO wifi_credentials (house_name, login, password) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE login = VALUES(login), password = VALUES(password)", fmt.Sprintf("mini_%d", i), login, password)
		if err != nil {
			log.Fatalf("Failed to insert WIFI credentials for house %d: %v", i, err)
		}
	}
}

func createWiFiTable(db *sql.DB) (err error) {
	// Создаем и заполняем таблицу wifi_credentials, если она не существует
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS wifi_credentials (
		house_name VARCHAR(255) PRIMARY KEY, -- Указываем максимальную длину ключа
		login TEXT NOT NULL,
		password TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatalf("Failed to create wifi_credentials table: %v", err)
		return
	}

	fillWifiCredentials(db)

	return
}

func createUserTable(db *sql.DB) (err error) {
	// Создание таблицы юзеров, если ее не существуют
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTO_INCREMENT,
        chat_id INTEGER UNIQUE NOT NULL,
        apartment TEXT,
        phone_number TEXT
    );`)
	if err != nil {
		log.Fatal(err)
	}

	return
}

func createOrderTabled(db *sql.DB) (err error) {
	// Создание таблицы заказов, если ее не существует
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    user_id INTEGER NOT NULL,
    house_name TEXT NOT NULL,
    status TEXT NOT NULL,
    description TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);`)
	if err != nil {
		log.Fatal(err)
	}

	return
}
