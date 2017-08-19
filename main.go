package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

func main() {
	telegramkey := flag.String("telegramkey", "", "the telegram bot api key")
	weatherkey := flag.String("weatherkey", "", "the openweather api key")
	mongoAddress := flag.String("dbadd", "127.0.0.1:27017", "the mongodb address")
	mongoUser := flag.String("dbusr", "", "the mongodb username")
	mongopwd := flag.String("dbpwd", "", "the mongodb password")
	mongooption := flag.String("dbopt", "", "the mongodb option (optional)")
	ssl := flag.Bool("ssl", false, "enable the mongossl connection")
	inlinemode := flag.Bool("inline", false, "enable the bot to responde via webhook")
	chatid := flag.String("chatid", "", "your telegram chat id")
	fullchain := flag.String("chain", "", "the path for the chain (optional)")
	pvkey := flag.String("pvkey", "", "the path of your private key (optional)")
	flag.Parse()

	if *inlinemode {
		InitInlineMode(*weatherkey, *chatid, *telegramkey, *mongoAddress, *mongoUser, *mongopwd, *mongooption, *ssl, *fullchain, *pvkey)
	}
	ticker := NewJobTicker()
	tgmanager := NewTelegramManager(*telegramkey, nil)
	for {
		<-ticker.t.C
		ticker.updateJobTicker()
		weather, err := getDailyReport(*weatherkey, *chatid, *telegramkey)
		if err != nil {
			sendMessageError(*chatid, *telegramkey, err)
			continue
		}
		if *mongoUser != "" && *mongopwd != "" {
			err = dumpWeatherToDatabase(weather, *mongoAddress, *mongoUser, *mongopwd, *mongooption)
			if err != nil {
				sendMessageError(*chatid, *telegramkey, fmt.Errorf("Error while interacting with the databse: %s", err.Error()))
			}
		}
		err = tgmanager.CreateTelegramMessageAndSend(&weather, *chatid)
		if err != nil {
			sendMessageError(*chatid, *telegramkey, err)
			continue
		}
	}

}

func getDailyReport(weatherApiKey string, chatid, telegramkey string) (Response, error) {
	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=" + "Milan" + "," + "it" + "&units=metric&appid=" + weatherApiKey)
	if err != nil {
		return Response{}, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.Bytes()
	var weather Response
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return Response{}, err
	}
	return weather, nil
}

func sendMessageError(chatid, tgkey string, err error) {
	msg := TelegramMessage{
		Text:   err.Error(),
		ChatID: chatid,
	}
	encoded, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error while encoding the message: %s", err.Error())
	}
	_, err = http.Post("https://api.telegram.org/bot"+tgkey+"/sendMessage", "application/json", bytes.NewBuffer(encoded))
	if err != nil {
		fmt.Printf("Error while sending the message: %s\n", err.Error())
	}
}
