package main

import (
	"fmt"
	"os"
	"time"

	"github.com/betacraft/yaag/yaag"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/snowlyg/IrisAdminApi/config"
	"github.com/snowlyg/IrisAdminApi/database"
	"github.com/snowlyg/IrisAdminApi/files"
	"github.com/snowlyg/IrisAdminApi/models"
	"github.com/snowlyg/IrisAdminApi/routes"
)

func NewLogFile() *os.File {
	path := "./logs/"
	_ = files.CreateFile(path)
	filename := path + time.Now().Format("2006-01-02") + ".log"
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		color.Red(fmt.Sprintf("日志记录出错: %v", err))
	}

	return f
}

func NewApp() *iris.Application {
	api := iris.New()
	api.Logger().SetLevel(config.GetAppLoggerLevel())

	api.RegisterView(iris.HTML("resources", ".html"))

	db := database.GetGdb()
	db.AutoMigrate(
		&models.User{},
		&models.OauthToken{},
		&models.Role{},
		&models.Permission{},
	)

	iris.RegisterOnInterrupt(func() {
		_ = db.Close()
	})

	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware. //api 文档配置
		On:       true,
		DocTitle: config.GetAppName(),
		DocPath:  "./resources/apiDoc/index.html", //设置绝对路径
		BaseUrls: map[string]string{
			"Production": config.GetAppUrl(),
			"Staging":    "",
		},
	})

	routes.App(api) //注册 app 路由

	return api
}
