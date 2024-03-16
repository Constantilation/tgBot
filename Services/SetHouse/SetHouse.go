package SetHouse

import (
	"database/sql"
	tb "gopkg.in/telebot.v3"
	"tgBotElgora/DB/Handlers"
)

func createSetHouseBtns(db *sql.DB) (setHouseBtns []tb.Btn, shaleInlineKeyboard, miniHouseInlineKeyboard HouseInlineKeyboard, setHouse *tb.ReplyMarkup) {
	setHouse = &tb.ReplyMarkup{ResizeKeyboard: true}
	setHouseBtns = []tb.Btn{
		setHouse.Text(shaleName),     // 0
		setHouse.Text(miniHouseName), // 1
	}
	setHouse.Reply(
		setHouse.Row(setHouseBtns[0]),
		setHouse.Row(setHouseBtns[1]),
	)

	shaleNames := getHouseNames(db, shaleBDName)
	for _, name := range shaleNames {
		shaleInlineKeyboard.InlineKeyboard = append(shaleInlineKeyboard.InlineKeyboard, []tb.InlineButton{{Unique: name.Unique, Text: name.Text, Data: name.Data}})
	}

	miniHouseNames := getHouseNames(db, miniHouseBDName)
	for _, name := range miniHouseNames {
		miniHouseInlineKeyboard.InlineKeyboard = append(miniHouseInlineKeyboard.InlineKeyboard, []tb.InlineButton{{Unique: name.Unique, Text: name.Text, Data: name.Data}})
	}

	return
}

func HouseNumber(c tb.Context, house string, phoneNumber string, db *sql.DB) error {
	chatID := c.Chat().ID
	if err := Handlers.AddOrUpdateUserApartment(db, chatID, house, phoneNumber); err != nil {
		return err
	}

	return nil
}

func SetupSetHouseHandlers(b *tb.Bot, mainMenu *tb.ReplyMarkup, db *sql.DB) *tb.ReplyMarkup {
	setHouseBtns, shaleInlineKeyboard, miniHouseInlineKeyboard, setHouse := createSetHouseBtns(db)

	b.Handle(&setHouseBtns[0], func(c tb.Context) error {
		return c.Send(chooseShale, &tb.ReplyMarkup{
			InlineKeyboard: shaleInlineKeyboard.InlineKeyboard,
		})
	})

	b.Handle(&setHouseBtns[1], func(c tb.Context) error {
		return c.Send(chooseMiniHouse, &tb.ReplyMarkup{
			InlineKeyboard: miniHouseInlineKeyboard.InlineKeyboard,
		})
	})

	return setHouse
}
