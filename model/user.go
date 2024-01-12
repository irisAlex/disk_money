package model

import "time"

type VerifyUserInfo struct {
	Token string `json:"token"`
	Name  string `json:"user"`
}

type LoginUserInfo struct {
	Name   string `json:"user"`
	Cipher string `json:"cipher"`
}

type DownTime struct {
	User   string `json:"user"`
	FileId string `json:"file_id"`
}

type AccountInfo struct {
	Name        string    `json:"name" bson:"user" valid:"-"`
	Cipher      string    `json:"cipher" bson:"cipher" valid:"-"`
	Email       string    `json:"email" bson:"email" valid:"-"`
	RealIP      string    `json:"real_ip" bson:"real_ip"`
	ExpiredTime int64     `json:"expired_time" bson:"expired_time"`
	CreateTime  time.Time `json:"create_time" bson:"create_time"`
	UpdateTime  int64     `json:"update_time" ,bson:"update_time"`
	DayDownTime int       `json:"down_time" bson:"down_time"`
	Vip         int       `json:"vip" bson:"vip"`
	SetMeal     int       `json:"set_meal" bson:"set_meal"`
}

type Vip struct {
	User    string `json:"user"`
	SetMeal int    `json:"set_meal"`
}

func TableName() string {
	return "account"
}
