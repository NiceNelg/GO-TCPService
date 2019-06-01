package data_unit

/**
 * @Struct 发送数据单元
 * @Auther Nelg
 * @Date 2019.05.31
 */
type SendUnit struct {
	/*接收客户端编号*/
	ClientNumber string
	/*发送内容*/
	Content string
	/*上次发送时间*/
	Sendtime int64
	/*发送次数*/
	Sendnum int
}
