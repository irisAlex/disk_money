package engine

import (
	"context"
	"fmt"
	"log"
	tripartite "money/pkg"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func StartServer(ctx context.Context) {
	addr := fmt.Sprintf("%s:%d", "172.17.36.168", 8088)
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

	v1.POST("register", tripartite.Register)
	v1.POST("verify", tripartite.Verify)
	v1.POST("login", tripartite.Login)
	return nil
}

func newGinEngine() *gin.Engine {

	app := gin.New()
	app.Use(Cors())
	// add gzip mw
	// http client of prometheus don't decode gzip content, then curl and chrome can decode it.

	// register fgprof
	// go tool pprof --http=:6061 http://localhost:6060/debug/fgprof?seconds=10

	// regisger custom route
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
