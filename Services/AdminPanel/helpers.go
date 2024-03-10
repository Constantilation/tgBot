package AdminPanel

import (
	"log"
	config "tgBotElgora/Config"
)

func GetAdminID() int64 {
	adminID, err := config.GetInt64Config(config.AdminId)

	if err != nil {
		log.Fatal("ADMIN_ID is not set in .env")
	}

	return adminID
}
