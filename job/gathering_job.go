package job

import (
	"github.com/astaxie/beego/orm"
	//"strings"
	"github.com/astaxie/beego"
	"time"
)

	//var lists []*models.KcUserAddress
	//var offset int = 0
	//o := orm.NewOrm()
	//beego.Warn("gathering cron start")
	//for {
	//	num, err := o.QueryTable("kc_user_address").Filter("currency", "ETH").Limit(100, offset).All(&lists, "uid", "address")
	//	if err != nil {
	//		beego.Error(err)
	//		time.Sleep(time.Second * 10)
	//		continue
	//	}
	//	if num > 0 {
	//		for _, l := range lists {
	//			resp, err := ethServer(l.Address)
	//			if err != nil {
	//				ethSleep()
	//				continue
	//			}
	//			for _, trans := range resp.Result {
	//				if strings.ToLower(trans.To) == strings.ToLower(l.Address) {
	//					ethDiposit(l.Uid, trans)
	//				}
	//			}
	//			ethSleep()
	//		}
	//		offset += 100
	//	} else {
	//		offset = 0
	//		ethSleep()
	//	}
	//}

func gatheringSleep() {
	time.Sleep(time.Minute)
}

func GatheringJob() {
	var list orm.ParamsList
	o := orm.NewOrm()
	beego.Error("gathering cron start")
	for {
		num, err := o.QueryTable("kc_subscription_submission").Filter("status", 0).ValuesFlat(&list, "txid")
		if err != nil {
			beego.Error(err)
			time.Sleep(time.Second * 10)
			continue
		}
		if num > 0 {
			for _, l := range list {
				if l != "" {
					var c orm.ParamsList
					count, err := o.Raw("SELECT count(1) count FROM kc_subscription_submission WHERE txid = ? GROUP BY `order`", l).ValuesFlat(&c, "count")
					if err != nil {
						beego.Error(err)
						continue
					}
					if count > 1 {
						o.QueryTable("kc_subscription_submission").Filter("txid", l).Update(orm.Params{
							"Warn1": "重复HASH",
						})
					}
				}
			}
		}
		gatheringSleep()
	}
}
