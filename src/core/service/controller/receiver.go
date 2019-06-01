package controller

import (
	"core/config"
	"github.com/garyburd/redigo/redis"
	"lib/processing/processing"
	"tool"
)

/**
 * @Struct 重发配置
 * @Auther Nelg
 * @Date 2019.05.31
 */
type receive struct {
	/*数据处理队列key*/
	list string
}

/**
 * @Struct 接收者
 * @Auther Nelg
 * @Date 2019.05.31
 */
type Receiver struct {
	/*redis连接池*/
	redisPool *redis.Pool
	/*数据接收配置*/
	receive receive
	/*日志配置*/
	logconfig config.LogConfig
}

/**
 * @Function 初始化服务操作对象
 * @Auther Nelg
 * @Date 2019.05.30
 */
func InitReceiver(redisPool *redis.Pool, configure *config.Config) (receiver Receiver) {
	receiver = Receiver{
		redisPool: redisPool,
		receive: receive{
			list: configure.List.Handle,
		},
		logconfig: configure.Log,
	}
	return
}

/**
 * @Function 保存任务
 * @Auther Nelg
 * @Date 2019.05.30
 */
func (this *Receiver) SaveTask(clientNumber *string, dataArray [][]byte) {
	if len(dataArray) <= 0 {
		return
	}
	redisCli := this.redisPool.Get()
	defer redisCli.Close()
	/*解析数据结构*/
	for _, cmd := range dataArray {
		data, cliNumber, err := processing.Resolvepack(cmd)

		if err != nil {
			continue
		}
		/*设置客户端编号*/
		if *clientNumber == "" {
			*clientNumber = cliNumber
		}
		/*将包数据存入redis*/
		redisCli.Do("lpush", this.receive.list, data)

		/*记录日志*/
		if this.logconfig.Switch > 0 {
			tool.WriteLog(this.logconfig.StdPath, data, "[save_task]")
		}
	}
}
