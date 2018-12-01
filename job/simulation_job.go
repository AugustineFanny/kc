package job
//
//import (
//	"github.com/spf13/cast"
//	"kuangchi_backend/models"
//	"kuangchi_backend/utils"
//	"github.com/astaxie/beego/orm"
//	"time"
//	"errors"
//)
//
//type calculationRes struct {
//	Username      string   `json:"username"`
//	SubscriptionA float64  `json:"subscription_all"`   //认购总计
//	MiningA       float64  `json:"mining_all"`         //生息总计
//	ShareA        float64  `json:"share_all"`          //奖励总计
//	UnlockA       float64  `json:"unlock_all"`         //释放总计
//	Lock          float64  `json:"lock"`               //锁定总计
//	Active        float64  `json:"active"`             //流通总计
//	SubscriptionD float64  `json:"subscription_day"`   //当日认购
//	MiningD       float64  `json:"mining_day"`         //当日生息
//	ShareD        float64  `json:"share_day"`          //当日奖励
//	UnlockD       float64  `json:"unlock_day"`         //当日释放
//	LockD         float64  `json:"lock_day"`           //新增锁定
//	ActiveD       float64  `json:"active_day"`         //新增流通
//	InA           float64  `json:"in_all"`             //转入总计
//	OutA		  float64  `json:"out_all"`            //转出总计
//	AddlockA      float64  `json:"addlock_all"`        //转移总计
//	InD 		  float64  `json:"in_day"`             //当日转入
//	OutD          float64  `json:"out_day"`            //当日转出
//	AddlockD      float64  `json:"addlock_day"`        //当日转移
//	Reward1A      float64  `json:"reward1_all"`        //个人推广总计
//	Reward2A      float64  `json:"reward2_all"`        //竞赛推广总计
//	Reward1D      float64  `json:"reward1_day"`        //个人推广
//	Reward2D      float64  `json:"reward2_day"`        //竞赛推广
//	Date          string   `json:"date"`
//}
//
//func (u *calculationRes) format() {
//	u.SubscriptionA = utils.ShowFloat(u.SubscriptionA, 6)
//	u.MiningA       = utils.ShowFloat(u.MiningA, 6)
//	u.ShareA 		= utils.ShowFloat(u.ShareA, 6)
//	u.UnlockA 		= utils.ShowFloat(u.UnlockA, 6)
//	u.Lock 			= utils.ShowFloat(u.Lock, 6)
//	u.Active 		= utils.ShowFloat(u.Active, 6)
//	u.SubscriptionD = utils.ShowFloat(u.SubscriptionD, 6)
//	u.MiningD 		= utils.ShowFloat(u.MiningD, 6)
//	u.ShareD 		= utils.ShowFloat(u.ShareD, 6)
//	u.UnlockD 		= utils.ShowFloat(u.UnlockD, 6)
//	u.LockD 		= utils.ShowFloat(u.LockD, 6)
//	u.ActiveD 		= utils.ShowFloat(u.ActiveD, 6)
//	u.InA 			= utils.ShowFloat(u.InA, 6)
//	u.OutA 			= utils.ShowFloat(u.OutA, 6)
//	u.AddlockA 		= utils.ShowFloat(u.AddlockA, 6)
//	u.InD 			= utils.ShowFloat(u.InD, 6)
//	u.OutD 			= utils.ShowFloat(u.OutD, 6)
//	u.AddlockD 		= utils.ShowFloat(u.AddlockD, 6)
//	u.Reward1A 		= utils.ShowFloat(u.Reward1A, 6)
//	u.Reward1D 		= utils.ShowFloat(u.Reward1D, 6)
//	u.Reward2A 		= utils.ShowFloat(u.Reward2A, 6)
//	u.Reward2D 		= utils.ShowFloat(u.Reward2D, 6)
//}
//
//func handleCalculationInit(activities []*models.KcSActivity, res map[int]*calculationRes) {
//	for _, activity := range activities {
//		if res[activity.Uid] == nil {
//			res[activity.Uid] = &calculationRes{Username: activity.Username}
//		}
//		res[activity.Uid].SubscriptionA += activity.Subscription
//		res[activity.Uid].Lock += activity.Mining()
//		res[activity.Uid].SubscriptionD = activity.Subscription
//		res[activity.Uid].LockD = activity.Mining()
//		res[activity.Uid].InA += activity.In
//		res[activity.Uid].OutA += activity.Out
//		res[activity.Uid].AddlockA += activity.Lock
//		res[activity.Uid].InD = activity.In
//		res[activity.Uid].OutD = activity.Out
//		res[activity.Uid].AddlockD = activity.Lock
//		res[activity.Uid].ActiveD = activity.In - activity.Out - activity.Lock
//		res[activity.Uid].Active += activity.In - activity.Out - activity.Lock
//
//		if activity.In > 0 {
//			if res[activity.InSourceUid] == nil {
//				res[activity.InSourceUid] = &calculationRes{Username: activity.InSource}
//			}
//			res[activity.InSourceUid].ActiveD -= activity.In
//			res[activity.InSourceUid].Active -= activity.In
//		}
//		if activity.Out > 0 {
//			if res[activity.OutDestUid] == nil {
//				res[activity.OutDestUid] = &calculationRes{Username: activity.OutDest}
//			}
//			res[activity.OutDestUid].ActiveD += activity.Out
//			res[activity.OutDestUid].Active += activity.Out
//		}
//	}
//}
//
//func handleCalculationExpire(date time.Time, res map[int]*calculationRes) {
//	var params []orm.Params
//	o := orm.NewOrm()
//	qb, _ := orm.NewQueryBuilder("mysql")
//	qb.Select("uid", "subscription", "`lock`").From("kc_s_activity").Where("find_in_set(?, expire_dates)")
//	sql := qb.String()
//	o.Raw(sql, date.Format("2006-01-02")).Values(&params)
//	for _, node := range params {
//		uid := cast.ToInt(node["uid"].(string))
//		subscription := cast.ToFloat64(node["subscription"].(string))
//		lock := cast.ToFloat64(node["lock"].(string))
//		res[uid].Active += (subscription + lock) * 0.1
//		res[uid].UnlockD = (subscription + lock) * 0.1
//		res[uid].UnlockA += (subscription + lock) * 0.1
//		res[uid].Lock -= (subscription + lock) * 0.1
//	}
//}
//
//func handleCalculationMining(cur models.KcSCurrency, res map[int]*calculationRes) {
//	var rate float64 = 0  //待分配总量
//	var percentage float64 = 1 //获得百分比 控制总量在100000
//	for _, node := range res {
//		rate += node.Lock * cur.MInterestRate
//	}
//	if rate > 100000 {
//		percentage = 100000 / rate
//	}
//	for _, node := range res {
//		mining := utils.ShowFloat(node.Lock * cur.MInterestRate * percentage, 6)
//		node.MiningA += mining
//		node.MiningD = mining
//		node.Active += mining
//		node.ActiveD = mining
//	}
//}
//
//func handleCalculationShare(cur models.KcSCurrency, activities []*models.KcSActivity, res map[int]*calculationRes) {
//	process := []orm.Params{}
//	for _, activity := range activities {
//		if activity.InviterId != 0 {
//			process = append(process, orm.Params{
//				"uid": cast.ToString(activity.Uid),
//				"inviter_id": cast.ToString(activity.InviterId),
//				"amounts": cast.ToString(int(activity.Mining())),
//				"parents": activity.Parents,
//				"lid": "0",
//			})
//		}
//	}
//	shares, _ := handleShare(process)
//	var rate float64 = 0  //待分配总量
//	var percentage float64 = 1 //获得百分比 控制总量在100000
//	var phase1 = map[int][2]float64{}
//	for uidStr, shareS := range shares {
//		reward1 := shareS.ExtensionRate() * 0.05
//		factReward1 := reward1
//		reward2 := shareS.CompetitionRate() * cur.CInterestRate
//		factReward2 := reward2
//
//		uid := cast.ToInt(uidStr)
//		if res[uid] == nil {
//			sUser := models.GetSUserByUid(uid)
//			res[uid] = &calculationRes{Username: sUser.Username}
//		}
//
//		limit := res[uid].Lock * 3
//		if factReward1 + factReward2 > limit {
//			if factReward1 > limit {
//				factReward1 = limit
//				factReward2 = 0
//			} else {
//				factReward2 = limit - factReward1
//			}
//		}
//		rate += factReward2
//		phase1[uid] = [2]float64{factReward1, factReward2}
//	}
//	if rate > 100000 {
//		percentage = 100000 / rate
//	}
//	for uid, data := range phase1 {
//		reward1 := utils.ShowFloat(data[0], 6)
//		reward2 := utils.ShowFloat(data[1] * percentage, 6)
//		res[uid].Reward1D = reward1
//		res[uid].Reward2D = reward2
//		res[uid].Reward1A += reward1
//		res[uid].Reward2A += reward2
//		res[uid].ShareA += reward1 + reward2
//		res[uid].ShareD = reward1 + reward2
//		res[uid].Active += reward1 + reward2
//		res[uid].ActiveD += reward1 + reward2
//	}
//}
//
//func handleCalculation(date time.Time, calculationResMap map[int]*calculationRes) error {
//	for _, node := range calculationResMap {
//		node.SubscriptionD = 0
//		node.MiningD = 0
//		node.ShareD = 0
//		node.UnlockD = 0
//		node.LockD = 0
//		node.ActiveD = 0
//		node.InD = 0
//		node.OutD = 0
//		node.AddlockD = 0
//	}
//	var activities []*models.KcSActivity
//	o := orm.NewOrm()
//	o.QueryTable("kc_s_activity").Filter("date", date).All(&activities)
//	cur := models.GetSCurrency()
//	miningStartDate := time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, 60)
//	handleCalculationInit(activities, calculationResMap)
//	if date.Sub(miningStartDate) > 0 {
//		handleCalculationMining(cur, calculationResMap)
//	}
//	handleCalculationShare(cur, activities, calculationResMap)
//	handleCalculationExpire(date, calculationResMap)
//	for _, calculation := range calculationResMap {
//		calculation.Date = date.Format("2006-01-02")
//		calculation.format()
//		if calculation.Active < 0 {
//			return errors.New("用户余额不足：" + calculation.Username + ", " + calculation.Date)
//		}
//	}
//	return nil
//}
//
//func Calculation(username, startDateS, endDateS string) (res []*calculationRes, err error) {
//	startDate, err := cast.ToTimeE(startDateS)
//	if err != nil {
//		return []*calculationRes{}, err
//	}
//	footStoreStartDate := time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC)
//	timeDuration := startDate.Sub(footStoreStartDate)
//	if timeDuration < 0 {
//		return []*calculationRes{}, errors.New("日期应大于6月1日")
//	}
//	calculationResMap := map[int]*calculationRes{}
//	if username != "" {
//		sUser := models.GetSUserByUsername(username)
//		if sUser == nil {
//			return []*calculationRes{}, errors.New("无效用户")
//		}
//		endDate, err := cast.ToTimeE(endDateS)
//		if err != nil {
//			return []*calculationRes{}, err
//		}
//		for i := footStoreStartDate; endDate.Sub(i) >= 0; i = i.AddDate(0, 0, 1) {
//			err = handleCalculation(i, calculationResMap)
//			if startDate.Sub(i) <= 0 {
//				if calculationResMap[sUser.Id] != nil {
//					var calculation = new(calculationRes)
//					*calculation = *calculationResMap[sUser.Id]
//					res = append(res, calculation)
//				}
//			}
//		}
//	} else {
//		for i := footStoreStartDate; startDate.Sub(i) >= 0; i = i.AddDate(0, 0, 1) {
//			err = handleCalculation(i, calculationResMap)
//		}
//		for _, calculation := range calculationResMap {
//			res = append(res, calculation)
//		}
//	}
//	if res == nil {
//		return []*calculationRes{}, err
//	}
//	return res, err
//}
//
//func Stat(startDateS, endDateS string) (res []*calculationRes, err error) {
//	startDate, err := cast.ToTimeE(startDateS)
//	if err != nil {
//		return []*calculationRes{}, err
//	}
//	endDate, err := cast.ToTimeE(endDateS)
//	if err != nil {
//		return []*calculationRes{}, err
//	}
//	footStoreStartDate := time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC)
//	timeDuration := startDate.Sub(footStoreStartDate)
//	if timeDuration < 0 {
//		return []*calculationRes{}, errors.New("日期应大于6月1日")
//	}
//	calculationResMap := map[int]*calculationRes{}
//	for i := footStoreStartDate; endDate.Sub(i) >= 0; i = i.AddDate(0, 0, 1) {
//		err = handleCalculation(i, calculationResMap)
//		if startDate.Sub(i) <= 0 {
//			calculation := calculationRes{}
//			for _, c := range calculationResMap {
//				calculation.SubscriptionA += c.SubscriptionA
//				calculation.MiningA += c.MiningA
//				calculation.ShareA += c.ShareA
//				calculation.UnlockA += c.UnlockA
//				calculation.Lock += c.Lock
//				calculation.Active += c.Active
//				calculation.SubscriptionD += c.SubscriptionD
//				calculation.MiningD += c.MiningD
//				calculation.ShareD += c.ShareD
//				calculation.UnlockD += c.UnlockD
//				calculation.LockD += c.LockD
//				calculation.ActiveD += c.ActiveD
//				calculation.InA += c.InA
//				calculation.OutA += c.OutA
//				calculation.AddlockA += c.AddlockA
//				calculation.InD += c.InD
//				calculation.OutD += c.OutD
//				calculation.AddlockD += c.AddlockD
//				calculation.Reward1A += c.Reward1A
//				calculation.Reward1D += c.Reward1D
//				calculation.Reward2A += c.Reward2A
//				calculation.Reward2D += c.Reward2D
//				calculation.Date = i.Format("2006-01-02")
//				calculation.format()
//			}
//			res = append(res, &calculation)
//		}
//	}
//	if res == nil {
//		return []*calculationRes{}, err
//	}
//	return res, err
//}
