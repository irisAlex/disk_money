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
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", 8088)
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
	v1.GET("money", tripartite.Login)
	v1.GET("index", func(c *gin.Context) {
		c.HTML(200, "index.html", "nil")
		c.Abort()
	})
	v1.GET("account.php", tripartite.Register)
	v1.GET("mydisk.php", tripartite.Disk)
	return nil
}

func newGinEngine() *gin.Engine {

	app := gin.New()

	app.LoadHTMLGlob("html/*")
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
