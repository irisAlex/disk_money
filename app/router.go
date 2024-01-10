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
	v1.POST("fileUpload", api.UploadFile)
	v1.GET("downloadFile", api.DownloadFile)
	return nil
}
