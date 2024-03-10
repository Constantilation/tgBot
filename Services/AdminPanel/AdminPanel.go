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

func ActiveOrdersHandler(c tb.Context, db *sql.DB, house string) error {
	fmt.Println(house)
	orders, err := Handlers.GetActiveOrdersForHouse(db, house) // Функция для получения заказов для дома
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

func EvictHandler(c tb.Context, db *sql.DB, userId int64) error {
	// Выполняем логику выселения
	err := Handlers.EvictTenant(db, userId)
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
				Data:   "evict_" + strconv.FormatInt(tenant.UserID, 10),
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

		var inlineButtons [][]tb.InlineButton // Используем двумерный слайс для группировки кнопок

		if len(houses) == 0 {
			return c.Send("Нет активных заказов")
		}

		// Создаем inline кнопки для каждого дома
		for _, house := range houses {
			btn := tb.InlineButton{
				Unique: "active_order_" + house, // Уникальный идентификатор кнопки
				Text:   house,                   // Текст кнопки
				// Data используется для передачи информации обратно в коллбек
				Data: house, // Данные кнопки, которые будут использоваться в коллбеке
			}
			// Добавляем кнопку в список
			inlineButtons = append(inlineButtons, []tb.InlineButton{btn})
		}

		// Создаем inline клавиатуру с кнопками
		markup := &tb.ReplyMarkup{InlineKeyboard: inlineButtons}

		// Отправляем сообщение с inline клавиатурой
		err = c.Send(chooseHouseToSeeOrdersText, markup)
		if err != nil {
			log.Println(keyboardMessageError, err)
		}
		return nil
	})

	return adminPanel
}
