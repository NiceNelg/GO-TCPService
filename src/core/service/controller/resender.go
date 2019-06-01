package controller

import (
	"core/config"
	"encoding/json"
	"time"
	"tool"

	dataunit "core/service/data-unit"
	"github.com/garyburd/redigo/redis"
)

/**
 * @Struct 重发配置
 * @Auther Nelg
 * @Date 2019.05.31
 */
type resend struct {
	/*重发队列key*/
	list string
	/*重发时间*/
	sendtime int64
	/*重发时间间隔*/
	num int
}

/**
 * @Struct 重发者
 * @Auther Nelg
 * @Date 2019.05.31
 */
type Resender struct {
	/*redis连接池*/
	redisPool *redis.Pool
	/*重发配置*/
	resend resend
	/*发送配置*/
	send send
	/*日志配置*/
	logconfig config.LogConfig
}

/**
 * @Function 初始化服务操作对象
 * @Auther Nelg
 * @Date 2019.05.30
 */
func InitResender(redisPool *redis.Pool, configure *config.Config) (resender Resender) {
	resender = Resender{
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
 * @Function 重发协程
 * @Auther Nelg
 * @Date 2019.05.30
 */
func (this *Resender) StartResend() {
	for {
		failCmd := make([]string, 1)
		for {
			/*从redis连接池中取出连接*/
			redisCli := this.redisPool.Get()
			/*取出最新一条需要重新发送的信息*/
			redisData, err := redis.String(redisCli.Do("rpop", this.resend.list))
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

			/*检查发送次数、发送间隔是否符合*/
			if cmd.Sendnum <= this.resend.num && time.Now().Unix()-cmd.Sendtime >= this.resend.sendtime {
				redisCli := this.redisPool.Get()
				redisCli.Do("lpush", this.send.list+"_"+cmd.ClientNumber, redisData)
				redisCli.Close()
				if this.logconfig.Switch > 0 {
					tool.WriteLog(this.logconfig.StdPath, redisData, "[resend_pop]")
				}
				continue
			}

			failCmd = append(failCmd, redisData)
		}
		for _, sendUnit := range failCmd {
			if sendUnit == "" {
				continue
			}
			redisCli := this.redisPool.Get()
			redisCli.Do("lpush", this.resend.list, sendUnit)
			redisCli.Close()
		}
		time.Sleep(5 * time.Second)
	}
	return
}
