package main

type Response struct {
	Coord      Coordinate         `json:"coord"`
	Weather    []WeatherInfo      `json:"weather"`
	Generic    GenericInformation `json:"main"`
	Visibility int                `json:"visibility"`
	Wind       WindInformation    `json:"wind"`
	Clouds     CloudInformation   `json:"clouds`
	DayTime    int64              `json:"dt`
	Data       string             `json:"-"` // Internal: use for identify the data of the query
}

type Coordinate struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

type WeatherInfo struct {
	Id          int    `json:"id"`
	Main        string `json:"main`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type GenericInformation struct {
	Temperature float64 `json:"temp"`
	Pressure    float64 `json:"pressure"`
	Humidity    int     `json:"humidity"`
	Max         float64 `json:"temp_max"`
	Min         float64 `json:"temp_min"`
}

type WindInformation struct {
	Speed     float64 `json:"speed"`
	Direction float64 `json:"deg"`
}

type CloudInformation struct {
	All int `json:"all"`
}

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type TelegramUpdateRequest struct {
	Offset int `json:"offset"`
}

type TelegramUpdateObjectResponse struct {
	Result TelegramMessageResponse `json:"result"`
	Status bool                    `json:"ok"`
}
type TelegramMessageResponse struct {
	Message  Message `json:"message"`
	UpdateID int     `json:"update_id"`
}

//message
type Message struct {
	Updateid int      `json:"update_id"`
	Message  string   `json:"text"`
	Chat     Chat     `json:"chat"`
	GPS      Location `json:"location"`
}

type Location struct {
	Long float64 `json:"longitude"`
	Lat  float64 `json:"latitude"`
}

type Chat struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}
