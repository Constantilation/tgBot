package SetHouse

import (
	"database/sql"
	"fmt"
	tb "gopkg.in/telebot.v3"
	"tgBotElgora/DB/Handlers"
)

func createSetHouseBtns(db *sql.DB) (
	setHouseBtns, shaleBtns, miniHouseBtns []tb.Btn,
	setHouse, setShale, setMiniHouse *tb.ReplyMarkup,
) {
	setHouse = &tb.ReplyMarkup{ResizeKeyboard: true}
	setShale = &tb.ReplyMarkup{ResizeKeyboard: true}
	setMiniHouse = &tb.ReplyMarkup{ResizeKeyboard: true}

	setHouseBtns = []tb.Btn{
		setHouse.Text(shaleName),     // 0
		setHouse.Text(miniHouseName), // 1
	}
	setHouse.Reply(
		setHouse.Row(setHouseBtns[0]),
		setHouse.Row(setHouseBtns[1]),
	)

	shaleNames := getHouseNames(db, shaleName)
	for _, v := range shaleNames {
		shaleBtns = append(shaleBtns, setShale.Text(v))
	}
	backButtonShale := setShale.Text(backButton)
	shaleBtns = append(shaleBtns, backButtonShale)

	setShale.Reply(
		setShale.Row(shaleBtns[:len(shaleNames)/2]...),
		setShale.Row(shaleBtns[len(shaleNames)/2:len(shaleNames)]...),
		setShale.Row(backButtonShale),
	)

	miniHouseNames := getHouseNames(db, miniHouseName)
	for _, v := range miniHouseNames {
		miniHouseBtns = append(miniHouseBtns, setMiniHouse.Text(v))
	}

	backButtonMiniHouse := setMiniHouse.Text(backButton)
	miniHouseBtns = append(miniHouseBtns, backButtonMiniHouse)

	setMiniHouse.Reply(
		setMiniHouse.Row(miniHouseBtns[:len(miniHouseNames)]...),
		setMiniHouse.Row(backButtonMiniHouse),
	)

	return
}

func HouseNumber(c tb.Context, mainMenu *tb.ReplyMarkup, db *sql.DB) error {
	chatID := c.Sender().ID
	apartment := c.Text()

	if err := Handlers.AddOrUpdateUserApartment(db, chatID, apartment); err != nil {
		fmt.Println(err)
		return c.Send("Произошла ошибка при сохранении номера квартиры.")
	}

	return c.Send(fmt.Sprintf("Номер сохранен: %s. Теперь вы можете использовать другие команды.", apartment), mainMenu)
}

func SetupSetHouseHandlers(b *tb.Bot, mainMenu *tb.ReplyMarkup, db *sql.DB) *tb.ReplyMarkup {
	setHouseBtns,
		shaleBtns,
		miniHouseBtns,
		setHouse,
		setShale,
		setMiniHouse := createSetHouseBtns(db)

	b.Handle(&setHouseBtns[0], func(c tb.Context) error {
		return c.Send(chooseShale, setShale)
	})

	b.Handle(&setHouseBtns[1], func(c tb.Context) error {
		return c.Send(chooseMiniHouse, setMiniHouse)
	})

	for _, shaleBtn := range shaleBtns {
		b.Handle(&shaleBtn, func(c tb.Context) error {
			if c.Text() == backButton {
				return c.Send("Пожалуйста, выберите свой дом", setHouse)
			}
			return HouseNumber(c, mainMenu, db)
		})
	}

	for _, miniHouseBtn := range miniHouseBtns {
		b.Handle(&miniHouseBtn, func(c tb.Context) error {
			if c.Text() == backButton {
				return c.Send("Пожалуйста, выберите свой дом", setHouse)
			}
			return HouseNumber(c, mainMenu, db)
		})
	}

	return setHouse
}
