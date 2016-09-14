package routers

import (
	"github.com/astaxie/beego"
	"sensitive_word/controllers"
)

func init() {
	beego.Router("/sensitive", &controllers.MainController{}, "post:Sensitive")
}
