package model

type RegisterUserInfo struct {
	User  string `json:"user" bson:"user" valid:"-"`
	Pwd   string `json:"pwd" bson:"pwd" valid:"-"`
	Email string `json:"email" bson:"email" valid:"-"`
}

type VerifyUserInfo struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

type LoginUserInfo struct {
	User string `json:"user"`
	Pwd  string `json: "pwd"`
}

func TableName() string {
	return "account"
}
