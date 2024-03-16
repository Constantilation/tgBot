package main

import (
	"fmt"
	"github.com/joho/godotenv"
	tb "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	config "tgBotElgora/Config"
	"tgBotElgora/DB"
	"tgBotElgora/DB/Handlers"
	"tgBotElgora/Helpers"
	"tgBotElgora/Services/AdminPanel"
	"tgBotElgora/Services/MainMenu"
	"tgBotElgora/Services/SetHouse"
	"time"
)

var selectedHouses = make(map[int64]string)

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

	mainMenu := MainMenu.SetupMainMenuHandlers(b, db)

	var mainMenuHandler = func(c tb.Context) error {
		userID := c.Sender().ID
		_, err := Handlers.GetUserApartmentByChatID(db, userID)

		if err != nil {
			setHouse := SetHouse.SetupSetHouseHandlers(b, mainMenu, db)
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

		if !Helpers.IsAdmin(userID) {
			return mainMenuHandler(c)
		}
		adminPanel := AdminPanel.CreateAdminPanelHandlers(b, db)

		// Здесь можно добавить логику для отображения админской панели,
		// например, отправить сообщение с кнопками для управления заказами и пользователями
		return c.Send("Админ. панель", adminPanel)
	})

	b.Handle(tb.OnText, func(c tb.Context) error {
		text := c.Text()
		switch text {
		case AdminPanel.AdminPanelName:
			{
				if !Helpers.IsAdmin(c.Chat().ID) {
					return mainMenuHandler(c)
				} else {
					adminPanel := AdminPanel.CreateAdminPanelHandlers(b, db)
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

			if strings.Contains(parts[0], "order_done") {
				return AdminPanel.OrderDoneHandler(parts[1], db, c)
			}

			if strings.Contains(parts[0], "evict") {
				userIDStr := strings.TrimPrefix(parts[1], "evict_")
				userID, err := strconv.ParseInt(userIDStr, 10, 64)
				if err != nil {
					log.Println("Ошибка при конвертации ID пользователя:", err)
					return c.Respond(&tb.CallbackResponse{Text: "Произошла ошибка при обработке запроса."})
				}
				return AdminPanel.EvictHandler(c, db, userID) // Передаем userID в функцию выселения
			}

			if strings.Contains(parts[0], "house") {
				selectedHouses[c.Chat().ID] = parts[1]
				contactRequest := "Пожалуйста, поделитесь своим контактом, нажав на кнопку ниже:"
				contactButton := tb.ReplyButton{Text: "Поделиться своим контактом", Contact: true}
				c.Send(contactRequest, &tb.ReplyMarkup{ReplyKeyboard: [][]tb.ReplyButton{{contactButton}}, ResizeKeyboard: true})
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

	b.Handle(tb.OnContact, func(c tb.Context) error {
		phoneNumber := c.Message().Contact.PhoneNumber
		houseName := selectedHouses[c.Chat().ID]

		SetHouse.HouseNumber(c, houseName, phoneNumber, db)

		if err != nil {
			return c.Send("Произошла ошибка при сохранении номера квартиры.")
		}

		houseName, _ = Helpers.GetHouseName(houseName)
		c.Edit(fmt.Sprintf("Номер дома сохранен: %s. Теперь вы можете использовать другие команды.", houseName))
		return c.Send("Главное меню", mainMenu)
	})

	b.Start()
}
