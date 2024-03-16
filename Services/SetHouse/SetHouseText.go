package SetHouse

import (
	"database/sql"
	"fmt"
	tb "gopkg.in/telebot.v3"
	"tgBotElgora/DB/Handlers"
)

var (
	shaleName       = "Шале"
	shaleBDName     = "shale"
	miniHouseName   = "Мини-дом"
	miniHouseBDName = "mini"
	chooseShale     = "Выберите номер шале"
	chooseMiniHouse = "Выберите номер минидома"

	backButton = "Назад"
)

var SHALE_AMOUNT = 6
var MINI_HOUSE_AMOUNT = 3

type HouseInlineButton struct {
	Unique string
	Text   string
	Data   string
}

type HouseInlineKeyboard struct {
	InlineKeyboard [][]tb.InlineButton
}

func getHouseNames(db *sql.DB, key string) []HouseInlineButton {
	var names []HouseInlineButton
	var houseAmount int
	var displayName string

	switch key {
	case shaleBDName:
		houseAmount = SHALE_AMOUNT
		displayName = shaleName
	case miniHouseBDName:
		houseAmount = MINI_HOUSE_AMOUNT
		displayName = miniHouseName
	default:
		return []HouseInlineButton{}
	}

	// Запрашиваем доступность каждого дома из базы данных
	for i := 1; i <= houseAmount; i++ {
		houseIdentifier := fmt.Sprintf("%s_%s_%d", "house", key, i)
		houseData := fmt.Sprintf("%s_%d", key, i)
		houseDisplayName := fmt.Sprintf("%s %d", displayName, i)
		isOccupied, err := Handlers.GetHouseOccupied(db, houseData)
		fmt.Println(isOccupied, err)

		if err != nil && err != sql.ErrNoRows {
			fmt.Println(err)
			continue // В случае ошибки пропускаем дом
		}

		if !isOccupied {
			names = append(names, HouseInlineButton{
				Data:   houseData,
				Text:   houseDisplayName,
				Unique: houseIdentifier,
			})
		}
	}

	return names
}
