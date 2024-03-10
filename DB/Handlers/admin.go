package Handlers

import (
	"database/sql"
	"fmt"
)

type Tenant struct {
	UserID    int64
	HouseName string
}

type Order struct {
	ID          int
	HouseName   string
	Description string
}

func GetTenants(db *sql.DB) ([]Tenant, error) {
	rows, err := db.Query("SELECT chat_id, apartment FROM users WHERE apartment IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []Tenant
	for rows.Next() {
		var t Tenant
		if err := rows.Scan(&t.UserID, &t.HouseName); err != nil {
			return nil, err
		}
		tenants = append(tenants, t)
	}

	return tenants, nil
}

func EvictTenant(db *sql.DB, userID int64) error {
	fmt.Println(userID)
	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Удаляем все заказы, связанные с пользователем
	if _, err := tx.Exec("DELETE FROM orders WHERE user_id = ?", userID); err != nil {
		tx.Rollback() // Откатываем транзакцию в случае ошибки
		return err
	}

	// Удаляем пользователя
	if _, err := tx.Exec("DELETE FROM users WHERE chat_id = ?", userID); err != nil {
		tx.Rollback() // Откатываем транзакцию в случае ошибки
		return err
	}

	// Подтверждаем транзакцию
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func GetHousesWithActiveOrders(db *sql.DB) ([]string, error) {
	var houses []string

	rows, err := db.Query("SELECT DISTINCT house_name FROM orders WHERE status = 'pending'")
	if err != nil {
		return nil, fmt.Errorf("querying distinct house names with active orders: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var houseName string
		if err := rows.Scan(&houseName); err != nil {
			return nil, fmt.Errorf("scanning house name: %w", err)
		}
		houses = append(houses, houseName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over houses result set: %w", err)
	}

	return houses, nil
}

func GetActiveOrdersForHouse(db *sql.DB, houseName string) ([]Order, error) {
	var orders []Order

	rows, err := db.Query("SELECT id, house_name, description FROM orders WHERE house_name = ? AND status = 'pending'", houseName)
	if err != nil {
		return nil, fmt.Errorf("querying active orders for house: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.HouseName, &order.Description); err != nil {
			return nil, fmt.Errorf("scanning order: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating over orders result set: %w", err)
	}

	return orders, nil
}

func MarkOrderAsDone(db *sql.DB, orderID int) error {
	_, err := db.Exec("DELETE FROM orders WHERE id = ?", orderID)
	if err != nil {
		return fmt.Errorf("deleting order with ID %s: %w", orderID, err)
	}
	return nil
}
