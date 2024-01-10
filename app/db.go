package app

import "money/pkg/mongodb"

var Modeler mongodb.Modeler // global

func InitDB() {
	Modeler = mongodb.NewMongodb()
}
