package Helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	config "tgBotElgora/Config"
)

var houseTypeNames = map[string]string{
	"shale":     "Шале",
	"miniHouse": "Мини-дом",
}

func IsAdmin(userID int64) bool {
	adminID, err := config.GetInt64Config(config.AdminId)

	if err != nil {
		log.Fatal("ADMIN_ID is not set in .env")
	}

	return userID == adminID
}

func GetHouseName(fullName string) (string, error) {
	parts := strings.Split(fullName, "_")
	if len(parts) != 2 {
		return "", fmt.Errorf("неправильный формат ввода")
	}

	houseType := parts[0]
	houseNumber := parts[1]

	houseTypeName, ok := houseTypeNames[houseType]
	if !ok {
		return "", fmt.Errorf("неизвестный тип дома: %s", houseType)
	}

	return fmt.Sprintf("%s %s", houseTypeName, houseNumber), nil
}

func SendContactViaTelegramAPI(botToken string, chatID int64, firstName, lastName, phoneNumber string) {
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
