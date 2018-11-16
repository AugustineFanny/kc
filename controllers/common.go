package controllers

import (
	"kuangchi_backend/models"
	"kuangchi_backend/result"
	"kuangchi_backend/utils"
	"github.com/astaxie/beego"
	"path"
	"strconv"
	"io/ioutil"
	"strings"
	"os"
	"io"
)

type CommonController struct {
	beego.Controller
}

func (c *CommonController) Error(err interface{}) {
	switch err.(type) {
	case result.ApiResponse:
		c.Data["json"] = err
		c.ServeJSON()
	case error:
		c.Data["json"] = result.ErrMsg(err.(error).Error())
		c.ServeJSON()
	case string:
		c.Data["json"] = result.ErrMsg(err.(string))
		c.ServeJSON()
	case int:
		c.Data["json"] = result.ErrCode(err.(int))
		c.ServeJSON()
	default:
		c.Data["json"] = result.ErrCode(9999)
		c.ServeJSON()
	}
}

func (c *CommonController) Ok() {
	c.Data["json"] = result.Success("")
	c.ServeJSON()
}

func (c *CommonController) GetUser() *models.User {
	authorization := c.Ctx.Request.Header.Get("Authorization")
	username := utils.CheckToken(authorization)
	if username != "" {
		return models.GetUserByUsername(username)
	}
	return nil
}

func (c *CommonController) GetAdmin() *models.Admin {
	v := c.GetSession("admin")
	if v != nil {
		return v.(*models.Admin)
	}
	return nil
}

func (c *CommonController) SaveImg(id int, dir, field string, required bool) (string, error) {
	f, h, err := c.GetFile(field)
	if err != nil {
		if required == true {
			return "", result.ErrMsg("missing " + field)
		}
		return "", nil
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return "", result.ErrCode(100102)
	}
	if len(buf) > 1024 * 1024 * 2 {
		return "", result.ErrCode(100111)
	}
	ext := strings.ToLower(path.Ext(h.Filename))
	if !(ext == ".jpg" || ext == ".jpeg" || ext == ".png") {
		return "", result.ErrCode(100310)
	}
	filename := field + "_" + strconv.Itoa(id) + "_" + utils.TimestampString() + ext
	if err := c.SaveToFile(field, path.Join(dir, filename)); err != nil {
		beego.Error(err)
		return "", result.ErrCode(100102)
	}
	return filename, nil
}

func (c *CommonController) SaveImgs(id int, dir, field string, required bool) (string, error) {
	files, err:=c.GetFiles(field)
	if err != nil {
		if required == true {
			return "", result.ErrMsg("missing " + field)
		}
		return "", nil
	}
	res := []string{}
	for index, f := range files {
		file, err := f.Open()
		defer file.Close()
		if err != nil {
			return "", err
		}
		ext := strings.ToLower(path.Ext(f.Filename))
		if !(ext == ".jpg" || ext == ".jpeg" || ext == ".png") {
			return "", result.ErrCode(100310)
		}
		filename := field + "_" + strconv.Itoa(index) + "_" + utils.TimestampString() + ext
		dst, err := os.Create(path.Join(dir, filename))
		defer dst.Close()
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(dst, file); err != nil {
			return "", err
		}
		res = append(res, filename)
	}
	return strings.Join(res, ","), nil
}