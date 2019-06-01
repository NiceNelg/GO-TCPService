package base_handleunit

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	dataunit "core/service/data-unit"
	"github.com/garyburd/redigo/redis"
	"lib/processing/data/data_protocol808"
	"lib/processing/tool/tool_protocol808"
	"tool"
)

/**
 * @Struct 数据单元
 * @Auther Nelg
 * @Date 2019.05.30
 */
type HandleUnit struct {
	data_protocol808.Data
}

/**
 * @Function 数据包存入发送队列
 * @Auther Nelg
 * @Date 2019.05.30
 */
func (this *HandleUnit) SaveToSendList(redisPool *redis.Pool, sendList string) {
	cmd, _ := json.Marshal(this)
	//从redis中取出连接
	redisCli := redisPool.Get()
	//归还redis连接到redis连接池
	defer redisCli.Close()
	redisCli.Do("lpush", sendList+"_"+this.ClientNumber, string(cmd))
}

/**
 * @Function 组成发送数据
 * @Auther Nelg
 * @Date 2019.05.30
 */
func (this *HandleUnit) HandleSend() (send dataunit.SendUnit) {
	/*组成发送数据*/
	content := this.Body["ack_sign"] + this.Body["ack_sn"] + this.Body["result"]
	this.Attribute = tool.StrPad(strconv.Itoa(len(content)/2), "0", 4, "LEFT")
	data := this.Sign + this.Attribute + this.ClientNumber + this.Sn + content
	/*转换成[]byte*/
	dataByte, _ := hex.DecodeString(data)
	/*异或计算*/
	dataByte = append(dataByte, tool.BuildBCC(dataByte))
	/*转义*/
	dataByte = tool_protocol808.Escape(dataByte)
	this.Content = hex.EncodeToString(dataByte)
	/*组装发送单元*/
	send = dataunit.SendUnit{
		Content:  this.Content,
		Sendtime: 0,
		Sendnum:  0,
	}
	return
}
