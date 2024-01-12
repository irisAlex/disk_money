package api

import (
	"money/dao"
	"money/model"
	"money/pkg/prese"
	"time"

	"github.com/gin-gonic/gin"
)

func GrantVip(c *gin.Context) {
	vip := &model.Vip{}
	prese.ParseJSON(c, &vip)

	if vip.User == "" {
		prese.ResJSON(c, 400, "用户非法")
		return
	}

	u, err := dao.GetUserAccount(vip.User)
	if u.Name == "" || err != nil {
		prese.ResJSON(c, 400, "用户不存在")
		return
	}

	context := Context{}
	var expri int64
	switch vip.SetMeal {
	case 0:
		context.SetDiscountStrategy(RegularCustomerDiscount{})
	case 1:
		context.SetDiscountStrategy(MonthPremiumCustomerDiscount{})
	case 2:
		context.SetDiscountStrategy(QuaretrPremiumCustomerDiscount{})
	case 3:
		context.SetDiscountStrategy(YearPremiumCustomerDiscount{})
	case -1:
		context.SetDiscountStrategy(PermanentPremiumCustomerDiscount{})
	}

	expri = context.CalculateDiscount()
	u.ExpiredTime = time.Now().Add(time.Duration(expri) * 24 * time.Hour).Unix() //effects time after
	u.SetMeal = vip.SetMeal

	err = dao.UpdateAccount(u)
	if err != nil {
		prese.ResJSON(c, 500, "服务故障")
		return
	}

	prese.ResJSON(c, 200, "更新成功")
}
