package Handlers

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// AddOrUpdateUserApartment добавляет или обновляет номер квартиры пользователя.
func AddOrUpdateUserApartment(db *sql.DB, chatID int64, apartment string) error {
	// Проверяем, существует ли уже пользователь
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE chat_id = ?)", chatID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %v", err)
	}

	if exists {
		// Если пользователь существует, обновляем его данные
		_, err := db.Exec("UPDATE users SET apartment = ? WHERE chat_id = ?", apartment, chatID)
		if err != nil {
			return fmt.Errorf("error updating user apartment: %v", err)
		}
	} else {
		// Если пользователь не существует, добавляем новую запись
		_, err := db.Exec("INSERT INTO users (chat_id, apartment) VALUES (?, ?)", chatID, apartment)
		if err != nil {
			return fmt.Errorf("error inserting new user: %v", err)
		}
	}

	return nil
}

// GetUserApartmentByChatID возвращает номер квартиры пользователя по его chatID.
func GetUserApartmentByChatID(db *sql.DB, chatID int64) (string, error) {
	var apartmentNumber string
	err := db.QueryRow("SELECT apartment FROM users WHERE chat_id = ?", chatID).Scan(&apartmentNumber)
	if err != nil {
		return "", fmt.Errorf("unable to fetch apartment number for user %d: %v", chatID, err)
	}
	return apartmentNumber, nil
}

func GetWifiCredentialsByUserID(db *sql.DB, userID int64) (string, string, string, error) {
	var login, password string
	// Здесь предполагается, что у вас есть функция, которая возвращает номер дома пользователя
	houseNumber, err := GetUserApartmentByChatID(db, userID)
	if err != nil {
		return "", "", "", err
	}

	err = db.QueryRow("SELECT login, password FROM wifi_credentials WHERE house_name = ?", houseNumber).Scan(&login, &password)
	if err != nil {
		return "", "", "", fmt.Errorf("error fetching WIFI credentials for user %d: %v", userID, err)
	}

	return houseNumber, login, password, nil
}

// GetHouseOccupied проверяет, занят ли дом, по наличию его названия в таблице users.
func GetHouseOccupied(db *sql.DB, houseName string) (isOccupied bool, err error) {
	var apartmentName string

	err = db.QueryRow("SELECT apartment FROM users WHERE apartment = ?", houseName).Scan(&apartmentName)
	if err != nil {
		if err == sql.ErrNoRows {
			// Дом свободен, так как запись не найдена
			isOccupied = false
			err = nil // Очищаем ошибку, так как отсутствие записи не является ошибкой в данном контексте
		} else {
			// Произошла другая ошибка при запросе
			fmt.Printf("Error checking availability of %s: %v\n", houseName, err)
		}
	} else {
		// Дом занят, так как нашлась запись
		isOccupied = true
	}

	return
}

func AddOrder(db *sql.DB, userID int64, houseName, description, status string) error {
	_, err := db.Exec("INSERT INTO orders (user_id, house_name, status, description) VALUES (?, ?, ?, ?)",
		userID, houseName, status, description)
	if err != nil {
		return fmt.Errorf("failed to insert order: %v", err)
	}
	return nil
}
