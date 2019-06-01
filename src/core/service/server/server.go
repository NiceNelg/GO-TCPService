package server

import (
	"core/config"
	clientconn "core/service/client-conn"
	"core/service/controller"
	"fmt"
	"net"
	"os"
	"tool"
)

/**
 * @Struct 服务端
 * @Auther Nelg
 * @Date 2019.05.31
 */
type Server struct {
	/*监听的ip地址*/
	host string
	/*监听的端口*/
	port string
	/*所有配置信息*/
	config *config.Config
	/*处理者*/
	handler controller.Handler
	/*接收者*/
	receiver controller.Receiver
	/*发送者*/
	sender controller.Sender
}

/**
 * @Function 初始化服务
 * @Auther Nelg
 * @Date 2019.05.31
 */
func Init(configure *config.Config, handler controller.Handler, receiver controller.Receiver,
	sender controller.Sender) (serv Server) {
	serv = Server{
		host:     configure.Service.Host,
		port:     configure.Service.Port,
		config:   configure,
		handler:  handler,
		sender:   sender,
		receiver: receiver,
	}
	return
}

/**
 * @Function 开启tcp服务
 * @Auther Nelg
 * @Date 2019.05.30
 */
func (this *Server) StartTCPServer() {

	/*建立TCP服务器*/
	tcpAddr, err := net.ResolveTCPAddr("tcp", this.host+":"+this.port)
	if err != nil {
		fmt.Println("Start server error：", err)
		os.Exit(-1)
	}

	/*监听端口*/
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	for {
		/*接受连接请求*/
		tcpConn, _ := tcpListener.AcceptTCP()
		/*实例化客户端*/
		cli := clientconn.ClientConn{
			Conn:      tcpConn,
			Hearttime: this.config.Client.HeartTimeOut,
			LogConfig: this.config.Log,
		}

		if this.config.Log.Switch > 0 {
			fd, _ := cli.Conn.File()
			tool.WriteLog(this.config.Log.StdPath, fd.Name(), "[client_connect]")
		}
		/*新建设备协程，用于接收客户端数据，发送客户端数据*/
		go cli.DeviceCoroutines(this.handler, this.receiver, this.sender)
	}
}
