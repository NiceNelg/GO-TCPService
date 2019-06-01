package config

import (
	"github.com/aWildProgrammer/fconf"
	"log"
)

/**
 * @Struct 数据库配置信息
 * @Auther Nelg
 * @Date 2019.05.30
 */
type DatabaseConfig struct {
	Host     string
	Port     string
	Password string
}

/**
 * @Struct Mysql配置信息
 * @Auther Nelg
 * @Date 2019.05.30
 */
type MysqlConfig struct {
	DatabaseConfig
	Username string
	Database string
}

/**
 * @Struct Redis配置信息
 * @Auther Nelg
 * @Date 2019.05.30
 */
type RedisConfig struct {
	DatabaseConfig
	Database string
}

/**
 * @Struct 日志置信息
 * @Auther Nelg
 * @Date 2019.05.30
 */
type LogConfig struct {
	/*日志开关*/
	Switch int
	/*客户端日志路径*/
	StdPath string
	/*错误日志路径*/
	ErrPath string
}

/**
 * @Struct 队列配置信息
 * @Auther Nelg
 * @Date 2019.05.30
 */
type ListConfig struct {
	/*处理队列key名*/
	Handle string
	/*发送队列key名*/
	Send string
	/*重发队列key名*/
	Resend string
}

/**
 * @Struct 服务配置信息
 * @Auther Nelg
 * @Date 2019.05.30
 */
type ServiceConfig struct {
	/*监听IP*/
	Host string
	/*监听端口号*/
	Port string
}

/**
 * @Struct 客户端配置信息
 * @Auther Nelg
 * @Date 2019.05.30
 */
type ClientConfig struct {
	/*心跳超时断开连接秒数*/
	HeartTimeOut int
}

/**
 * @Struct 处理者配置信息
 * @Auther Nelg
 * @Date 2019.05.30
 */
type HandlerConfig struct {
	/*处理数据协程数*/
	Worknum int
	/*重发时间间隔*/
	Resendtime int64
	/*重发次数*/
	Resendnum int
}

/**
 * @Struct 硬件服务器配置信息
 * @Auther Nelg
 * @Date 2019.05.30
 */
type Config struct {
	Log     LogConfig
	Redis   RedisConfig
	Mysql   MysqlConfig
	List    ListConfig
	Service ServiceConfig
	Client  ClientConfig
	Handler HandlerConfig
}

/**
 * @Function 获取配置
 * @Auther Nelg
 * @Date 2019.05.30
 */
func GetConfig() (configure Config) {

	/*读取配置文件*/
	configFile, err := fconf.NewFileConf("/home/wwwroot/liankong_device/device.ini")
	if err != nil {
		log.Fatalln("read configure file failr\n")
	}

	/*读取日志相关配置*/
	configure.Log.Switch, _ = configFile.Int("log.switch")
	configure.Log.StdPath = configFile.String("log.std_log_path")
	configure.Log.ErrPath = configFile.String("log.error_log_path")

	/*读取redis相关配置*/
	configure.Redis.Host = configFile.String("redis.host")
	configure.Redis.Port = configFile.String("redis.port")
	configure.Redis.Password = configFile.String("redis.password")
	configure.Redis.Database = configFile.String("redis.db")

	/*读取mysql相关配置*/
	configure.Mysql.Host = configFile.String("mysql.host")
	configure.Mysql.Port = configFile.String("mysql.port")
	configure.Mysql.Username = configFile.String("mysql.username")
	configure.Mysql.Password = configFile.String("mysql.password")
	configure.Mysql.Database = configFile.String("mysql.database")

	/*读取队列相关配置*/
	configure.List.Handle = configFile.String("list.handle")
	configure.List.Send = configFile.String("list.send")
	configure.List.Resend = configFile.String("list.resend")

	/*读取服务相关配置*/
	configure.Service.Host = configFile.String("service.host")
	configure.Service.Port = configFile.String("service.port")

	/*读取客户端相关配置*/
	configure.Client.HeartTimeOut, _ = configFile.Int("client.heart_time_out")

	/*读取处理者相关配置*/
	configure.Handler.Worknum, _ = configFile.Int("handler.worknum")
	configure.Handler.Resendtime, _ = configFile.Int64("handler.resendtime")
	configure.Handler.Resendnum, err = configFile.Int("handler.resendnum")
	return
}
