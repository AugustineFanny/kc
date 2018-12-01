package controllers
//
//import (
//	"kuangchi_backend/models"
//	"kuangchi_backend/result"
//	"kuangchi_backend/job"
//	"encoding/json"
//	"github.com/spf13/cast"
//	"strings"
//	"path"
//	"encoding/csv"
//	"io"
//	"time"
//)
//
//// @router /simulation/users [post]
//func (u *AdminController) CreateSUer() {
//	var form struct {
//		Username    string  `json:"username"`
//		Parent      string  `json:"parent"`
//	}
//	json.Unmarshal(u.Ctx.Input.RequestBody, &form)
//	if err := models.CreateSUser(form.Username, form.Parent); err != nil {
//		u.Error(err)
//		return
//	}
//	u.Ok()
//}
//
//// @router /simulation/batch/users [post]
//func (u *AdminController) BatchSUsers() {
//	f, h, err := u.GetFile("file")
//	if err != nil {
//		u.Error("缺少文件")
//		return
//	}
//	defer f.Close()
//	ext := strings.ToLower(path.Ext(h.Filename))
//	if !(ext == ".csv") {
//		u.Error(100310)
//		return
//	}
//	reader := csv.NewReader(f)
//	for {
//		record, err := reader.Read()
//		if err == io.EOF {
//			break
//		} else if err != nil {
//			u.Error(err)
//			return
//		}
//		if len(record) == 1 {
//			if err := models.CreateSUser(record[0], ""); err != nil {
//				u.Error(err)
//				return
//			}
//		}
//		if len(record) >= 2 {
//			if err := models.CreateSUser(record[0], record[1]); err != nil {
//				u.Error(err)
//				return
//			}
//		}
//	}
//	u.Ok()
//}
//
//// @router /simulation/users [get]
//func (u *AdminController) GetSUsers() {
//	u.Data["json"] = result.Success(models.GetSUsers())
//	u.ServeJSON()
//}
//
//// @router /simulation/users [delete]
//func (u *AdminController) DeleteAllSUsers() {
//	if err := models.DeleteAllSUsers(); err != nil {
//		u.Error(err)
//		return
//	}
//	u.Ok()
//}
//
//// @router /simulation/user/:uid [post]
//func (u *AdminController) ChangeSUser() {
//	var form struct {
//		Parent      string  `json:"parent"`
//	}
//	json.Unmarshal(u.Ctx.Input.RequestBody, &form)
//	var err error
//	uid, err := u.GetInt(":uid")
//	if err != nil {
//		u.Error(err)
//		return
//	}
//	if err := models.ChangeSUser(uid, form.Parent); err != nil {
//		u.Error(err)
//		return
//	}
//	u.Ok()
//}
//
//// @router /simulation/user/:uid [delete]
//func (u *AdminController) DeleteSUser() {
//	uid, err := u.GetInt(":uid")
//	if err != nil {
//		u.Error(err)
//		return
//	}
//	if err := models.DeleteSUser(uid); err != nil {
//		u.Error(err)
//		return
//	}
//	u.Ok()
//}
//
//// @router /simulation/activities [post]
//func (u *AdminController) CreateSActivity() {
//	var form struct {
//		Username     string    `json:"username"`
//		Subscription float64   `json:"subscription"`
//		Lock         float64   `json:"lock"`
//		In           float64   `json:"in"`
//		InSource     string    `json:"in_source"`
//		Out          float64   `json:"out"`
//		OutDest      string    `json:"out_dest"`
//		Date         string    `json:"date"`
//	}
//	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &form); err != nil {
//		u.Error(err)
//		return
//	}
//	activity := models.KcSActivity{
//		Username: form.Username,
//		Subscription: form.Subscription,
//		Lock: form.Lock,
//		In: form.In,
//		InSource: form.InSource,
//		Out: form.Out,
//		OutDest: form.OutDest,
//		Date: cast.ToTime(form.Date),
//	}
//	if err := models.CreateSActivity(activity); err != nil {
//		u.Error(err)
//		return
//	}
//	u.Ok()
//}
//
//// @router /simulation/batch/activities [post]
//func (u *AdminController) BatchSActivities() {
//	f, h, err := u.GetFile("file")
//	if err != nil {
//		u.Error("缺少文件")
//		return
//	}
//	defer f.Close()
//	ext := strings.ToLower(path.Ext(h.Filename))
//	if !(ext == ".csv") {
//		u.Error(100310)
//		return
//	}
//	reader := csv.NewReader(f)
//	for {
//		record, err := reader.Read()
//		if err == io.EOF {
//			break
//		} else if err != nil {
//			u.Error(err)
//			return
//		}
//		if len(record) != 8 {
//			u.Error("必须八列：用户名，认购，转移，转入，转入来源用户，转出，转出目标用户，日期")
//			return
//		}
//		if record[0] == "" {
//			continue
//		}
//		date, err := time.Parse("2006/1/_2", record[7])
//		if err != nil {
//			u.Error("日期格式：2018/6/1")
//			return
//		}
//
//		activity := models.KcSActivity{
//			Username: 	  cast.ToString(record[0]),
//			Subscription: cast.ToFloat64(record[1]),
//			Lock: 		  cast.ToFloat64(record[2]),
//			In: 		  cast.ToFloat64(record[3]),
//			InSource: 	  cast.ToString(record[4]),
//			Out: 		  cast.ToFloat64(record[5]),
//			OutDest: 	  cast.ToString(record[6]),
//			Date: 		  date,
//		}
//		if err := models.CreateSActivity(activity); err != nil {
//			u.Error(err)
//			return
//		}
//	}
//	u.Ok()
//}
//
//// @router /simulation/activities [get]
//func (u *AdminController) GetSActivities() {
//	u.Data["json"] = result.Success(models.GetSActivities())
//	u.ServeJSON()
//}
//
//// @router /simulation/activities [delete]
//func (u *AdminController) DeleteAllSActivities() {
//	if err := models.DeleteAllSActivities(); err != nil {
//		u.Error(err)
//		return
//	}
//	u.Ok()
//}
//
//// @router /simulation/activity/:id [post]
//func (u *AdminController) ChangeSActivity() {
//	var form struct {
//		Subscription float64   `json:"subscription"`
//		Lock         float64   `json:"lock"`
//		In           float64   `json:"in"`
//		InSource     string    `json:"in_source"`
//		Out          float64   `json:"out"`
//		OutDest    string      `json:"out_dest"`
//		Date         string    `json:"date"`
//	}
//	if err := json.Unmarshal(u.Ctx.Input.RequestBody, &form); err != nil {
//		u.Error(err)
//		return
//	}
//	var err error
//	id, err := u.GetInt(":id")
//	if err != nil {
//		u.Error(err)
//		return
//	}
//	activity, err := models.GetSActivity(id)
//	if err != nil {
//		u.Error("目标不存在")
//		return
//	}
//	activity.Subscription = form.Subscription
//	activity.Lock = form.Lock
//	activity.In = form.In
//	activity.InSource = form.InSource
//	activity.Out = form.Out
//	activity.OutDest = form.OutDest
//	activity.Date = cast.ToTime(form.Date)
//	if err := models.ChangeSActivity(*activity); err != nil {
//		u.Error(err)
//		return
//	}
//	u.Ok()
//}
//
//// @router /simulation/activity/:id [delete]
//func (u *AdminController) DeleteSActivity() {
//	id, err := u.GetInt(":id")
//	if err != nil {
//		u.Error(err)
//		return
//	}
//	if err := models.DeleteSActivity(id); err != nil {
//		u.Error(err)
//		return
//	}
//	u.Ok()
//}
//
//// @router /simulation/currency [get]
//func (u *AdminController) GetSCurrency() {
//	u.Data["json"] = result.Success(models.GetSCurrency())
//	u.ServeJSON()
//}
//
//// @router /simulation/currency [post]
//func (u *AdminController) ChangeSCurrency() {
//	var form models.KcSCurrency
//	json.Unmarshal(u.Ctx.Input.RequestBody, &form)
//	if err := models.ChangeSCurrency(form.MInterestRate, form.CInterestRate); err != nil {
//		u.Error(err)
//		return
//	}
//	u.Ok()
//}
//
//// @router /simulation/calculation [get]
//func (u *AdminController) Calculation() {
//	username := u.GetString("username")
//	startDateS := u.GetString("start_date")
//	endDateS := u.GetString("end_date")
//	res, err := job.Calculation(username, startDateS, endDateS)
//	if err != nil {
//		u.Error(result.ApiResponse{100102, err.Error(), res})
//		return
//	}
//	u.Data["json"] = result.Success(res)
//	u.ServeJSON()
//}
//
//// @router /simulation/stat [get]
//func (u *AdminController) SimulationStat() {
//	startDateS := u.GetString("start_date")
//	endDateS := u.GetString("end_date")
//	res, err := job.Stat(startDateS, endDateS)
//	if err != nil {
//		u.Error(result.ApiResponse{100102, err.Error(), res})
//		return
//	}
//	u.Data["json"] = result.Success(res)
//	u.ServeJSON()
//}
