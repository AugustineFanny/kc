package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"errors"
	"strings"
	"time"
	"fmt"
)

func GetSUsers() (list []*KcSUser) {
	o := orm.NewOrm()
	if _, err := o.QueryTable("kc_s_user").All(&list); err != nil {
		beego.Error(err)
	}
	return list
}

func GetSUserByUid(uid int) *KcSUser {
	o := orm.NewOrm()
	sUser := KcSUser{Id: uid}
	if err := o.Read(&sUser); err != nil {
		return nil
	}
	return &sUser
}

func GetSUserByUsername(username string) *KcSUser {
	o := orm.NewOrm()
	sUser := KcSUser{Username: username}
	if err := o.Read(&sUser, "username"); err != nil {
		return nil
	}
	return &sUser
}

func CreateSUser(username, parent string) error {
	sUser := KcSUser{Username: username}
	o := orm.NewOrm()
	if parent != "" {
		p := GetSUserByUsername(parent)
		if p == nil {
			return errors.New("无效邀请人")
		}
		sUser.InviterId = p.Id
		sUser.Parents = p.DoAsParents()
		sUser.Inviter = p.Username
		sUser.Pnames = p.DoAsParentsName()
	}
	if _, err := o.Insert(&sUser); err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func ChangeSUser(uid int, parent string) error {
	sUser := KcSUser{Id: uid}
	o := orm.NewOrm()
	if err := o.Read(&sUser); err != nil {
		return err
	}
	if sUser.Username == parent {
		return errors.New("推荐人不能为自己")
	}
	if parent == "" {
		sUser.InviterId = 0
		sUser.Parents = ""
		sUser.Inviter = ""
		sUser.Pnames = ""
	} else {
		p := GetSUserByUsername(parent)
		if p == nil {
			return errors.New("无效邀请人")
		}
		names := strings.Split(p.Pnames, ",")
		for _, name := range names {
			if name == sUser.Username {
				return errors.New("推荐人不能为自己的子节点")
			}
		}
		sUser.InviterId = p.Id
		sUser.Parents = p.DoAsParents()
		sUser.Inviter = p.Username
		sUser.Pnames = p.DoAsParentsName()
	}
	if _, err := o.Update(&sUser); err != nil {
		return err
	}
	o.QueryTable("kc_s_activity").Filter("uid", uid).Update(orm.Params{
		"InviterId": sUser.InviterId,
		"Parents": sUser.Parents,
		"Inviter": sUser.Inviter,
		"Pnames": sUser.Pnames,
	})
	return nil
}

func DeleteSUser(uid int) error {
	o := orm.NewOrm()
	o.Begin()
	if _, err := o.QueryTable("kc_s_user").Filter("id", uid).Delete(); err !=nil {
		o.Rollback()
		return err
	}
	if _, err := o.Raw(`DELETE FROM kc_s_user WHERE find_in_set(?, parents)`, uid).Exec(); err != nil {
		beego.Error(uid)
		o.Rollback()
		return err
	}
	if _, err := o.QueryTable("kc_s_activity").Filter("uid", uid).Delete(); err != nil {
		beego.Error(uid)
		o.Rollback()
		return err
	}
	if _, err := o.Raw(`DELETE FROM kc_s_activity WHERE find_in_set(?, parents)`, uid).Exec(); err != nil {
		beego.Error(uid)
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func GetSActivities() (list []*KcSActivity) {
	o := orm.NewOrm()
	if _, err := o.QueryTable("kc_s_activity").All(&list); err != nil {
		beego.Error(err)
	}
	return list
}

func CreateSActivity(activity KcSActivity) error {
	if activity.Subscription < 0 {
		return errors.New("认购大于0")
	}
	if activity.In < 0 {
		return errors.New("转入大于0")
	}
	if activity.Out < 0 {
		return errors.New("转出大于0")
	}
	cur := GetSCurrency()
	if activity.Lock != 0 && activity.Lock < cur.MinLock {
		return errors.New(fmt.Sprintf("最小锁仓%dFET", cur.MinLock))
	}
	startDate := time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC)
	if activity.Date.Sub(startDate) < 0 {
		return errors.New("日期应大于6月1日")
	}
	footstoneExpireDate := startDate.AddDate(0, 0, 60 + 60)
	activity.ExpireDate = activity.Date.AddDate(0, 0, 60)
	if activity.ExpireDate.Sub(footstoneExpireDate) < 0 {
		activity.ExpireDate = footstoneExpireDate
	}
	activity.handleExpireDates()
	p := GetSUserByUsername(activity.Username)
	if p == nil {
		return errors.New("无效用户")
	}
	if activity.In > 0 {
		if activity.InSource == p.Username {
			return errors.New("来源用户不能为自己")
		}
		inSourceUser := GetSUserByUsername(activity.InSource)
		if inSourceUser == nil {
			return errors.New("无效来源用户")
		}
		activity.InSourceUid = inSourceUser.Id
	}
	if activity.Out > 0 {
		if activity.OutDest == p.Username {
			return errors.New("目标用户不能为自己")
		}
		outSourceUser := GetSUserByUsername(activity.OutDest)
		if outSourceUser == nil {
			return errors.New("无效目标用户")
		}
		activity.OutDestUid = outSourceUser.Id
	}
	activity.Uid = p.Id
	activity.InviterId = p.InviterId
	activity.Parents = p.Parents
	activity.Inviter = p.Inviter
	activity.Pnames = p.Pnames
	o := orm.NewOrm()
	if _, err := o.Insert(&activity); err != nil {
		beego.Error(err)
		return err
	}
	return nil
}

func GetSActivity(id int) (*KcSActivity, error) {
	o := orm.NewOrm()
	activity := KcSActivity{Id: id}
	if err := o.Read(&activity); err != nil {
		return nil, err
	}
	return &activity, nil
}

func ChangeSActivity(activity KcSActivity) error {
	cur := GetSCurrency()
	if activity.Lock != 0 && activity.Lock < cur.MinLock {
		return errors.New(fmt.Sprintf("最小锁仓%dFET", cur.MinLock))
	}
	startDate := time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC)
	if activity.Date.Sub(startDate) < 0 {
		return errors.New("日期应大于6月1日")
	}
	o := orm.NewOrm()
	footstoneExpireDate := startDate.AddDate(0, 0, 60 + 60)
	activity.ExpireDate = activity.Date.AddDate(0, 0, 60)
	if activity.ExpireDate.Sub(footstoneExpireDate) < 0 {
		activity.ExpireDate = footstoneExpireDate
	}
	activity.handleExpireDates()
	if activity.In > 0 {
		if activity.InSource == activity.Username {
			return errors.New("来源用户不能为自己")
		}
		inSourceUser := GetSUserByUsername(activity.InSource)
		if inSourceUser == nil {
			return errors.New("无效来源用户")
		}
		activity.InSourceUid = inSourceUser.Id
	}
	if activity.Out > 0 {
		if activity.OutDest == activity.Username {
			return errors.New("目标用户不能为自己")
		}
		outSourceUser := GetSUserByUsername(activity.OutDest)
		if outSourceUser == nil {
			return errors.New("无效目标用户")
		}
		activity.OutDestUid = outSourceUser.Id
	}
	if _, err := o.Update(&activity); err != nil {
		return err
	}
	return nil
}

func DeleteSActivity(id int) error {
	o := orm.NewOrm()
	if _, err := o.QueryTable("kc_s_activity").Filter("id", id).Delete(); err !=nil {
		o.Rollback()
		return err
	}
	return nil
}

func ChangeSCurrency(miningInterestRate, competitionInterestRate float64) error {
	o := orm.NewOrm()
	cur := KcSCurrency{Currency: "FET"}
	if err := o.Read(&cur, "Currency"); err != nil {
		beego.Error(err)
		cur.MinLock = 100
		cur.CInterestRate = miningInterestRate
		cur.MInterestRate = competitionInterestRate
		if _, err := o.Insert(&cur); err != nil {
			beego.Error(err)
		}
	} else {
		cur.MinLock = 100
		cur.MInterestRate = miningInterestRate
		cur.CInterestRate = competitionInterestRate
		if _, err := o.Update(&cur); err != nil {
			beego.Error(err)
		}
	}
	return nil
}

func GetSCurrency() KcSCurrency {
	o := orm.NewOrm()
	cur := KcSCurrency{Currency: "FET"}
	if err := o.Read(&cur, "Currency"); err != nil {
		beego.Error(err)
		cur.MinLock = 100
		cur.CInterestRate = 0.001
		cur.MInterestRate = 0.001
		if _, err := o.Insert(&cur); err != nil {
			beego.Error(err)
		}
	}
	return cur
}

func DeleteAllSUsers() error {
	o := orm.NewOrm()
	o.Begin()
	if _, err := o.Raw(`DELETE FROM kc_s_user`).Exec(); err !=nil {
		o.Rollback()
		return err
	}
	if _, err := o.Raw(`DELETE FROM kc_s_activity`).Exec(); err !=nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func DeleteAllSActivities() error {
	o := orm.NewOrm()
	if _, err := o.Raw(`DELETE FROM kc_s_activity`).Exec(); err !=nil {
		o.Rollback()
		return err
	}
	return nil
}
