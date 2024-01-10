package dao

import (
	"errors"
	"money/model"
	"money/pkg/mongodb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Modeler = mongodb.NewMongodb()
	UserTab = model.TableName()
)

func VerifyUser(u, p string) (*model.LoginUserInfo, error) {
	var (
		user   = new(model.LoginUserInfo)
		filter = bson.M{"user": u, "pwd": p}
	)
	err := Modeler.FindOne(UserTab, filter, bson.M{}, user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFoundRecord")
	}
	if err != nil {
		return nil, err
	}

	return user, nil

}

func GetUserInfo(name, email string) (*model.RegisterUserInfo, error) {
	var (
		svc    = new(model.RegisterUserInfo)
		filter = bson.M{"user": name, "email": email}
	)

	err := Modeler.FindOne(UserTab, filter, bson.M{}, svc)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFoundRecord")
	}
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func InsertUser(u *model.RegisterUserInfo) error {
	if err := Modeler.InsertOne(UserTab, u); err != nil {

		return err
	}
	return nil
}
