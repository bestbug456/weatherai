package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func consumeTelegramMessage(higher int, telegramKey string) error {
	tgRequest := TelegramUpdateRequest{
		Offset: higher + 1,
	}
	encoded, err := json.Marshal(tgRequest)
	if err != nil {
		return err
	}
	_, err = http.Post("https://api.telegram.org/bot"+telegramKey+"/getUpdates", "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return err
	}
	return nil
}

func createTelegramMessageAndSend(weather Response, telegramKey, chatid string) error {

	var tgMsg TelegramMessage
	switch weather.Weather[0].Main {
	case "Clear":
		tgMsg = TelegramMessage{
			ChatID: chatid,
			Text:   "Today is a good day!",
		}
		break
	case "Clouds":
		tgMsg = TelegramMessage{
			ChatID: chatid,
			Text:   "Today is a good day, just some cloud in the air!",
		}
		break
	case "Drizzle":
		tgMsg = TelegramMessage{
			ChatID: chatid,
			Text:   "Today is a drizzle day! Close your window and bring with you a umbrella",
		}
		break
	case "Rain":
		tgMsg = TelegramMessage{
			ChatID: chatid,
			Text:   "Today is a rainy day! Close your window and bring with you a umbrella, and a kway!",
		}
		break
	case "Thunderstorm":
		tgMsg = TelegramMessage{
			ChatID: chatid,
			Text:   "Today is a thunderstorm day! Close your window and bring with you a umbrella and possibly don't stay outside!",
		}
		break
	case "Snow":
		tgMsg = TelegramMessage{
			ChatID: chatid,
			Text:   "Today is a snow day! Do you wanna build a snow man?",
		}
		break
	case "Extreme":
		tgMsg = TelegramMessage{
			ChatID: chatid,
			Text:   "Today is a day tagged as \"extreme\" day! Don't go outside home if possible.",
		}
		break
	}
	tgMsg.Text += fmt.Sprintln("More information: Max Temperature ", weather.Generic.Max, ", Min Temperature ", weather.Generic.Min, ", Humidity: ", weather.Generic.Humidity, ", Cloudiness: ", weather.Clouds.All, "%, Wind speed: ", weather.Wind.Speed)
	sendMessageOnTelegram(tgMsg, telegramKey)
	return nil
}

func sendMessageOnTelegram(msg TelegramMessage, telegramKey string) error {
	encoded, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Error while encoding the message: %s", err.Error())
	}
	_, err = http.Post("https://api.telegram.org/bot"+telegramKey+"/sendMessage", "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return fmt.Errorf("Error while sending the message: %s", err.Error())
	}
	return nil
}
