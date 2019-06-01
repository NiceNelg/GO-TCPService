package controller

import (
	"core/config"
	"database/sql"
	"lib/processing/processing"
	"time"
	"tool"

	"github.com/garyburd/redigo/redis"
)

/**
 * @Struct 服务操作者
 * @Auther Nelg
 * @Date 2019.05.30
 */
type Handler struct {
	/*mysql连接池*/
	db *sql.DB
	/*redis连接池*/
	redisPool *redis.Pool
	/*处理数据队列key*/
	handleList string
	/*处理者协程数*/
	worknum int
	/*发送队列key*/
	sendList string
	/*日志信息*/
	logConfig config.LogConfig
}

/**
 * @Function 初始化服务操作对象
 * @Auther Nelg
 * @Date 2019.05.30
 */
func InitHandler(db *sql.DB, redisPool *redis.Pool, configure *config.Config) (handler Handler) {
	handler = Handler{
		redisPool:  redisPool,
		db:         db,
		handleList: configure.List.Handle,
		worknum:    configure.Handler.Worknum,
		sendList:   configure.List.Send,
		logConfig:  configure.Log,
	}

	return
}

/**
 * @Function 开启数据处理协程
 * @Auther Nelg
 * @Date 2019.05.30
 */
func (this *Handler) StartHandler() {
	for i := 0; i < this.worknum; i++ {
		go this.invoke()
	}
	return
}

/**
 * @Function 处理业务
 * @Auther Nelg
 * @Date 2019.05.30
 */
func (this *Handler) invoke() {
	for {
		/*从redis中取出数据*/
		redisCli := this.redisPool.Get()
		redisData, err := redis.String(redisCli.Do("rpop", this.handleList))
		/*归还redis连接到redis连接池*/
		redisCli.Close()
		if err != nil || redisData == "" {
			time.Sleep(1 * time.Second)
			continue
		}

		if this.logConfig.Switch > 0 {
			tool.WriteLog(this.logConfig.StdPath, redisData, "[handle_pkg]")
		}
		/*获取业务对象，单独提到业务库中方便扩展*/
		unit, err := processing.GetBusinessObject(redisData)
		if err != nil {
			continue
		}
		/*处理业务*/
		unit = unit.HandleBusiness(this.db)
		/*若没有返回处理对象则代表此数据包不需要下发*/
		if unit == nil {
			continue
		}
		/*组成发送数据*/
		send := unit.HandleSend()
		if this.logConfig.Switch > 0 {
			tool.WriteLog(this.logConfig.StdPath, redisData, "[send_add]")
		}

		/*存入redis队列*/
		redisCli = this.redisPool.Get()
		_, err = redisCli.Do("lpush", this.sendList+"_"+send.ClientNumber, send.Content)
		redisCli.Close()
		if err != nil {
			continue
		}
	}
	return
}
