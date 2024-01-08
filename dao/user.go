package dao

import (
	"errors"

	"money/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	modeler = mongodb.NewMongodb() // global
	UserTab = "user"
)

type RegisterUserInfo struct {
	User  string `json:"user" bson:"user" valid:"-"`
	Pwd   string `json:"pwd" bson:"pwd" valid:"-"`
	Email string `json:"email" bson:"email" valid:"-"`
}

type LoginUserInfo struct {
	User string `json:"user"`
	Pwd  string `json: "pwd"`
}

type VerifyUserInfo struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

func VerifyUser(u, p string) (*LoginUserInfo, error) {
	var (
		user   = new(LoginUserInfo)
		filter = bson.M{"user": u, "pwd": p}
	)
	err := modeler.FindOne(UserTab, filter, bson.M{}, user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFoundRecord")
	}
	if err != nil {
		return nil, err
	}

	return user, nil

}

func GetUserInfo(name, email string) (*RegisterUserInfo, error) {
	var (
		svc    = new(RegisterUserInfo)
		filter = bson.M{"user": name, "email": email}
	)

	err := modeler.FindOne(UserTab, filter, bson.M{}, svc)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFoundRecord")
	}
	if err != nil {
		return nil, err
	}

	return svc, nil
}

type user struct {
}

func InsertUser(u *RegisterUserInfo) error {
	if err := modeler.InsertOne(UserTab, u); err != nil {

		return err
	}
	return nil
}
