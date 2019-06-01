package data

/**
 * @Struct 数据包
 * @Auther Nelg
 * @Date 2019.05.30
 */
type Data struct {
	/*客户端编号*/
	ClientNumber string `json:"client_number"`
	/*数据体（已解析）*/
	Body map[string]string `json:"-"`
	/*未解包数据（对于16进制数据需转换成可视化数据）*/
	Content string `json:"content"`
}
