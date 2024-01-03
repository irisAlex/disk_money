package tripartite

// HandleDownloadFile 下载文件

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func HandleDownloadFile(c *gin.Context) {

	content, err := New().SetHeaders().SetBody("action=get_vip_fl&file_id=4564078").Post("http://www.xunniuwp.com/ajax.php")

	fmt.Print(err)

	if err != nil {
		fmt.Println(content.String())
	}
	fmt.Println(content.String())
}
