package app

import (
	api "money/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(app *gin.Engine) error {
	var (
		v1 = app.Group("/")
	)

	v1.POST("register", api.Register)
	v1.POST("verify", api.Verify)
	v1.POST("login", api.Login)
	v1.POST("file_upload", api.UploadFile)
	v1.GET("download_file", api.DownloadFile)
	v1.POST("cash_vip", api.GrantVip)
	v1.GET("get_card_key", api.GenerateCardKey)

	v1.POST("chat", api.Chat)
	return nil
}
