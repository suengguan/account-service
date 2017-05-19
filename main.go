package main

import (
	_ "app-service/account-service/routers"
	"app-service/account-service/service"

	daoApi "api/dao_service"
	"github.com/astaxie/beego"
)

func main() {
	var err error

	var cfg = beego.AppConfig
	daoApi.ActionDaoApi.Init(cfg.String("ActionDaoService"))
	daoApi.AlgorithmDaoApi.Init(cfg.String("AlgorithmDaoService"))
	daoApi.BussinessDaoApi.Init(cfg.String("BussinessDaoService"))
	daoApi.ResourceDaoApi.Init(cfg.String("ResourceDaoService"))
	daoApi.UserDaoApi.Init(cfg.String("UserDaoService"))

	var svc service.AccountService
	err = svc.CreateAdmin()
	if err != nil {
		beego.Debug("create admin failed! reason:", err)
		return
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
