package tripartite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"money/pkg/mongodb"

	"github.com/gin-gonic/gin"
)

const (
	prefix = "workflow"
	// UserIDKey 存储上下文中的键(用户ID)
	UserIDKey = prefix + "/user-id"
	// TraceIDKey 存储上下文中的键(跟踪ID)
	TraceIDKey = prefix + "/trace-id"
	// ResBodyKey 存储上下文中的键(响应Body数据)
	ResBodyKey = prefix + "/res-body"
)

var modeler = mongodb.NewMongodb() // global

type register struct {
	User  string `json:"user" bson:"user" valid:"-"`
	Pwd   string `json:"pwd" bson:"pwd" valid:"-"`
	Email string `json:"email" bson:"email" valid:"-"`
}

func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
}

func JSONMarshal(obj interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	err := enc.Encode(obj)
	if err != nil {
		return nil, err
	}

	// json.NewEncoder.Encode adds a final '\n', json.Marshal does not.
	// Let's keep the default json.Marshal behaviour.
	res := b.Bytes()
	if len(res) >= 1 && res[len(res)-1] == '\n' {
		res = res[:len(res)-1]
	}
	return res, nil
}

func ResJSON(c *gin.Context, status int, v interface{}) {
	buf, err := JSONMarshal(v)
	if err != nil {
		panic(err)
	}
	c.Set(ResBodyKey, buf)
	c.Data(status, "application/json; charset=utf-8", buf)
	c.Abort()
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
	ParseJSON(c, s)

	err := modeler.InsertOne("register", s)
	if err != nil {
		log.Println(err)
		ResJSON(c, 400, "客户端错误")
	}

	ResJSON(c, 200, "注册成功")
}
