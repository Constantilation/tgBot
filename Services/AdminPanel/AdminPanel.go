package AdminPanel

import (
	"database/sql"
	"fmt"
	tb "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	"tgBotElgora/DB/Handlers"
)

func showAdminPanel() (tenantsBtn, activeOrdersBtn tb.Btn, markup *tb.ReplyMarkup) {
	markup = &tb.ReplyMarkup{ResizeKeyboard: true}
	tenantsBtn = markup.Text(tenantsBtnName)
	activeOrdersBtn = markup.Text(ActiveOrdersBtnName)
	// Добавьте другие кнопки админской панели здесь

	markup.Reply(
		markup.Row(tenantsBtn),
		markup.Row(activeOrdersBtn),
	)

	return
}

func ActiveOrdersHandler(c tb.Context, db *sql.DB) error {
	orders, err := Handlers.GetActiveOrdersForHouse(db, c.Text()) // Функция для получения заказов для дома
	if err != nil {
		log.Println(GetOrdersForHouseError, err)
		return c.Send(GetOrdersForHouseError)
	}

	for _, order := range orders {
		markup := &tb.ReplyMarkup{}
		doneButton := tb.InlineButton{
			Unique: "order_done", // Уникальный ключ обработчика
			Text:   Done,
			Data:   "order_done_" + strconv.Itoa(order.ID), // Где orderID - это ID заказа
		}
		markup.InlineKeyboard = [][]tb.InlineButton{{doneButton}}

		c.Send(fmt.Sprintf("Заказ %d: %s", order.ID, order.Description), markup)
	}

	return err
}

func EvictHandler(c tb.Context, db *sql.DB) error {
	// Выполняем логику выселения
	err := Handlers.EvictTenant(db, c.Chat().ID)
	if err != nil {
		log.Println(err)
		return c.Respond(&tb.CallbackResponse{Text: tenantEvictionError})
	}

	// Отправляем подтверждение об успешном выселении
	c.Respond(&tb.CallbackResponse{Text: tenantEvicted})

	return c.Edit(tenantEvicted)
}

func OrderDoneHandler(orderId string, db *sql.DB, c tb.Context) error {
	// Извлекаем ID заказа из второй части
	orderIDStr := strings.TrimPrefix(orderId, "order_done_")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		fmt.Println("Ошибка при конвертации ID заказа:", err)
		return c.Respond(&tb.CallbackResponse{Text: orderHandleError})
	}

	// Ваша логика по обработке выполнения заказа
	err = Handlers.MarkOrderAsDone(db, orderID)
	if err != nil {
		fmt.Println(orderHandleError, err)
		return c.Respond(&tb.CallbackResponse{Text: orderHandleError})
	}

	// Отправляем ответ на callback, что заказ выполнен
	c.Respond(&tb.CallbackResponse{Text: orderHandleDone})
	// Опционально, обновляем сообщение для устранения кнопки
	return c.Edit(orderDone)
}

func CreateAdminPanelHandlers(b *tb.Bot, db *sql.DB) *tb.ReplyMarkup {
	tenantsBtn, activeOrdersBtn, adminPanel := showAdminPanel()
	// Обработчик кнопки "Текущие жильцы"
	b.Handle(&tenantsBtn, func(c tb.Context) error {
		// Извлекаем список жильцов из БД
		tenants, err := Handlers.GetTenants(db)
		if err != nil {
			log.Println(err)
			return c.Send(tenantsGetError, adminPanel)
		}

		if len(tenants) == 0 {
			return c.Send(tenantsAreMissing, adminPanel)
		}

		for _, tenant := range tenants {
			evictionBtn := tb.InlineButton{
				Unique: "evict", // Уникальный идентификатор кнопки с использованием названия дома
				Text:   tenantsEvict,
				Data:   "evict_" + tenant.HouseName,
			}
			markup := &tb.ReplyMarkup{}
			markup.InlineKeyboard = [][]tb.InlineButton{{evictionBtn}}

			// Отправляем сообщение с названием дома и кнопкой "Выселение"
			c.Send(fmt.Sprintf("Дом: %s", tenant.HouseName), markup)
		}

		return nil
	})

	b.Handle(&activeOrdersBtn, func(c tb.Context) error {
		houses, err := Handlers.GetHousesWithActiveOrders(db)
		if err != nil {
			log.Println(houseGetError, err)
			return c.Send(houseGetError)
		}

		markup := tb.ReplyMarkup{ResizeKeyboard: true} // Инициализация клавиатуры
		var buttons []tb.Btn                           // Слайс для кнопок

		// Создание кнопок для каждого дома
		for _, house := range houses {
			btn := markup.Text(house)      // Создаем кнопку
			buttons = append(buttons, btn) // Добавляем кнопку в слайс
		}

		// Добавляем кнопку "Назад" отдельно
		backButton := markup.Text(AdminPanelName)
		buttons = append(buttons, backButton)

		// Добавление всех кнопок в клавиатуру
		markup.Reply(
			// Для каждой кнопки создаем отдельную строку
			markup.Row(buttons...),
		)

		// Отправляем сообщение с клавиатурой
		err = c.Send(chooseHouseToSeeOrdersText, &markup)
		if err != nil {
			log.Println(keyboardMessageError, err)
		}
		return nil
	})

	return adminPanel
}
