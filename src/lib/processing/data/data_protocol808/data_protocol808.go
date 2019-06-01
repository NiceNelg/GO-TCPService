package data_protocol808

import "lib/processing/data/data"

/**
 * @Struct 数据头部
 * @Auther Nelg
 * @Date 2019.05.30
 */
type DataHeader struct {
	/*命令标识*/
	Sign string `json:"sign"`
	/*消息属性*/
	Attribute string `json:"attribute"`
	/*命令流水*/
	Sn string `json:"sn"`
}

/**
 * @Struct 数据包
 * @Auther Nelg
 * @Date 2019.05.30
 */
type Data struct {
	/*数据头部*/
	DataHeader
	data.Data
}
