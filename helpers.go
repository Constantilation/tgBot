package main

import (
	"log"
	config "tgBotElgora/Config"
)

func isAdmin(userID int64) bool {
	adminID, err := config.GetInt64Config(config.AdminId)

	if err != nil {
		log.Fatal("ADMIN_ID is not set in .env")
	}

	return userID == adminID
}
