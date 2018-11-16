package controllers

import (
	"kuangchi_backend/result"
	"kuangchi_backend/models"
)

type UserController struct {
	CommonController
}

// @router /user/:username [get]
func (u *UserController) Info() {
	username := u.GetString(":username")
	user, err := models.GetUserDetail(username)
	if err != nil {
		u.Error(err)
		return
	}
	u.Data["json"] = result.Success(user)
	u.ServeJSON()
}
