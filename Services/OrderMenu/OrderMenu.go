package OrderMenu

import (
	"database/sql"
	"fmt"
	tb "gopkg.in/telebot.v3"
	"log"
	config "tgBotElgora/Config"
	"tgBotElgora/DB/Handlers"
)

func createOrderServices() (*tb.ReplyMarkup, []tb.Btn) {
	orderServices := &tb.ReplyMarkup{ResizeKeyboard: true}
	btns := []tb.Btn{
		orderServices.Text(bathVatName),         // 0
		orderServices.Text(extraCleaningName),   // 1
		orderServices.Text(hookahName),          // 2
		orderServices.Text(airportTransferName), // 3
		orderServices.Text(igniteGrillName),     // 4
		orderServices.Text(backToMenuName),      // 5
	}

	orderServices.Reply(
		orderServices.Row(btns[0], btns[1]),
		orderServices.Row(btns[2], btns[3]),
		orderServices.Row(btns[4], btns[5]),
	)

	return orderServices, btns
}

func orderButtonHandler(b *tb.Bot, serviceName string, unique string, db *sql.DB) tb.InlineButton {
	orderButton := tb.InlineButton{
		Unique: unique,
		Text:   "Заказать",
		Data:   serviceName,
	}

	fmt.Println(serviceName, "here 2", orderButton)

	b.Handle(&orderButton, func(c tb.Context) error {
		adminID, err := config.GetInt64Config(config.AdminId)
		if err != nil {
			log.Fatal("ADMIN_ID is not set in .env")
		}

		// Логика выполнения заказа
		userID := c.Sender().ID
		apartmentNumber, err := Handlers.GetUserApartmentByChatID(db, userID)
		if err != nil {
			log.Println(err)
			return c.Send("Пожалуйста, сначала введите номер вашей квартиры.")
		}

		err = Handlers.AddOrder(db, userID, apartmentNumber, serviceName, "pending")
		if err != nil {
			log.Println(err)
			return c.Send("Произошла ошибка при сохранении заказа.")
		}

		orderMessage := fmt.Sprintf("Заказ услуги: %s. Дом %s", serviceName, apartmentNumber)
		b.Send(tb.ChatID(adminID), orderMessage)
		return c.Send(orderSent)
	})

	return orderButton
}

func SetupOrderHandlers(b *tb.Bot, returnToMainMenu *tb.ReplyMarkup, db *sql.DB) *tb.ReplyMarkup {
	orderService, btns := createOrderServices()

	b.Handle(&btns[0], func(c tb.Context) error {
		orderButton := orderButtonHandler(b, bathVatNameService, "order_bath_vat", db)

		return c.Send(bathVatHandlerText, &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{orderButton}}})
	})

	b.Handle(&btns[1], func(c tb.Context) error {
		orderButton := orderButtonHandler(b, extraCleaningNameService, "order_extra_cleaning", db)

		return c.Send(extraCleaningHandlerText, &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{orderButton}}})
	})

	b.Handle(&btns[2], func(c tb.Context) error {
		orderButton := orderButtonHandler(b, hookahName, "order_hookah", db)

		return c.Send(hookahHandlerText, &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{orderButton}}})
	})

	b.Handle(&btns[3], func(c tb.Context) error {
		orderButton := orderButtonHandler(b, airportTransferName, "order_transfer", db)

		return c.Send(airportTransferHandlerText, &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{orderButton}}})
	})

	b.Handle(&btns[4], func(c tb.Context) error {
		orderButton := orderButtonHandler(b, igniteGrillName, "order_grill", db)

		return c.Send(igniteGrillHandlerText, &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{orderButton}}})
	})

	b.Handle(&btns[5], func(c tb.Context) error {
		fmt.Println(returnToMainMenu)
		return c.Send("Главное меню", returnToMainMenu)
	})

	return orderService
}
