package handle_interface

import (
	dataunit "core/service/data-unit"
	"database/sql"
)

/**
 * @Interface 对数据进行操作的接口
 * @Auther Nelg
 * @Date 2019.05.30
 */
type HandleInterface interface {
	/*处理数据包业务*/
	HandleBusiness(*sql.DB) HandleInterface
	/*生成发送的数据包内容*/
	HandleSend() dataunit.SendUnit
}
