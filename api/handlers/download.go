package api

// HandleDownloadFile 下载文件

import (
	"money/dao"
	"money/pkg/aes"
	"money/pkg/client"
	"money/pkg/prese"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	D = 0
	M = 1
	Q = 2
	Y = 3
	P = -1
)

func DownloadFile(c *gin.Context) {
	file_id_x := c.Query("file_id")

	Token := c.Request.Header.Get("token")
	if Token == "" {
		prese.ResJSON(c, 400, "非法请求")
		return
	}

	verfiyInfo, _ := aes.GcmDecrypt(Token)

	info := strings.Split(verfiyInfo, "|")

	if len(info) < 1 {
		prese.ResJSON(c, 400, "非法请求")
		return
	}

	if file_id_x == "" {
		prese.ResJSON(c, 500, "文件Id 不能为空")
		return
	}

	account, err := dao.GetUserAccount(info[0])

	if err != nil {
		prese.ResJSON(c, 400, "非法请求")
		return
	}

	context := Context{}
	switch account.SetMeal {
	case D:
		context.SetDiscountStrategy(RegularCustomerDiscount{})
	case M:
		context.SetDiscountStrategy(MonthPremiumCustomerDiscount{})

	case Q:
		context.SetDiscountStrategy(QuaretrPremiumCustomerDiscount{})
	case Y:
		context.SetDiscountStrategy(YearPremiumCustomerDiscount{})
	case P:
		context.SetDiscountStrategy(PermanentPremiumCustomerDiscount{})
	}

	if !context.IsPermission(account.DayDownTime) {
		prese.ResJSON(c, 400, "以达最大文件下载次数")
		return
	}

	content, erro := getDownFileLink(file_id_x)
	if erro != nil {
		prese.ResJSON(c, 400, "以达最大文件下载次数")
		return
	}

	account.DayDownTime += 1

	if dao.UpdateAccount(account) != nil {
		prese.ResJSON(c, 500, "服务异常")
		return
	}
	prese.ResJSON(c, 200, resDownData{
		Ct: content,
	})
}

func getDownFileLink(fix_id string) (string, error) {
	content, err := client.New().SetHeaders().SetBody("action=get_vip_fl&file_id=" + fix_id).Post("http://www.xunniuwp.com/ajax.php")

	if err != nil {
		return "", err
	}
	return content.String(), nil
}
