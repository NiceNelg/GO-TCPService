package main

import (
	"core/config"
	"core/mysqlpool"
	"core/redispool"
	"core/service/controller"
	"core/service/server"
)

func main() {

	/*获取配置项*/
	configure := config.GetConfig()

	/*初始化redis*/
	redisPool := redispool.NewPool(
		configure.Redis.Host+":"+configure.Redis.Port,
		configure.Redis.Password,
		configure.Redis.Database,
	)

	/*初始化数据库*/
	db := mysqlpool.Init(
		configure.Mysql.Username,
		configure.Mysql.Password,
		configure.Mysql.Host,
		configure.Mysql.Port,
		configure.Mysql.Database,
	)
	defer db.Close()

	/*实例化操作者*/
	handler := controller.InitHandler(db, redisPool, &configure)
	/*开启数据处理队列*/
	handler.StartHandler()

	/*实例化重发者*/
	resender := controller.InitResender(redisPool, &configure)
	/*开启数据处理队列*/
	resender.StartResend()

	/*实例化接收者*/
	receiver := controller.InitReceiver(redisPool, &configure)
	/*实例化发送者*/
	sender := controller.InitSender(redisPool, &configure)

	/*初始化服务*/
	serv := server.Init(&configure, handler, receiver, sender)

	/*开始TCP服务*/
	serv.StartTCPServer()

	/*开始Websocket服务*/
}
