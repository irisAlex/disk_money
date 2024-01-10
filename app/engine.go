package app

import (
	"money/middleware"

	"github.com/gin-gonic/gin"
)

func NewGinEngine() *gin.Engine {

	app := gin.New()
	app.Use(middleware.Cors())

	app.Use(middleware.LoggerMiddleware())

	err := RegisterRouter(app)
	handleError(err)

	// swagger
	// if dir := cfg.Http.Swagger; dir != "" {
	// 	app.Static("/swagger", dir)
	// }

	return app
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
