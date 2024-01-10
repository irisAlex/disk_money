package api

// HandleDownloadFile 下载文件

import (
	"money/pkg/client"
	"money/pkg/prese"

	"github.com/gin-gonic/gin"
)

func DownloadFile(c *gin.Context) {

	file_id_x := c.Query("file_id")

	if file_id_x == "" {
		prese.ResJSON(c, 500, "文件Id 不能为空")
		return
	}

	content, err := client.New().SetHeaders().SetBody("action=get_vip_fl&file_id=" + file_id_x).Post("http://www.xunniuwp.com/ajax.php")

	if err != nil {
		prese.ResJSON(c, 500, "下载文件失败")
		return
	}
	prese.ResJSON(c, 200, resDownData{
		Ct: content.String(),
	})
}
