package engine

import (
	"context"
	"fmt"
	"log"

	"money/api"
	"money/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StartServer(ctx context.Context) {
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8088)
	serv := &http.Server{
		Addr:         addr,
		Handler:      newGinEngine(),
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//go func() {
	log.Printf("http server start, listen addr: [%s]", addr)
	var err error

	err = serv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Print(err.Error())
	}
	//}()
}

func RegisterRouter(app *gin.Engine) error {
	var (
		v1 = app.Group("/")
	)

	v1.POST("register", api.Register)
	v1.POST("verify", api.Verify)
	v1.POST("login", api.Login)
	return nil
}

func newGinEngine() *gin.Engine {

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
