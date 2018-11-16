package routers

import (
	"kuangchi_backend/models"
	"kuangchi_backend/result"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"kuangchi_backend/utils"
	"encoding/json"
)

func auth(ctx *context.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	username := utils.CheckToken(authorization)
	if username == "" {
		ctx.Output.JSON(result.ErrCode(100103), false, false)
		return
	}
	user := models.GetUserByUsername(username)

	if user == nil || user.Status == 0 {
		ctx.Output.JSON(result.ErrCode(100317), false, false)
	}

}

func twoFactor(ctx *context.Context) {
	authorization := ctx.Request.Header.Get("Authorization")
	username := utils.CheckToken(authorization)
	if username == "" {
		ctx.Output.JSON(result.ErrCode(100103), false, false)
		return
	}
	user := models.GetUserByUsername(username)

	if user.TfOpened == true {
		var form models.TfCodeDto
		json.Unmarshal(ctx.Input.RequestBody, &form)
		if !utils.DefaultOptions.Authenticate(user.TfSecret, form.TfCode) {
			ctx.Output.JSON(result.ErrCode(100108), false, false)
		}
	}
}

func admin(ctx *context.Context) {
	if ctx.Input.Session("admin") == nil {
		ctx.Output.JSON(result.ErrCode(100103), false, false)
	}
}

func superAdmin(ctx *context.Context) {
	v := ctx.Input.Session("admin")
	admin := v.(*models.Admin)
	if admin.Super != 1 {
		ctx.Output.JSON(result.ErrCode(100107), false, false)
	}
}

func init() {
	beego.InsertFilter("/kc/self/*", beego.BeforeRouter, auth)
	beego.InsertFilter("/kc/wallet/*", beego.BeforeRouter, auth)

	//beego.InsertFilter("/kc/ws", beego.BeforeRouter, wsAuth)

	beego.InsertFilter("/kc/admin/*", beego.BeforeRouter, admin)

	beego.InsertFilter("/kc/admin/super/*", beego.BeforeRouter, superAdmin)
}
