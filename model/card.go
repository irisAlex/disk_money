package model

type Card struct {
	CardSecret string `json:"card_secret" bson:"card_secret"`
	SetMeal    int    `json:"set_meal" bson:"set_meal"`
	Available  bool   `json:"available" bson:"available"`
}

func (c Card) TableName() string {
	return "card"
}
