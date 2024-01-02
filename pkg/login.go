package tripartite

import (
	"fmt"
	"log"
	"money/pkg/mongodb"

	"github.com/gin-gonic/gin"
)

var modeler = mongodb.NewMongodb() // global

type register struct {
	User  string `json:"user" bson:"user" valid:"-"`
	Pwd   string `json:"pwd" bson:"pwd" valid:"-"`
	Email string `json:"email" bson:"email" valid:"-"`
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
	var s = &register{}
	s.User = "libin"
	s.Pwd = "123456"
	s.Email = "127@qq.com"
	err := modeler.InsertOne("register", s)
	if err != nil {
		log.Println(err)
	}
}
