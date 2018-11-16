// package controllers

// import (
// 	"kuangchi_backend/job"
// 	"kuangchi_backend/result"
// 	"encoding/json"
// )

// // @router /public/btcdo-script [get]
// func (u *PublicController) GetBtcdoScript() {
// 	u.Data["json"] = result.Success(map[string]interface{}{
// 		"symbol": job.Symbol,
// 		"amount": job.Amount,
// 		"floatation": job.Floatation,
// 		"sellMinAmount": job.SellMinAmount,
// 		"buyMinAmount": job.BuyMinAmount,
// 		"rest": float64(job.Rest) / 1000,
// 		"startFlag": job.StartFlag,
// 	})
// 	u.ServeJSON()
// }

// // @router /public/btcdo-script-logs [get]
// func (u *PublicController) GetBtcdoScriptLogs() {
// 	u.Data["json"] = result.Success(job.TradeLogs.Pop())
// 	u.ServeJSON()
// }

// // @router /public/btcdo-script [post]
// func (u *PublicController) BtcdoScript() {
// 	var form struct {
// 		Symbol        string
// 		Amount        float64
// 		Floatation    float64
// 		SellMinAmount float64
// 		BuyMinAmount  float64
// 		Rest          float64
// 	}
// 	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &form); err != nil {
// 		u.Error("参数 有误")
// 		return
// 	}
// 	if form.Symbol != "ETH_USDT" &&
// 		form.Symbol != "BDB_USDT" &&
// 		form.Symbol != "IOST_BDB" &&
// 		form.Symbol != "BDB_ETH" {
// 		u.Error("交易对 有误")
// 		return
// 	}
// 	if form.Amount <= 0 {
// 		u.Error("数量 必须大于0")
// 		return
// 	}
// 	if form.Floatation <= 0 {
// 		u.Error("价格 必须大于0")
// 		return
// 	}
// 	if form.SellMinAmount < 0 {
// 		u.Error("卖单数量限制 必须大于等于0")
// 		return
// 	}
// 	if form.BuyMinAmount < 0 {
// 		u.Error("买单数量限制 必须大于等于0")
// 		return
// 	}
// 	if form.Rest < 0.2 {
// 		u.Error("间隔时间 必须大于0.2秒")
// 		return
// 	}
// 	job.Symbol = form.Symbol
// 	job.Amount = form.Amount
// 	job.Floatation = form.Floatation
// 	job.SellMinAmount = form.SellMinAmount
// 	job.BuyMinAmount = form.BuyMinAmount
// 	job.Rest = int(form.Rest * 1000)
// 	u.Ok()
// }

// // @router /public/btcdo-script-start [post]
// func (u *PublicController) BtcdoScriptStart() {
// 	var form struct {
// 		ApiKey        string
// 		ApiSecret     string
// 		Passwd        string
// 	}
// 	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &form); err != nil {
// 		u.Error("参数 有误")
// 		return
// 	}
// 	if form.Passwd != "5b42b754002ac907d9e5c325" {
// 		u.Error("密码 有误")
// 		return
// 	}
// 	job.API_KEY = form.ApiKey
// 	job.API_SECRET = form.ApiSecret
// 	job.StartFlag = true
// 	job.TradeLogs.Push("start")
// 	u.Ok()
// }

// // @router /public/btcdo-script-stop [post]
// func (u *PublicController) BtcdoScriptStop() {
// 	job.StartFlag = false
// 	job.TradeLogs.Push("stop")
// 	u.Ok()
// }
