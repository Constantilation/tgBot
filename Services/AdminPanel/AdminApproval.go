package AdminPanel

import (
	"fmt"
	"log"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

func requestAdminApproval(c tb.Context, b *tb.Bot) (acceptBtn, rejectBtn tb.InlineButton) {
	adminID := GetAdminID()
	msg := fmt.Sprintf("Пользователь %s (%d) хочет получить доступ к боту.", c.Sender().FirstName, c.Sender().ID)

	acceptBtn = tb.InlineButton{
		Unique: "accept",
		Text:   "Да",
		Data:   strconv.FormatInt(c.Sender().ID, 10),
	}
	rejectBtn = tb.InlineButton{
		Unique: "reject",
		Text:   "Нет",
		Data:   strconv.FormatInt(c.Sender().ID, 10),
	}
	inlineKeys := [][]tb.InlineButton{{acceptBtn, rejectBtn}}

	_, err := b.Send(tb.ChatID(adminID), msg, &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
	if err != nil {
		return tb.InlineButton{}, tb.InlineButton{}
	}

	c.Send("Дождитесь одобрения администратора")

	return
}

func CreateAdminApprovalHandlers(c tb.Context, b *tb.Bot, setHouse *tb.ReplyMarkup) {

	livingRulesHandlerText := "📌Заезд в 15ч, выезд до 12ч, если не оговорены другие условия\n" +
		"📌Пожалуйста, не шумите на улице и общей территории после 23ч\n" +
		"📌Курение в домах в том числе кальянов, электронных сигарет и прочих электронных гаджетов запрещено\n" +
		"📌При проживании больше 5-ти суток предоставляется бесплатная уборка дома на 4-ые сутки"
	acceptBtn, rejectBtn := requestAdminApproval(c, b)

	b.Handle(&acceptBtn, func(c tb.Context) error {
		userID, err := strconv.ParseInt(c.Data(), 10, 64)
		if err != nil {
			log.Println("Ошибка при преобразовании данных кнопки в ID пользователя:", err)
			return err
		}

		rulesMarkup := &tb.ReplyMarkup{ResizeKeyboard: true}
		acknowledgeBtn := rulesMarkup.Text("ОЗНАКОМЛЕН")
		rulesMarkup.Reply(rulesMarkup.Row(acknowledgeBtn))

		_, err = b.Send(tb.ChatID(userID), livingRulesHandlerText, rulesMarkup)
		if err != nil {
			log.Println("Ошибка при отправке правил проживания пользователю:", err)
			return err
		}
		return nil
	})

	b.Handle("ОЗНАКОМЛЕН", func(c tb.Context) error {
		return c.Send("Добро пожаловать! Теперь вы можете пользоваться ботом.", setHouse)
	})

	b.Handle(&rejectBtn, func(c tb.Context) error {
		userID, err := strconv.ParseInt(c.Data(), 10, 64)
		if err != nil {
			log.Println("Ошибка при преобразовании данных кнопки в ID пользователя:", err)
			return err
		}
		_, err = b.Send(tb.ChatID(userID), "В доступе отказано.")
		if err != nil {
			log.Println("Ошибка при отправке сообщения пользователю:", err)
			return err
		}
		return nil
	})
}
