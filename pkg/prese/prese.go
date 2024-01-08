package prese

import (
	"money/pkg/aes"

	"github.com/gin-gonic/gin"
)

const (
	prefix = "disk"
	// ResBodyKey 存储上下文中的键(响应Body数据)
	ResBodyKey = prefix + "/res-body"
)

func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
}

func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := aes.JSONMarshal(v)
	if err != nil {
		panic(err)
	}
	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
}
