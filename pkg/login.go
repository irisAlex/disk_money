package tripartite

import (
	"money/pkg/mongodb"
	"strconv"
	"strings"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	prefix = "disk"
	// ResBodyKey 存储上下文中的键(响应Body数据)
	ResBodyKey = prefix + "/res-body"
)

var modeler = mongodb.NewMongodb() // global
var UserTab = "user"

type register struct {
	User  string `json:"user" bson:"user" valid:"-"`
	Pwd   string `json:"pwd" bson:"pwd" valid:"-"`
	Email string `json:"email" bson:"email" valid:"-"`
}

type resData struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

// type email struct {
// 	Email
// }

func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return err
	}
	return nil
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

type login struct {
	User string `json:"user"`
	Pwd  string `json: "pwd"`
}

func Login(c *gin.Context) {
	lg := &login{}
	ParseJSON(c, lg)
	if lg.User == "" || lg.Pwd == "" {
		ResJSON(c, 400, "用户名或密码错误不能为空")
		return
	}

	userI, err := verifyUserInfo(lg.User, lg.Pwd)
	if userI == nil || err != nil {
		ResJSON(c, 400, "用户名或密码错误")
		return
	}

	ResJSON(c, 200, resData{
		User:  lg.User,
		Token: getToken(lg.User),
	})
}

type verify struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

func Verify(c *gin.Context) {
	var (
		v       = verify{}
		effTime int64
		err     error
	)
	ParseJSON(c, &v)

	if v.Token == "" {
		ResJSON(c, 500, "非法请求")
		return
	}

	verfiyInfo, _ := GcmDecrypt(v.Token)

	info := strings.Split(verfiyInfo, "|")

	if len(info) < 1 {
		ResJSON(c, 500, "非法请求")
		return
	}

	effTime, err = strconv.ParseInt(info[1], 10, 64)
	if err != nil {
		ResJSON(c, 500, "非法请求")
		return
	}

	if info[0] != v.User && effTime < time.Now().Unix() {
		ResJSON(c, 400, "token 非法")
		return
	}

	ResJSON(c, 200, "success")
}

func Register(c *gin.Context) {
	var ri = &register{}
	ParseJSON(c, ri)
	if ri.User == "" || ri.Pwd == "" || ri.Email == "" {
		ResJSON(c, 400, "注册信息不能为空")
		return
	}
	svc, err := getUserInfo(ri.User, ri.Email)
	if svc != nil && err == nil {
		ResJSON(c, 400, "用户已经存在或email已被注册")
		return
	}
	if modeler.InsertOne(UserTab, ri) != nil {
		ResJSON(c, 500, "客户端错误")
		return
	}
	ResJSON(c, 200, resData{
		Token: getToken(ri.User),
		User:  ri.User,
	})
}

func verifyUserInfo(u, p string) (*login, error) {
	var (
		user   = new(login)
		filter = bson.M{"user": u, "pwd": p}
	)
	err := modeler.FindOne(UserTab, filter, bson.M{}, user)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFoundRecord")
	}
	if err != nil {
		return nil, err
	}

	return user, nil

}

func getUserInfo(name, email string) (*register, error) {
	var (
		svc    = new(register)
		filter = bson.M{"user": name, "email": email}
	)

	err := modeler.FindOne(UserTab, filter, bson.M{}, svc)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("ErrNotFoundRecord")
	}
	if err != nil {
		return nil, err
	}

	return svc, nil
}
