package api

import (
	"money/dao"
	"money/model"
	"money/pkg/aes"
	"money/pkg/prese"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var account = &model.AccountInfo{}
	prese.ParseJSON(c, &account)

	if account.Name == "" || account.Cipher == "" || account.Email == "" {
		prese.ResJSON(c, 400, "注册信息不能为空")
		return
	}

	ip := c.Request.Header.Get("X-Forwarded-For")

	// 如果 X-Forwarded-For 为空，则检查 X-Real-IP 头部
	if ip == "" {
		ip = c.Request.Header.Get("X-Real-IP")
	}

	// 如果仍然为空，则从请求中获取 RemoteAddr
	if ip == "" {
		ip = c.Request.RemoteAddr
	}

	if ip == "" {
		prese.ResJSON(c, 400, "ip 非法")
		return
	}

	u, err := dao.GetUserAccount(account.Name)

	if err == nil && u.Name != "" {
		prese.ResJSON(c, 400, "用户已经存在")
		return
	}

	remail, err := dao.GetEmail(account.Email)

	if err == nil && remail.Email == account.Email {
		prese.ResJSON(c, 400, "email已被注册")
		return
	}

	account.CreateTime = time.Now()
	account.RealIP = ip

	// 将哈希值转换为十六进制字符串
	aes.Md5(&account.Cipher)

	dao.InsertUser(account)

	prese.ResJSON(c, 200, resData{
		Token:    aes.GetToken(account.Name),
		User:     account.Name,
		Email:    account.Email,
		SetMeal:  0,
		ExriTime: 0,
	})
}

type resData struct {
	Token    string `json:"token"`
	User     string `json:"user"`
	Email    string `json:"email"`
	SetMeal  int    `json:"set_meal"`
	ExriTime int64  `json:"exri_time"`
}

type resDownData struct {
	Ct string `json:"content"`
}

// type email struct {
// 	Email
// }

func Login(c *gin.Context) {
	lg := &model.LoginUserInfo{}
	prese.ParseJSON(c, lg)
	if lg.Name == "" || lg.Cipher == "" {
		prese.ResJSON(c, 400, "用户名或密码错误不能为空")
		return
	}

	aes.Md5(&lg.Cipher)

	account, err := dao.VerifyUser(lg.Name, lg.Cipher)
	if account == nil || err != nil {
		prese.ResJSON(c, 400, "用户名或密码错误")
		return
	}

	prese.ResJSON(c, 200, resData{
		User:     lg.Name,
		Email:    account.Email,
		Token:    aes.GetToken(lg.Name),
		ExriTime: account.ExpiredTime,
	})
}

func Verify(c *gin.Context) {
	var (
		v       = model.VerifyUserInfo{}
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

	if info[0] != v.Name && effTime < time.Now().Unix() {
		prese.ResJSON(c, 400, "token 非法")
		return
	}

	prese.ResJSON(c, 200, "success")
}
