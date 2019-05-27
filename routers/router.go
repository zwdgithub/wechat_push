package routers

import (
	"github.com/astaxie/beego"
	"wechat_push/controllers"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	beego.Router("/message", &controllers.MainController{}, "get:CheckSignature")
	beego.Router("/message", &controllers.MainController{}, "post:Message")
}
