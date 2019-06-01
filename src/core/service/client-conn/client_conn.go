package client_conn

import (
	"core/config"
	"core/service/controller"
	"lib/processing/processing"
	"net"
	"time"
	"tool"
)

/**
 * @Struct 客户端连接
 * @Auther Nelg
 * @Date 2019.05.31
 */
type ClientConn struct {
	/*客户端连接套接字*/
	Conn *net.TCPConn
	/*客户端编号*/
	number string
	/*未完整的数据包*/
	buffer []byte
	/*心跳时间*/
	Hearttime int
	/*日志配置信息*/
	LogConfig config.LogConfig
}

/**
 * @Function 设备协程
 * @Auther Nelg
 * @Date 2019.05.30
 */
func (this *ClientConn) DeviceCoroutines(handler controller.Handler, receiver controller.Receiver,
	sender controller.Sender) {
	/*协程结束后关闭连接*/
	defer func() {
		err := this.Conn.Close()
		if err != nil && this.LogConfig.Switch > 0 {
			tool.WriteLog(this.LogConfig.ErrPath, err.Error(), "[client_close_error]")
		} else if this.LogConfig.Switch > 0 {
			fd, _ := this.Conn.File()
			tool.WriteLog(this.LogConfig.StdPath, fd.Name(), "[client_close]")
		}
	}()
	/*开启发送协程*/
	go sender.Send(&this.number, this.Conn)
	/*接收数据*/
	for {
		/*设置连接存活时间*/
		err := this.Conn.SetReadDeadline(time.Now().Add(time.Duration(this.Hearttime) * time.Second))
		if err != nil {
			break
		}

		/*读取数据*/
		readData := make([]byte, 2048)
		length, err := this.Conn.Read(readData)
		if err != nil {
			break
		} else if length <= 0 {
			continue
		}
		if this.LogConfig.Switch > 0 {
			tool.WriteLog(this.LogConfig.StdPath, string(readData), "[receive_pkg]")
		}

		/*包数据分割*/
		var dataArray [][]byte
		dataArray, this.buffer = processing.Cutpack(readData[0:length], this.buffer)

		/*存入处理队列*/
		go receiver.SaveTask(&this.number, dataArray)
	}

}
