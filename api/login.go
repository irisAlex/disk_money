package api

import (
	"money/dao"
	"money/pkg/aes"
	"money/pkg/prese"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var ri = &dao.RegisterUserInfo{}
	prese.ParseJSON(c, ri)
	if ri.User == "" || ri.Pwd == "" || ri.Email == "" {
		prese.ResJSON(c, 400, "注册信息不能为空")
		return
	}
	svc, err := dao.GetUserInfo(ri.User, ri.Email)
	if svc != nil && err == nil {
		prese.ResJSON(c, 400, "用户已经存在或email已被注册")
		return
	}

	dao.InsertUser(ri)

	prese.ResJSON(c, 200, resData{
		Token: aes.GetToken(ri.User),
		User:  ri.User,
	})
}

type resData struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

// type email struct {
// 	Email
// }

func Login(c *gin.Context) {
	lg := &dao.LoginUserInfo{}
	prese.ParseJSON(c, lg)
	if lg.User == "" || lg.Pwd == "" {
		prese.ResJSON(c, 400, "用户名或密码错误不能为空")
		return
	}

	userI, err := dao.VerifyUser(lg.User, lg.Pwd)
	if userI == nil || err != nil {
		prese.ResJSON(c, 400, "用户名或密码错误")
		return
	}

	prese.ResJSON(c, 200, resData{
		User:  lg.User,
		Token: aes.GetToken(lg.User),
	})
}

func Verify(c *gin.Context) {
	var (
		v       = dao.VerifyUserInfo{}
		effTime int64
		err     error
	)
	prese.ParseJSON(c, &v)

	if v.Token == "" {
		prese.ResJSON(c, 500, "非法请求")
		return
	}

	verfiyInfo, _ := aes.GcmDecrypt(v.Token)

	info := strings.Split(verfiyInfo, "|")

	if len(info) < 1 {
		prese.ResJSON(c, 500, "非法请求")
		return
	}

	effTime, err = strconv.ParseInt(info[1], 10, 64)
	if err != nil {
		prese.ResJSON(c, 500, "非法请求")
		return
	}

	if info[0] != v.User && effTime < time.Now().Unix() {
		prese.ResJSON(c, 400, "token 非法")
		return
	}

	prese.ResJSON(c, 200, "success")
}
