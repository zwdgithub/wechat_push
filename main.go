package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/toolbox"
	"log"
	"wechat_push/controllers"
	"wechat_push/models"
	_ "wechat_push/routers"
)

func main() {
	controllers.UserInit()
	tk1 := toolbox.NewTask("tk1", "0 55 * * * *", func() error { models.RefreshAccessToken(); return nil })
	err := tk1.Run()
	if err != nil {
		log.Fatalln("task run failed ...")
	}
	beego.Run()
	log.Println("server start success...")
}
