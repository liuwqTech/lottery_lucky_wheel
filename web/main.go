package main

import (
	"Lucky_Wheel/bootstrap"
	"Lucky_Wheel/web/routes"
	"fmt"
)

var port = 8080

func newApp() *bootstrap.Bootstrapper {
	// 初始化应用
	app := bootstrap.New("抽奖大转盘", "小新")
	app.Bootstrap()
	app.Configure(routes.Configure)
	return app
}

func main() {
	app := newApp()
	app.Listen(fmt.Sprintf(":%d", port))
}
