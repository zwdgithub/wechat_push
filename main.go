package main

import (
	"github.com/astaxie/beego"
	"log"
	_ "wechat_push/routers"
)

func main() {
	beego.Run()
	log.Println("server start success...")
}
