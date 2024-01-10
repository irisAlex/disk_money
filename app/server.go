package app

import (
	"fmt"
	"money/pkg/log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine
}

func InitServer() *Server {

	return &Server{
		engine: NewGinEngine(),
	}
}

func (s *Server) Start() {
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8088)
	serv := &http.Server{
		Addr:         addr,
		Handler:      s.engine,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//go func() {
	log.Infof("http server start, listen addr: [%s]", addr)
	var err error

	err = serv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Error(err.Error())
	}
}
