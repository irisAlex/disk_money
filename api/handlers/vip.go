package api

import (
	"money/dao"
	"money/model"
	"money/pkg/aes"
	"money/pkg/prese"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 设置密钥
const KEY = "money10000@#$%li"

type resCardData struct {
	SetMeal  int   `json:"set_meal"`
	ExriTime int64 `json:"exri_time"`
}

func GrantVip(c *gin.Context) {
	vip := &model.Vip{}
	prese.ParseJSON(c, &vip)

	if vip.User == "" || vip.CardKey == "" {
		prese.ResJSON(c, 400, "参数有问题")
		return
	}

	u, err := dao.GetUserAccount(vip.User)
	if err != nil || u.Name == "" {
		prese.ResJSON(c, 400, "用户不存在")
		return
	}

	context := Context{}
	var expri int64
	decryptedCardNumber, err := aes.VipCardDecrypt(vip.CardKey, KEY)
	if err != nil {
		prese.ResJSON(c, 400, "解密失败")
		return
	}
	card, err := dao.GetCardInfoByKey(decryptedCardNumber)

	if err != nil {
		prese.ResJSON(c, 400, "密钥不存在")
		return
	}

	if !card.Available {
		prese.ResJSON(c, 400, "此卡密已经被兑换过了或不存在，请换其他的卡密")
		return
	}

	switch card.SetMeal {
	case 1:
		context.SetDiscountStrategy(MonthPremiumCustomerDiscount{})
	case 2:
		context.SetDiscountStrategy(QuaretrPremiumCustomerDiscount{})
	case 3:
		context.SetDiscountStrategy(YearPremiumCustomerDiscount{})
	case 4:
		context.SetDiscountStrategy(PermanentPremiumCustomerDiscount{})
	}

	expri = context.CalculateDiscount()
	u.ExpiredTime = time.Now().Add(time.Duration(expri) * 24 * time.Hour).Unix() //effects time after
	u.SetMeal = card.SetMeal
	err = dao.UpdateAccount(u)
	if err != nil {
		prese.ResJSON(c, 500, "服务故障")
		return
	}

	card.Available = false
	err = dao.UpdateCardInfo(card)
	if err != nil {
		prese.ResJSON(c, 500, "服务故障")
		return
	}

	prese.ResJSON(c, 200, &resCardData{
		ExriTime: u.ExpiredTime,
		SetMeal:  u.SetMeal,
	})
}

type ResCardKey struct {
	Key string `json:"key"`
}

func GenerateCardKey(c *gin.Context) {
	tye := c.Query("card_type")
	if tye == "" {
		prese.ResJSON(c, 500, "非法访问")
		return
	}

	cardNumber, err := aes.GenerateRandomString(16)
	if err != nil {
		prese.ResJSON(c, 500, "生成充值卡号时发生错误:")
		return
	}
	// 加密充值卡号
	encryptedCardNumber, err := aes.VipCardEncrypt(cardNumber, KEY)
	if err != nil {
		prese.ResJSON(c, 500, "加密时发生错误:")
		return
	}

	Number, _ := strconv.Atoi(tye)
	dao.InsertCard(&model.Card{
		CardSecret: cardNumber,
		SetMeal:    Number,
		Available:  true,
	})

	prese.ResJSON(c, 200, &ResCardKey{
		Key: encryptedCardNumber,
	})

}
