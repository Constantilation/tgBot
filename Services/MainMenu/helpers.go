package MainMenu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Используйте внешнюю конфигурацию для хранения токена
func sendContactViaTelegramAPI(botToken string, chatID int64, firstName, lastName, phoneNumber string) {
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendContact", botToken)
	requestBody, err := json.Marshal(map[string]interface{}{
		"chat_id":      chatID,
		"first_name":   firstName,
		"last_name":    lastName,
		"phone_number": phoneNumber,
	})
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Response from Telegram API:", string(body))
}
