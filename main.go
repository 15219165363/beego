package main

import (
	_ "quickstart/routers"
	"quickstart/proxy"
	"github.com/astaxie/beego"
)


func main() {

	//设置日志文件
	beego.SetLogger("file", `{"filename":"/zzw/test/beego.log"}`)

	//设置日志等级
	level, err := beego.AppConfig.Int("log.level")
	if err != nil {
		beego.Error("Log Level Config Not Found!")
		return
	}
	beego.SetLevel(level)

	//期望日志输出调用的文件名和文件行号
	beego.SetLogFuncCall(true)	

	pport := beego.AppConfig.DefaultString("pport", "8088")
	err = proxy.StartServer(pport)
	if err != nil {
		beego.Error("Failed to Start Proxy  Server!")
		return
	}	
	beego.Run()
}

