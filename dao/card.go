package dao

import (
	"errors"
	"money/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetCardInfoByKey(key string) (*model.Card, error) {
	var (
		card   = new(model.Card)
		filter = bson.M{"card_secret": key}
	)

	err := Modeler.FindOne(card.TableName(), filter, bson.M{}, card)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFoundRecord")
	}
	if err != nil {
		return nil, err
	}

	return card, nil
}

func UpdateCardInfo(card *model.Card) error {
	filter := bson.M{"card_secret": card.CardSecret}
	if err := Modeler.Upsert(card.TableName(), filter, card); err != nil {
		return err
	}
	return nil
}

func InsertCard(card *model.Card) error {
	if err := Modeler.InsertOne(card.TableName(), card); err != nil {
		return err
	}
	return nil
}
