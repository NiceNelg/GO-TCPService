package controller

import (
	"core/config"
	"encoding/json"
	"net"
	"time"
	"tool"

	dataunit "core/service/data-unit"
	"github.com/garyburd/redigo/redis"
	"lib/processing/processing"
)

/**
 * @Struct 发送配置
 * @Auther Nelg
 * @Date 2019.05.31
 */
type send struct {
	/*发送队列key*/
	list string
}

/**
 * @Struct 发送者
 * @Auther Nelg
 * @Date 2019.05.31
 */
type Sender struct {
	/*redis连接池*/
	redisPool *redis.Pool
	/*发送配置*/
	send send
	/*重发配置*/
	resend resend
	/*日志配置*/
	logconfig config.LogConfig
}

/**
 * @Function 初始化服务操作对象
 * @Auther Nelg
 * @Date 2019.05.30
 */
func InitSender(redisPool *redis.Pool, configure *config.Config) (sender Sender) {
	sender = Sender{
		redisPool: redisPool,
		send: send{
			list: configure.List.Send,
		},
		resend: resend{
			list:     configure.List.Resend,
			sendtime: configure.Handler.Resendtime,
			num:      configure.Handler.Resendnum,
		},
		logconfig: configure.Log,
	}
	return
}

/**
 * @Function 发送协程
 * @Auther Nelg
 * @Date 2019.05.30
 */
func (this *Sender) Send(clientNumber *string, conn *net.TCPConn) {
	for {
		/*客户端还没上传设备号*/
		if *clientNumber == "" {
			time.Sleep(time.Second)
			continue
		}
		failCmd := make([]string, 1)
		for {
			/*从redis连接池中取出连接*/
			redisCli := this.redisPool.Get()
			/*取出最新一条需要发送的信息*/
			redisData, err := redis.String(redisCli.Do("rpop", this.send.list+"_"+*clientNumber))
			/*归还redis连接到redis连接池*/
			redisCli.Close()
			if err != nil {
				break
			}

			/*还原发送数据单元对象*/
			var cmd dataunit.SendUnit
			err = json.Unmarshal([]byte(redisData), &cmd)
			if err != nil {
				continue
			}

			/*重组装发送数据*/
			pkg, err := processing.RebuildSendPackage(cmd.Content)
			if err != nil {
				continue
			}

			_, err = conn.Write(pkg)
			if err != nil {
				if this.logconfig.Switch > 0 {
					tool.WriteLog(this.logconfig.StdPath, redisData, "[resend_fail]")
					tool.WriteLog(this.logconfig.ErrPath, err.Error(), "[resend_fail]")
				}
				//发送失败，检测发送次数是否已超出
				cmd.Sendnum++
				if cmd.Sendnum <= this.resend.num {
					cmd.Sendtime = time.Now().Unix()
					addCmd, _ := json.Marshal(cmd)
					failCmd = append(failCmd, string(addCmd))
				}
				if this.logconfig.Switch > 0 {
					tool.WriteLog(this.logconfig.StdPath, redisData, "[resend_add]")
				}
				return
			}
			if this.logconfig.Switch > 0 {
				tool.WriteLog(this.logconfig.StdPath, redisData, "[send_success]")
			}
		}
		for _, sendUnit := range failCmd {
			if sendUnit == "" {
				continue
			}
			redisCli := this.redisPool.Get()
			redisCli.Do("lpush", this.resend.list, sendUnit)
			redisCli.Close()
		}
		time.Sleep(1 * time.Second)
	}
	return
}
