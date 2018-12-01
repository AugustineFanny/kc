package job
//
//import (
//	"github.com/astaxie/beego/orm"
//	"github.com/astaxie/beego"
//	"github.com/spf13/cast"
//	"kuangchi_backend/models"
//	"kuangchi_backend/utils"
//	"strings"
//	"time"
//)
//
//func unlockOrder(cur *models.KcCurrency) {
//	var lockeds []models.KcLocked
//	o := orm.NewOrm()
//	nums, err := o.QueryTable("kc_locked").
//		Filter("currency", cur.Currency).Filter("status", 0).Filter("expire_date__lte", time.Now()).All(&lockeds)
//	if err != nil {
//		beego.Error(err)
//		return
//	}
//	if nums == 0 {
//		return
//	}
//	for _, locked := range lockeds {
//		o.Begin()
//		unlockAmount := locked.TotalAmount / 10
//		if locked.Amount < unlockAmount {
//			unlockAmount = locked.Amount
//			locked.Amount = 0
//		} else {
//			locked.Amount -= unlockAmount
//		}
//		locked.UnlockNum += 1
//		if locked.Amount == 0 || locked.UnlockNum == 10 {
//			locked.Status = 1
//		} else {
//			locked.ExpireDate = locked.ExpireDate.AddDate(0, 0, 30)
//		}
//		if _, err := o.Update(&locked, "Amount", "UnlockNum", "Status", "ExpireDate"); err != nil {
//			beego.Error(err)
//			o.Rollback()
//			return
//		}
//		wallet := models.KcWallet{Uid: locked.Uid, Currency: locked.Currency}
//		if err := o.ReadForUpdate(&wallet, "Uid", "Currency"); err != nil {
//			beego.Error(err)
//			o.Rollback()
//			return
//		}
//		wallet.LockAmount -= unlockAmount
//		wallet.MiningAmount -= unlockAmount
//		if wallet.LockAmount < 0 {
//			beego.Error("error: id:", wallet.Id, wallet.LockAmount, wallet.MiningAmount)
//			wallet.LockAmount = 0
//			wallet.MiningAmount = 0
//		}
//		if _, err := o.Update(&wallet, "LockAmount", "MiningAmount"); err != nil {
//			beego.Error(err)
//			return
//		}
//		o.Commit()
//	}
//}
//
//func expireSubscription() {
//	o := orm.NewOrm()
//	_, err := o.Raw(`UPDATE kc_subscription SET status = 4 WHERE status = 0 AND HOUR(timediff(now(), create_time)) > 24`).Exec()
//	if err != nil {
//		beego.Error(err)
//	}
//}
//
//func getValidNodes(currency string) []orm.Params {
//	var maps []orm.Params
//	o := orm.NewOrm()
//	sql := `SELECT uid, FLOOR(mining_amount) amount FROM kc_wallet WHERE currency = ? AND mining_amount >= 1 ORDER BY mining_amount`
//	_, err := o.Raw(sql, currency).Values(&maps)
//	if err != nil {
//		return []orm.Params{}
//	}
//	return maps
//}
//
//func mining(cur *models.KcCurrency) {
//	nodes := getValidNodes(cur.Currency)
//	var rate float64 = 0  //待分配总量
//	var percentage float64 = 1 //获得百分比 控制总量在100000
//	for _, node := range nodes {
//		rate += cast.ToFloat64(node["amount"].(string)) * cur.MInterestRate
//	}
//	if rate > 100000 {
//		percentage = 100000 / rate
//	}
//	o := orm.NewOrm()
//	o.Begin()
//	hasError := false
//	for _, node := range nodes {
//		uid := cast.ToInt(node["uid"].(string))
//		amount := cast.ToFloat64(node["amount"].(string))
//		reward := utils.ShowFloat(amount * cur.MInterestRate * percentage, 6)
//		mining := models.KcMining{
//			Uid: uid,
//			Currency: cur.Currency,
//			Reward: reward,
//			Mining: amount,
//			Rate: rate,
//			Percentage: percentage,
//			InterestRate: cur.MInterestRate,
//			Description: "mining",
//		}
//		miningId, err := o.Insert(&mining)
//		if err != nil {
//			beego.Error(err)
//			hasError = true
//			break
//		}
//		remark := cast.ToString(miningId) //备注记录kc_mining id
//		if err := models.AddAmount(o, uid, reward, cur.Currency, "mining", remark); err != nil {
//			beego.Error(err)
//			hasError = true
//			break
//		}
//	}
//	cur.HodlLast = time.Now()
//	if _, err := o.Update(cur, "HodlLast"); err != nil {
//		beego.Error(err)
//		hasError = true
//	}
//	if hasError {
//		o.Rollback()
//		models.SendSMSToMultiAdmin("锁仓奖励发放失败" + time.Now().Format("2006-01-02 15:04:05"))
//	} else {
//		o.Commit()
//		models.SendSMSToMultiAdmin("锁仓奖励发放成功" + time.Now().Format("2006-01-02 15:04:05"))
//	}
//}
//
////func mining(cur *models.KcCurrency) {
////	nodes := getValidNodes(cur.Currency)
////	l := 1
////	mx := 0
////	for index, node := range nodes {
////		if index != 0 && nodes[index]["amount"] == nodes[index - 1]["amount"] {
////			node["m"] = nodes[index - 1]["m"]
////		} else {
////			node["m"] = l
////		}
////		mx += node["m"].(int)
////		l++
////	}
////	o := orm.NewOrm()
////	o.Begin()
////	hasError := false
////	for _, node := range nodes {
////		uid, _ := strconv.Atoi(node["uid"].(string))
////		amount := utils.ShowFloat((float64(node["m"].(int)) / float64(mx)) * 100000, 6)
////		if err := models.AddAmount(o, uid, amount, cur.Currency, "mining", node["amount"].(string)); err != nil {
////			beego.Error(err)
////			hasError = true
////			break
////		}
////	}
////	cur.HodlLast = time.Now()
////	if _, err := o.Update(cur, "HodlLast"); err != nil {
////		beego.Error(err)
////		hasError = true
////	}
////	if hasError {
////		o.Rollback()
////		models.SendSMSToMultiAdmin("锁仓奖励发放失败" + time.Now().Format("2006-01-02 15:04:05"))
////	} else {
////		o.Commit()
////		models.SendSMSToMultiAdmin("锁仓奖励发放成功" + time.Now().Format("2006-01-02 15:04:05"))
////	}
////}
//
//func getValidLocked(currency string) []orm.Params {
//	var maps []orm.Params
//	o := orm.NewOrm()
//	sql := `SELECT kc_locked.id lid, FLOOR(amount) amounts, user.id uid, user.inviter_id, user.parents
//			FROM kc_locked LEFT JOIN user ON kc_locked.uid = user.id
//			WHERE kc_locked.currency = ? AND kc_locked.share = 0 AND kc_locked.amount >= 1`
//	_, err := o.Raw(sql, currency).Values(&maps)
//	if err != nil {
//		beego.Error(err)
//		return []orm.Params{}
//	}
//	return maps
//}
//
//func handleShare(maps []orm.Params) (map[string]*models.ShareStruct, []string) {
//	shares := map[string]*models.ShareStruct{}
//	ids := []string{}
//	for _, m := range maps {
//		uid := m["uid"].(string)
//		inviterId := m["inviter_id"].(string)
//		if shares[inviterId] == nil {
//			shares[inviterId] = &models.ShareStruct{}
//			shares[inviterId].Init()
//		}
//		shares[inviterId].AppendExtension(m["amounts"].(string))
//		parents := strings.Split(m["parents"].(string), ",")
//		parentsLen := len(parents)
//		for index, parent := range parents {
//			track := uid
//			if index < parentsLen - 1 {
//				track = parents[index + 1]
//			}
//			if shares[parent] == nil {
//				shares[parent] = &models.ShareStruct{}
//				shares[parent].Init()
//			}
//			shares[parent].AddCompetition(track, m["amounts"].(string))
//		}
//		ids = append(ids, m["lid"].(string))
//	}
//	return shares, ids
//}
//
//func sl2il(stringList []string) []interface{} {
//	res := make([]interface{}, len(stringList))
//	for i, v := range stringList {
//		res[i] = v
//	}
//	return res
//}
//
////获取未分配推广奖励的锁仓总量
////func getLockedNum(cur *models.KcCurrency) float64 {
////	var list orm.ParamsList
////	o := orm.NewOrm()
////	sql := `SELECT SUM(amount) amounts, currency FROM kc_locked WHERE kc_locked.currency = ? AND kc_locked.share = 0 AND kc_locked.amount >= 1 GROUP BY currency`
////	if _, err := o.Raw(sql, cur.Currency).ValuesFlat(&list, "amounts"); err != nil {
////		beego.Error(err)
////		return 0
////	}
////	if len(list) == 0 {
////		beego.Error("出问题")
////		return 0
////	}
////	amounts, err := strconv.ParseFloat(list[0].(string), 64)
////	if err != nil {
////		beego.Error(err)
////		return 0
////	}
////	return amounts
////}
//
//func share(cur *models.KcCurrency) {
//	ss := getValidLocked(cur.Currency)
//	shares, ids := handleShare(ss)
//	var rate float64 = 0  //待分配总量
//	var percentage float64 = 1 //获得百分比 控制总量在100000
//	var phase1 = map[int][5]float64{}
//	for uidStr, shareS := range shares {
//		reward1 := shareS.ExtensionRate() * 0.05
//		factReward1 := reward1
//		reward2 := shareS.CompetitionRate() * cur.CInterestRate
//		factReward2 := reward2
//
//		uid := cast.ToInt(uidStr)
//		user := models.GetUserById(uid)
//		if user == nil {
//			beego.Error("有问题", uid)
//			continue
//		}
//		miningAmount := models.GetWallet(user, cur.Currency)[0].MiningAmount
//		limit := miningAmount * 3
//		if factReward1 + factReward2 > limit {
//			if factReward1 > limit {
//				factReward1 = limit
//				factReward2 = 0
//			} else {
//				factReward2 = limit - factReward1
//			}
//		}
//		rate += factReward2
//		//列表内容：实得个人推广，实得竞赛推广，锁仓，应得个人推广，应得竞赛推广
//		phase1[uid] = [5]float64{factReward1, factReward2, miningAmount, reward1, reward2}
//	}
//	if rate > 100000 {
//		percentage = 100000 / rate
//	}
//	hasError := false
//	o := orm.NewOrm()
//	o.Begin()
//	for uid, data := range phase1 {
//		reward1 := utils.ShowFloat(data[0], 6)
//		reward2 := utils.ShowFloat(data[1] * percentage, 6)
//		reward := utils.ShowFloat(reward1 + reward2, 6)
//		miningAmount := utils.ShowFloat(data[2], 6)
//		remark := cast.ToString(data[3]) + "," + cast.ToString(data[4])
//		mining := models.KcMining{
//			Uid: uid,
//			Currency: cur.Currency,
//			Reward: reward,
//			Reward1: reward1,
//			Reward2: reward2,
//			Mining: miningAmount,
//			Rate: rate,
//			Percentage: percentage,
//			InterestRate: cur.CInterestRate,
//			Description: "share",
//			Remark: remark,
//		}
//		miningId, err := o.Insert(&mining)
//		if err != nil {
//			beego.Error(err)
//			hasError = true
//			break
//		}
//		if reward == 0 {
//			//无锁仓 留记录 但不给币
//			continue
//		}
//		if err := models.AddAmount(o, uid, reward1, cur.Currency, "share", cast.ToString(miningId)); err != nil {
//			beego.Error(err)
//			hasError = true
//			break
//		}
//		if reward2 != 0 {
//			if err := models.AddAmount(o, uid, reward2, cur.Currency, "competition", cast.ToString(miningId)); err != nil {
//				beego.Error(err)
//				hasError = true
//				break
//			}
//		}
//	}
//	if len(ids) > 0 {
//		if _, err := o.QueryTable("kc_locked").Filter("id__in", sl2il(ids)...).Update(orm.Params{
//			"share": 1,
//		}); err != nil {
//			beego.Error(err)
//			hasError = true
//		}
//	}
//	cur.ShareLast = time.Now()
//	if _, err := o.Update(cur, "ShareLast"); err != nil {
//		beego.Error(err)
//		hasError = true
//	}
//	if hasError {
//		o.Rollback()
//		models.SendSMSToMultiAdmin("推广奖励发放失败" + time.Now().Format("2006-01-02 15:04:05"))
//	} else {
//		o.Commit()
//		models.SendSMSToMultiAdmin("推广奖励发放成功" + time.Now().Format("2006-01-02 15:04:05"))
//	}
//
//	//var mx float64 = 0
//	//ps := map[string]float64{}
//	//for uidStr, amounts := range shares {
//	//	ps[uidStr] = hashrate(amounts)
//	//	mx += ps[uidStr]
//	//}
//	//o := orm.NewOrm()
//	//o.Begin()
//	//hasError := false
//	//for uidStr, m := range ps {
//	//	uid, _ := strconv.Atoi(uidStr)
//	//	//奖池 = 新增锁仓量 * 20%
//	//	spreadPool := getLockedNum(cur) * 0.2
//	//	amount := utils.ShowFloat((m / mx) * spreadPool, 6)
//	//	user := models.GetUserById(uid)
//	//	if user == nil {
//	//		beego.Error("有问题")
//	//		continue
//	//	}
//	//	miningAmount := models.GetWallet(user, cur.Currency)[0].MiningAmount
//	//	if miningAmount == 0 {
//	//		continue
//	//	}
//	//	if amount > miningAmount * 3 {
//	//		amount = miningAmount * 3
//	//	}
//	//	remark := strings.Join(shares[uidStr], ",")
//	//	if err := models.AddAmount(o, uid, amount, cur.Currency, "share", remark); err != nil {
//	//		beego.Error(err)
//	//		hasError = true
//	//		break
//	//	}
//	//	if _, err := o.QueryTable("kc_locked").Filter("id__in", sl2il(ids)...).Update(orm.Params{
//	//		"share": 1,
//	//	}); err != nil {
//	//		beego.Error(err)
//	//		hasError = true
//	//	}
//	//}
//	//cur.ShareLast = time.Now()
//	//if _, err := o.Update(cur, "ShareLast"); err != nil {
//	//	beego.Error(err)
//	//	hasError = true
//	//}
//	//if hasError {
//	//	o.Rollback()
//	//	models.SendSMSToMultiAdmin("推广奖励发放失败" + time.Now().Format("2006-01-02 15:04:05"))
//	//} else {
//	//	o.Commit()
//	//	models.SendSMSToMultiAdmin("推广奖励发放成功" + time.Now().Format("2006-01-02 15:04:05"))
//	//}
//}
//
//func PredistributionMining() []orm.Params {
//	cur := models.GetCurrency("FET")
//	if cur == nil {
//		beego.Error("no FET")
//		return []orm.Params{}
//	}
//	nodes := getValidNodes(cur.Currency)
//	var rate float64 = 0  //待分配总量
//	var percentage float64 = 1 //获得百分比 控制总量在100000
//	for _, node := range nodes {
//		rate += cast.ToFloat64(node["amount"].(string)) * cur.MInterestRate
//	}
//	if rate > 100000 {
//		percentage = 100000 / rate
//	}
//	for _, node := range nodes {
//		amount :=  utils.ShowFloat(cast.ToFloat64(node["amount"].(string)) * cur.MInterestRate * percentage, 6)
//		node["mining"] = amount
//	}
//	return nodes
//}
//
//func PredistributionShareOld() []orm.Params {
//	cur := models.GetCurrency("FET")
//	ss := getValidLocked(cur.Currency)
//	if len(ss) == 0 {
//		return []orm.Params{}
//	}
//	//shares, _ := handleShare(ss)
//	var mx float64 = 0
//	ps := map[string]float64{}
//	//for uidStr, amounts := range shares {
//	//	ps[uidStr] = hashrate(amounts)
//	//	mx += ps[uidStr]
//	//}
//	res := []orm.Params{}
//	for uidStr, m := range ps {
//		amount := utils.ShowFloat(( m / mx ) * 100000, 6)
//		res = append(res, map[string]interface{}{
//			"uid": uidStr,
//			"share": amount,
//			"amount": m,
//		})
//	}
//	return res
//}
//
//func PredistributionShare() []orm.Params {
//	cur := models.GetCurrency("FET")
//	if cur == nil {
//		beego.Error("no FET")
//		return []orm.Params{}
//	}
//	ss := getValidLocked(cur.Currency)
//	shares, _ := handleShare(ss)
//	res := []orm.Params{}
//	for uidStr, shareS := range shares {
//		extensionRate := shareS.ExtensionRate()
//		competitionRate := shareS.CompetitionRate()
//		res = append(res, orm.Params{
//			"uid": uidStr,
//			"extension_rate": extensionRate,
//			"competition_rate": competitionRate,
//			"extension_reward": extensionRate * 0.05,
//			"competition_reward": competitionRate * cur.CInterestRate,
//			})
//	}
//	return res
//}
//
//func MiningJob() {
//	beego.Debug("mining start")
//	c := time.Tick(time.Minute)
//	for now := range c {
//		cur := models.GetCurrency("FET")
//		if cur == nil {
//			beego.Error("error: mining job closed")
//			beego.Error("       no currency")
//			return
//		}
//		//到期解锁 每期解锁10% 1点以后解锁
//		if now.Hour() >= 1 {
//			utils.RedisClient.Set("UnlockOrder", time.Now().Format("2006-01-02 15:04:05"), 10 * 60)
//			beego.Warn("start unlock order")
//			unlockOrder(cur)
//			expireSubscription()
//			utils.RedisClient.Delete("UnlockOrder")
//		}
//
//		//if now.Sub(cur.HodlLast).Minutes() >= 24 * 60 {
//		if now.Sub(cur.HodlLast).Minutes() >= 60 {
//			if _, err := utils.RedisClient.GetString("Hodl"); err != nil {
//				utils.RedisClient.Set("Hodl", time.Now().Format("2006-01-02 15:04:05"), 10 * 60)
//				beego.Warn("start mining")
//				mining(cur)
//				beego.Warn("mining over")
//				utils.RedisClient.Delete("Hodl")
//			}
//		}
//		if now.Sub(cur.ShareLast).Minutes() >= 60 {
//			if _, err := utils.RedisClient.GetString("Share"); err != nil {
//				utils.RedisClient.Set("Share", time.Now().Format("2006-01-02 15:04:05"), 10 * 60)
//				beego.Warn("start share")
//				share(cur)
//				beego.Warn("share over")
//				utils.RedisClient.Delete("Share")
//			}
//		}
//	}
//}
