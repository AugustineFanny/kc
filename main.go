package main

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/session/redis"
	_ "github.com/go-sql-driver/mysql"

	"kuangchi_backend/job"
	_ "kuangchi_backend/routers"
)

func init() {
	job.Start()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
