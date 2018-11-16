package models

import (
	"strconv"
	"math"
	"github.com/astaxie/beego"
	"github.com/spf13/cast"
	"kuangchi_backend/utils"
)

type UserDto struct {
	Ident         string
	Username      string
	Password      string //密码
	Code          string //图形验证码
	Captcha       string //邮箱或手机验证码
	CheckPassword string `json:"check_password"`
	NewPassword   string `json:"new_password"` //新密码
	InviteCode    string `json:"invite_code"` //邀请人code
	CountryCode   string `json:"country_code"`
}

type TfCodeDto struct {
	TfCode string `json:"tf_code"`
}

type InfoDto struct {
	Address     string    `json:"address"` //地址
	Desc        string    `json:"desc"        orm:"type(text)"` //简介
	//社交媒体
	Weibo       string    `json:"weibo"`
	Facebook    string    `json:"facebook"`
	Twitter     string    `json:"twitter"`
	Wechat      string    `json:"wechat"`
	Telegram    string    `json:"telegram"`
}

type StatusDto struct {
	Id     int
	Status int
	Desc   string
}

type TradeDto struct {
	Code   string
	Price  float64
	Amount float64
	Remark string
	FundPassword string `json:"fund_password"`
}

type Withdraw struct {
	Address string
}

type Detail struct {
	GccNum      int
	GccWithdraw string
}

type ShareStruct struct {
	Extension      []string           //个人推广
	Competition    map[string]int64   //竞赛推广
}

func (ss *ShareStruct) Init() {
	ss.Extension = []string{}
	ss.Competition = map[string]int64{}
}

func (ss *ShareStruct) AppendExtension(amountStr string) {
	ss.Extension = append(ss.Extension, amountStr)
}

func (ss *ShareStruct) AddCompetition(key, amountStr string) {
	amounts, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		beego.Error(err)
	}
	if ss.Competition == nil {
		ss.Competition = map[string]int64{}
	}
	ss.Competition[key] += amounts
}

func (ss *ShareStruct) ExtensionRate() float64 {
	var res int64 = 0
	for _, amount := range ss.Extension {
		amountFloat := cast.ToInt64(amount)
		res += amountFloat
	}
	return float64(res)
}

//大于1万 提取1万乘10 其余按原值计算；小于1万 直接乘10
func (ss *ShareStruct) plus(amount int64) int64 {
	if amount <= 10000 {
		return amount * 10
	}
	return amount + 90000
}

func (ss *ShareStruct) CompetitionRate() float64 {
	var aMax int64 = 0
	var res int64 = 0
	for _, amount := range ss.Competition {
		if amount > aMax {
			aMax = amount
		}
		res += ss.plus(amount)
	}
	ceil := math.Cbrt(float64(aMax))
	return utils.ShowFloat(float64(res - ss.plus(aMax)) + ceil, 2)
}
