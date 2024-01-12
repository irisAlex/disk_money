package api

const (
	DAT       int64 = 1
	MONTH     int64 = 30
	QUARTER   int64 = 90
	YEAR      int64 = 365
	PERMANENT int64 = 100000
)

// DiscountStrategy 接口定义了折扣策略
type DiscountStrategy interface {
	ApplyDiscount() int64
	ApplyPermission(int) bool
}

// RegularCustomerDiscount 实现了普通会员的折扣策略
type RegularCustomerDiscount struct{}

func (r RegularCustomerDiscount) ApplyDiscount() int64 {
	return DAT
}

func (r RegularCustomerDiscount) ApplyPermission(dayDdownTime int) bool {
	if dayDdownTime > 0 {
		return false
	}
	return true
}

// MonthPremiumCustomerDiscount 实现了月高级会员的折扣策略
type MonthPremiumCustomerDiscount struct{}

func (m MonthPremiumCustomerDiscount) ApplyDiscount() int64 {
	return MONTH // 10% 折扣
}

func (m MonthPremiumCustomerDiscount) ApplyPermission(dayDdownTime int) bool {
	if dayDdownTime > 50 {
		return false
	}
	return true
}

type YearPremiumCustomerDiscount struct{}

func (y YearPremiumCustomerDiscount) ApplyDiscount() int64 {
	return YEAR // 10% 折扣
}

func (y YearPremiumCustomerDiscount) ApplyPermission(dayDdownTime int) bool {
	if dayDdownTime > 200 {
		return false
	}
	return true
}

type QuaretrPremiumCustomerDiscount struct{}

func (q QuaretrPremiumCustomerDiscount) ApplyDiscount() int64 {
	return QUARTER // 10% 折扣
}

func (q QuaretrPremiumCustomerDiscount) ApplyPermission(dayDdownTime int) bool {
	if dayDdownTime > 150 {
		return false
	}
	return true
}

type PermanentPremiumCustomerDiscount struct{}

func (p PermanentPremiumCustomerDiscount) ApplyDiscount() int64 {
	return PERMANENT // 10% 折扣
}

func (p PermanentPremiumCustomerDiscount) ApplyPermission(dayDdownTime int) bool {
	if dayDdownTime > 600 {
		return false
	}
	return true
}

// Context 封装了折扣策略
type Context struct {
	discountStrategy DiscountStrategy
}

// SetDiscountStrategy 方法用于设置折扣策略
func (c *Context) SetDiscountStrategy(strategy DiscountStrategy) {
	c.discountStrategy = strategy
}

// CalculateDiscount 方法使用当前的折扣策略计算折扣后的金额
func (c *Context) CalculateDiscount() int64 {
	return c.discountStrategy.ApplyDiscount()
}

func (c *Context) IsPermission(downTime int) bool {
	return c.discountStrategy.ApplyPermission(downTime)
}
