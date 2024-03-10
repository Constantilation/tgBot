package SetHouse

import (
	"database/sql"
	"fmt"
	"tgBotElgora/DB/Handlers"
)

var (
	shaleName     = "Шале"
	miniHouseName = "Мини-дом"

	chooseShale     = "Выберите номер шале"
	chooseMiniHouse = "Выберите номер минидома"

	backButton = "Назад"
)

var SHALE_AMOUNT = 6
var MINI_HOUSE_AMOUNT = 3

func getHouseNames(db *sql.DB, key string) []string {
	var names []string
	var houseAmount int

	switch key {
	case shaleName:
		houseAmount = SHALE_AMOUNT
	case miniHouseName:
		houseAmount = MINI_HOUSE_AMOUNT
	default:
		return []string{}
	}

	// Запрашиваем доступность каждого дома из базы данных
	for i := 1; i <= houseAmount; i++ {
		houseName := fmt.Sprintf("%s %d", key, i)
		isOccupied, err := Handlers.GetHouseOccupied(db, houseName)

		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err)
			continue // В случае ошибки пропускаем дом
		}

		if !isOccupied {
			names = append(names, houseName)
		}
	}

	return names
}
