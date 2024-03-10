package AdminPanel

import (
	"fmt"
	tb "gopkg.in/telebot.v3"
	"strconv"
)

func requestAdminApproval(c tb.Context, b *tb.Bot) (acceptBtn, rejectBtn tb.InlineButton) {
	adminID := GetAdminID() // Получаем ID админа из .env файла или конфигурации
	msg := fmt.Sprintf("Пользователь %s (%d) хочет получить доступ к боту.", c.Sender().FirstName, c.Sender().ID)

	// Создаем inline кнопки для одобрения/отклонения
	acceptBtn = tb.InlineButton{
		Unique: "accept",
		Text:   "Да",
		Data:   strconv.FormatInt(c.Sender().ID, 10), // Используем ID пользователя как данные кнопки
	}
	rejectBtn = tb.InlineButton{
		Unique: "reject",
		Text:   "Нет",
		Data:   strconv.FormatInt(c.Sender().ID, 10),
	}
	inlineKeys := [][]tb.InlineButton{{acceptBtn, rejectBtn}}

	// Отправляем сообщение админу с inline кнопками
	_, err := b.Send(tb.ChatID(adminID), msg, &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
	if err != nil {
		return tb.InlineButton{}, tb.InlineButton{}
	}

	c.Send("Дождитесь одобрения администратора")

	return
}

func CreateAdminApprovalHandlers(c tb.Context, b *tb.Bot, setHouse *tb.ReplyMarkup) {
	acceptBtn, rejectBtn := requestAdminApproval(c, b)

	b.Handle(&acceptBtn, func(c tb.Context) error {
		return c.Send("Пожалуйста, выберите свой дом", setHouse)
	})

	b.Handle(&rejectBtn, func(c tb.Context) error {
		return c.Send("В доступе отказано")
	})
}
