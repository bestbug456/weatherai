package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/mgo.v2"
)

type RequestManager struct {
	WeatherApiKey string
	Chatid        string
	Telegramkey   string
	s             *mgo.Session
}

func InitInlineMode(weatherApiKey, chatid, telegramkey, address, username, password, option string, enableSSL bool, pvkeypath, chainpath string) {
	manager := NewRequestManager(weatherApiKey, chatid, telegramkey, address, username, password, option, enableSSL)
	// Enable 'secret' handler
	// only you and telegram
	// should know the telegram key.
	http.HandleFunc("/"+telegramkey, manager.NewRequest)
	http.HandleFunc("/ping", manager.Ping)
	if pvkeypath != "" && chainpath != "" {
		err := http.ListenAndServeTLS(":443", pvkeypath, chainpath, nil)
		if err != nil {
			fmt.Printf("HTTPS unable to listen and serve on address: %s cause error: %s.\n", "0.0.0.0:443", err.Error())
			os.Exit(1)
		}
	}
	log.Printf("running in insecure mode (chain is: %s and key is %s)", chainpath, pvkeypath)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("HTTPS unable to listen and serve on address: %s cause error: %s.\n", "0.0.0.0:8080", err.Error())
		os.Exit(1)
	}
}

func NewRequestManager(weatherApiKey, chatid, telegramKey, address, username, password, option string, enableSSL bool) *RequestManager {
	s, err := NewDbSession(address, username, password, option, enableSSL)
	if err != nil {
		sendMessageError(chatid, telegramKey, err)
		s = nil
	}
	return &RequestManager{
		WeatherApiKey: weatherApiKey,
		Chatid:        chatid,
		Telegramkey:   telegramKey,
		s:             s,
	}
}

func (rm *RequestManager) Ping(w http.ResponseWriter, r *http.Request) {
	/* return value */
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("pong")
}

func (rm *RequestManager) NewRequest(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request TelegramMessageResponse
	err := decoder.Decode(&request)
	if err != nil {
		sendMessageError(rm.Chatid, rm.Telegramkey, err)
		return
	}
	defer r.Body.Close()
	if request.Message.Chat.Id == 0 {
		sendMessageError(rm.Chatid, rm.Telegramkey, fmt.Errorf("receive message without chat id."))
		return
	}
	if request.Message.GPS.Lat == 0.0 && request.Message.GPS.Long == 0.0 {
		chatid := strconv.FormatInt(request.Message.Chat.Id, 10)
		msg := TelegramMessage{
			ChatID: chatid,
			Text:   "Send to me your position and I will tell to you the weather :)",
		}
		err = sendMessageOnTelegram(msg, rm.Telegramkey)
		if err != nil {
			sendMessageError(rm.Chatid, rm.Telegramkey, err)
		}
		return
	}
	// Get the weather
	url := fmt.Sprintf("%s%v%s%v%s%s", "https://api.openweathermap.org/data/2.5/weather?lat=", request.Message.GPS.Lat, "&lon=", request.Message.GPS.Long, "&units=metric&appid=", rm.WeatherApiKey)
	resp, err := http.Get(url)
	if err != nil {
		sendMessageError(rm.Chatid, rm.Telegramkey, err)
		return
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.Bytes()
	var weather Response
	err = json.Unmarshal(body, &weather)
	if err != nil {
		sendMessageError(rm.Chatid, rm.Telegramkey, err)
		return
	}
	// Dump information to the database
	// if is enabled
	if rm.s != nil {
		s := rm.s.Copy()
		err = dumpWeatherToDatabaseInline(weather, s)
		s.Close()
		if err != nil {
			sendMessageError(rm.Chatid, rm.Telegramkey, err)
		}
	}
	chatid := strconv.FormatInt(request.Message.Chat.Id, 10)
	err = createTelegramMessageAndSend(weather, rm.Telegramkey, chatid)
	if err != nil {
		sendMessageError(rm.Chatid, rm.Telegramkey, err)
		return
	}

}
