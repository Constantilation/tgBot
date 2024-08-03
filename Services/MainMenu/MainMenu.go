package MainMenu

import (
	"database/sql"
	"fmt"
	"log"
	config "tgBotElgora/Config"
	"tgBotElgora/DB/Handlers"
	"tgBotElgora/Helpers"
	"tgBotElgora/Services/OrderMenu"

	"github.com/joho/godotenv"
	tb "gopkg.in/telebot.v3"
)

func createMainMenu() (*tb.ReplyMarkup, []tb.Btn) {
	mainMenu := &tb.ReplyMarkup{ResizeKeyboard: true}
	btns := []tb.Btn{
		mainMenu.Text(aboutUsName),      // 0
		mainMenu.Text(wiFiName),         // 1
		mainMenu.Text(ourServicesName),  // 2
		mainMenu.Text(livingRulesName),  // 3
		mainMenu.Text(usefulInfoName),   // 4
		mainMenu.Text(contactAdminName), // 5
		mainMenu.Text(reviewUsName),     // 6
	}

	mainMenu.Reply(
		mainMenu.Row(btns[0], btns[1]),
		mainMenu.Row(btns[2], btns[3]),
		mainMenu.Row(btns[4], btns[5]),
		mainMenu.Row(btns[6]),
	)

	return mainMenu, btns
}

func SetupMainMenuHandlers(b *tb.Bot, db *sql.DB) *tb.ReplyMarkup {
	mainMenu, btns := createMainMenu()
	orderService := OrderMenu.SetupOrderHandlers(b, mainMenu, db)

	b.Handle(&btns[0], func(c tb.Context) error {
		return c.Send(aboutUsHandlerText)
	})

	b.Handle(&btns[1], func(c tb.Context) error {
		chatID := c.Sender().ID
		houseNumber, login, password, err := Handlers.GetWifiCredentialsByUserID(db, chatID)

		if err != nil {
			log.Println(err)
			return c.Send(errorText)

		}

		houseName, _ := Helpers.GetHouseName(houseNumber)

		return c.Send(fmt.Sprintf(wiFiText, houseName, login, password))
	})

	b.Handle(&btns[2], func(c tb.Context) error {
		return c.Send(ourServicesHandlerText, orderService)
	})

	b.Handle(&btns[3], func(c tb.Context) error {
		return c.Send(livingRulesHandlerText)
	})

	b.Handle(&btns[4], func(c tb.Context) error {
		return c.Send(usefulInfoHandlerText)
	})

	b.Handle(&btns[5], func(c tb.Context) error {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
			return err
		}

		adminName := config.GetConfig(config.AdminName)
		adminLastName := config.GetConfig(config.AdminLastName)
		adminPhone := config.GetConfig(config.AdminPhone)
		botToken := config.GetConfig(config.TelegramBotToken)

		Helpers.SendContactViaTelegramAPI(botToken, c.Chat().ID, adminName, adminPhone, &adminLastName)

		return c.Send(contactAdminHandlerText)
	})

	b.Handle(&btns[6], func(c tb.Context) error {
		return c.Send(reviewUsHandlerText)
	})

	return mainMenu
}
