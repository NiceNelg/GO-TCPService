package handleunit

import (
	handleinterface "core/service/handle-interface"
	"database/sql"
	"lib/processing/data/data_protocol808"
	basehandleunit "lib/processing/processing-unit/protocol808/base-handleunit"
	"time"
)

type Authcheck struct {
	basehandleunit.HandleUnit
}

/**
 * @Function 数据内容解析
 * @Auther Nelg
 * @Date 2019.05.30
 */
func AuthcheckInit(cmd data_protocol808.Data) (authcheck *Authcheck) {
	content := []rune(cmd.Content)
	cmd.Body = make(map[string]string, 6)
	cmd.Body["authkey"] = string(content[:8])
	cmd.Body["std_type"] = string(content[8:16])
	cmd.Body["car_type"] = string(content[16:18])
	cmd.Body["CCID"] = string(content[18:58])
	cmd.Body["now_version_len"] = string(content[58:60])
	cmd.Body["now_version"] = string(content[60:])
	authcheck = new(Authcheck)
	authcheck.Data = cmd
	return
}

/**
 * @Function 业务处理	需要数据下发的时候请返回HandleInterface
 * @Auther Nelg
 * @Date 2019.05.30
 */
func (this *Authcheck) HandleBusiness(db *sql.DB) (sendCmd handleinterface.HandleInterface) {
	authCheck := &Authcheck{}
	authCheck.Sign = "8102"
	authCheck.Sn = this.Sn
	authCheck.ClientNumber = this.ClientNumber
	authCheck.Body = make(map[string]string, 4)
	authCheck.Body["ack_sn"] = this.Sn
	authCheck.Body["ack_sign"] = this.Sign
	authCheck.Body["result"] = "00"
	/*获取当前时间*/
	authCheck.Body["time"] = time.Now().Format("060102150405")

	//保存指令
	//commandModel := model.CommandModelInit(db)
	//id := commandModel.SaveCommand()
	//sn := strconv.FormatInt(id, 16)
	//this.Sn = tool.StrPad(sn, "0", 4, "LEFT")
	////更新指令流水号
	//commandModel.SaveSn(id, this.Sn)
	sendCmd = authCheck
	return
}
