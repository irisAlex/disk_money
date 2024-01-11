package model

import "time"

type RegisterUserInfo struct {
	User        string    `json:"user" bson:"user" valid:"-"`
	Pwd         string    `json:"pwd" bson:"pwd" valid:"-"`
	Email       string    `json:"email" bson:"email" valid:"-"`
	RealIP      string    `json:"real_ip" bson:"real_ip"`
	ExpiredTime time.Time `json:"expired_time" bson:"expired_time"`
	CreateTime  time.Time `json:"create_time" bson:"create_time"`
	DownTime    int       `json:"down_time" bson:"down_time"`
	Vip         int       `json:"vip" bson:"vip"`
	SetMeal     int       `json:"set_meal" bson:"set_meal"`
}

type VerifyUserInfo struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

type LoginUserInfo struct {
	User string `json:"user"`
	Pwd  string `json: "pwd"`
}

type AccountInfo struct {
	User        string    `json:"user" bson:"user" valid:"-"`
	Pwd         string    `json:"pwd" bson:"pwd" valid:"-"`
	Email       string    `json:"email" bson:"email" valid:"-"`
	RealIP      string    `json:"real_ip" bson:"real_ip"`
	ExpiredTime time.Time `json:"expired_time" bson:"expired_time"`
	CreateTime  time.Time `json:"create_time" bson:"create_time"`
	DownTime    int       `json:"down_time" bson:"down_time"`
	Vip         int       `json:"vip" bson:"vip"`
	SetMeal     int       `json:"set_meal" bson:"set_meal"`
}

func TableName() string {
	return "account"
}
