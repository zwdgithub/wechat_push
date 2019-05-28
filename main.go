package main

import (
	"github.com/astaxie/beego"
	"log"
	"wechat_push/controllers"
	_ "wechat_push/routers"
)

func main() {
	controllers.UserInit()
	beego.Run()
	log.Println("server start success...")
}
