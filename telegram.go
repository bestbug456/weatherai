package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/goml/gobrain"
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

type TelegramManager struct {
	ApyKey string
	Brain  *gobrain.FeedForward
}

func NewTelegramManager(telegramkey string, brain *gobrain.FeedForward) *TelegramManager {
	return &TelegramManager{
		ApyKey: telegramkey,
		Brain:  brain,
	}
}

func (tm *TelegramManager) CreateTelegramMessageAndSend(weather *Response, chatid string) error {

	var etiquetteResult string
	if tm.Brain != nil {
		etiquetteResult = tm.prediction(weather)
	} else {
		etiquetteResult = tm.standartEtiquette(weather)
	}

	tgMsg := TelegramMessage{
		Text:   etiquetteResult,
		ChatID: chatid,
	}

	tgMsg.Text += fmt.Sprintln("More information: Max Temperature ", weather.Generic.Max, ", Min Temperature ", weather.Generic.Min, ", Humidity: ", weather.Generic.Humidity, ", Cloudiness: ", weather.Clouds.All, "%, Wind speed: ", weather.Wind.Speed)
	tm.sendMessageOnTelegram(tgMsg)
	return nil
}

func (tm *TelegramManager) prediction(weather *Response) string {
	// Prepare the various data
	input := []float64{weather.Generic.Temperature, weather.Generic.Max, weather.Generic.Min, float64(weather.Clouds.All), float64(weather.Generic.Humidity)}
	// Execute the prediction
	ris, err := tm.Brain.Update(input)
	if err != nil {
		return tm.standartEtiquette(weather)
	}
	weather.Etiquetteweather = fmt.Sprintf("%v", ris)
	if ris[0] > 0.5 {
		return "Our artificial inteligence say it will be a \"bad\" day!"
	}
	if ris[1] > 0.5 {
		return "Our artificial inteligence say it will be a \"Good\" day!"
	}
	return tm.standartEtiquette(weather)
}

func (tm *TelegramManager) standartEtiquette(weather *Response) string {
	switch weather.Weather[0].Main {
	case "Clear":
		return "Today is a good day!"
	case "Clouds":
		return "Today is a good day, just some cloud in the air!"
	case "Drizzle":
		return "Today is a drizzle day! Close your window and bring with you a umbrella"
	case "Rain":
		return "Today is a rainy day! Close your window and bring with you a umbrella, and a kway!"
	case "Thunderstorm":
		return "Today is a thunderstorm day! Close your window and bring with you a umbrella and possibly don't stay outside!"
	case "Snow":
		return "Today is a snow day! Do you wanna build a snow man?"
	case "Extreme":
		return "Today is a day tagged as \"extreme\" day! Don't go outside home if possible."
	}
	return "WTF is this weather?!"
}
func (tm *TelegramManager) sendMessageOnTelegram(msg TelegramMessage) error {
	encoded, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("Error while encoding the message: %s", err.Error())
	}
	_, err = http.Post("https://api.telegram.org/bot"+tm.ApyKey+"/sendMessage", "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		return fmt.Errorf("Error while sending the message: %s", err.Error())
	}
	return nil
}
