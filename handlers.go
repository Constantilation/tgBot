package main

import (
	"database/sql"
	tb "gopkg.in/telebot.v3"
	"tgBotElgora/Services/MainMenu"
	"tgBotElgora/Services/SetHouse"
)

func setupHandlers(b *tb.Bot, db *sql.DB) *tb.ReplyMarkup {
	mainMenu := MainMenu.SetupMainMenuHandlers(b, db)
	SetHouse.SetupSetHouseHandlers(b, mainMenu, db)

	return mainMenu
}
