package main

import (
	"github.com/joho/godotenv"
	tb "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	config "tgBotElgora/Config"
	"tgBotElgora/DB"
	"tgBotElgora/DB/Handlers"
	"tgBotElgora/Services/AdminPanel"
	"tgBotElgora/Services/SetHouse"
	"time"
)

func main() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
		return
	}

	botToken := config.GetConfig(config.TelegramBotToken)

	b, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	db := DB.InitDB()
	defer db.Close()

	mainMenu := setupHandlers(b, db)
	setHouse := SetHouse.SetupSetHouseHandlers(b, mainMenu, db)
	adminPanel := AdminPanel.CreateAdminPanelHandlers(b, db)

	var mainMenuHandler = func(c tb.Context) error {
		userID := c.Sender().ID
		_, err := Handlers.GetUserApartmentByChatID(db, userID)

		if err != nil {
			AdminPanel.CreateAdminApprovalHandlers(c, b, setHouse)
		} else {
			// Если пользователь найден в базе данных, показать главное меню
			log.Println("User found in database, showing the main menu.")
			return c.Send("Добро пожаловать обратно!", mainMenu)
		}

		return err
	}

	b.Handle("/start", func(c tb.Context) error {
		return mainMenuHandler(c)
	})

	b.Handle("/admin_panel", func(c tb.Context) error {
		userID := c.Sender().ID

		if !isAdmin(userID) {
			return mainMenuHandler(c)
		}

		// Здесь можно добавить логику для отображения админской панели,
		// например, отправить сообщение с кнопками для управления заказами и пользователями
		return c.Send("Админ. панель", adminPanel)
	})

	b.Handle(tb.OnText, func(c tb.Context) error {
		text := c.Text()
		switch text {
		case AdminPanel.AdminPanelName:
			{
				if !isAdmin(c.Chat().ID) {
					return mainMenuHandler(c)
				} else {
					return c.Send(AdminPanel.AdminPanelName, adminPanel)
				}
			}
		default:
			return c.Send("Для вызова меню напишите /start")
		}

		return err
	})

	b.Handle(tb.OnCallback, func(c tb.Context) error {
		data := c.Callback().Data
		parts := strings.Split(data, "|")

		if len(parts) > 1 {
			if strings.HasPrefix(parts[1], "order_done_") {
				return AdminPanel.OrderDoneHandler(parts[1], db, c)
			}

			if strings.HasPrefix(parts[1], "evict_") {
				userIDStr := strings.TrimPrefix(parts[1], "evict_")
				userID, err := strconv.ParseInt(userIDStr, 10, 64)
				if err != nil {
					log.Println("Ошибка при конвертации ID пользователя:", err)
					return c.Respond(&tb.CallbackResponse{Text: "Произошла ошибка при обработке запроса."})
				}
				return AdminPanel.EvictHandler(c, db, userID) // Передаем userID в функцию выселения
			}

			if strings.Contains(parts[0], "active_order") {
				err = AdminPanel.ActiveOrdersHandler(c, db, parts[1])
				c.Edit("Актуальные заказы")
				return c.Respond(&tb.CallbackResponse{Text: "Команда выполнена"})
			}

		} else {
			return c.Send("Неизвестная ошибка")
		}

		return nil
	})

	b.Start()
}
