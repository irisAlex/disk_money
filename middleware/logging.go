package middleware

import (
	"bytes"
	"time"

	"money/pkg/log"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		log.Infof("http request method %s url %s", c.Request.Method, c.Request.URL)

		start := time.Now()
		writer := &CustomResponseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		statusCode := c.Copy().Writer.Status()
		if statusCode == 200 || statusCode == 201 {
			return
		}

		const rsize = 50
		resp := SubFrontString(writer.body.String(), rsize)
		log.Warnf("http response method %s url %s, code: %v, cost: %v, resp: %v",
			c.Request.Method, c.Request.URL, statusCode,
			time.Since(start).String(),
			resp,
		)
	}
}

func SubFrontString(source string, end int) string {
	var (
		r      = []rune(source)
		length = len(r)
	)

	if end > length {
		return source
	}
	return string(r[:end])
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
