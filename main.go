package main

import (
	"github.com/astaxie/beego"
	_ "sensitive_word/routers"
	_ "sensitive_word/service"
)

func main() {
	beego.Run()
}
