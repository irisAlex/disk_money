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
)

func VerifyUser(u, p string) (*model.AccountInfo, error) {
	var (
		user   = new(model.AccountInfo)
		filter = bson.M{"user": u, "cipher": p}
	)
	err := Modeler.FindOne(user.TableName(), filter, bson.M{}, user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFoundRecord")
	}
	if err != nil {
		return nil, err
	}

	return user, nil

}

func GetUserAccount(name string) (*model.AccountInfo, error) {
	var (
		svc    = new(model.AccountInfo)
		filter = bson.M{"user": name}
	)

	err := Modeler.FindOne(svc.TableName(), filter, bson.M{}, svc)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFoundRecord")
	}
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func GetEmail(email string) (*model.AccountInfo, error) {
	var (
		svc    = new(model.AccountInfo)
		filter = bson.M{"email": email}
	)

	err := Modeler.FindOne(svc.TableName(), filter, bson.M{}, svc)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFoundRecord")
	}
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func InsertUser(u *model.AccountInfo) error {
	if err := Modeler.InsertOne(u.TableName(), u); err != nil {

		return err
	}
	return nil
}

func UpdateAccount(account *model.AccountInfo) error {
	filter := bson.M{"user": account.Name}
	if err := Modeler.Upsert(account.TableName(), filter, account); err != nil {
		return err
	}
	return nil
}
