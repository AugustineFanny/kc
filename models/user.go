package models

import (
	"github.com/astaxie/beego/orm"
	"kuangchi_backend/result"
	"github.com/astaxie/beego"
)

type userDetail struct {
	Username       string
	Email          string
}

func GetUserDetail(username string) (map[string]interface{}, error) {
	o := orm.NewOrm()
	user := User{Username: username}
	if err := o.Read(&user, "Username"); err != nil {
		return nil, result.ErrCode(100102)
	}
	profile := Profile{Uid: user.Id}
	if err := o.Read(&profile, "Uid"); err != nil {
		beego.Error(err)
		return nil, result.ErrCode(100102)
	}
	detail := map[string]interface{}{
		"username": user.Username,
		"avatar": user.Avatar,
		"email": user.Verified("email"),
		"mobile": user.Verified("mobile"),
		"role": user.Verified("role"),
		"kyc": user.Verified("kyc"),
		"trade_times": profile.TradeTimes,
		"average_pass": profile.AveragePass,
		"trusted": profile.Trusted,
	}
	return detail, nil
}
