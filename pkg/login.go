package tripartite

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type login struct {
	User string
	Pwd  string
}

func Login(c *gin.Context) {

	content, err := New().SetHeaders().SetBody("action=get_vip_fl&file_id=4564078").Post("http://www.xunniuwp.com/ajax.php")

	fmt.Print(err)

	if err != nil {
		fmt.Println(content.String())
	}
	fmt.Println(content.String())
}

func Register(c *gin.Context) {
	c.HTML(200, "login.html", "register")
}

func Vip(c *gin.Context) {
	c.HTML(200, "vip.html", "vip")
}

func Disk(c *gin.Context) {
	c.HTML(200, "disk.html", "disk")
}
