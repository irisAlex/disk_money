package api

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	// form, err := c.MultipartForm()
	// if err != nil {
	//    c.String(http.StatusBadRequest, fmt.Sprintf("解析请求体失败: %s", err.Error()))
	//    return
	// }
	// // 处理文件
	// files := form.File["files"]
	// for _, file := range files {
	//    err = c.SaveUploadedFile(file, "<你的文件存储路径>")
	//    if err != nil {
	// 	  c.String(http.StatusBadRequest, fmt.Sprintf("上传文件失败: %s", err.Error()))
	// 	  return
	//    }
	// }

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	basePath := "./upload/"
	filename := basePath + filepath.Base(file.Filename)
	if err := c.SaveUploadedFile(file, filename); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("文件 %s 上传成功 ", file.Filename))
}
