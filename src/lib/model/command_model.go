package model

import (
	"database/sql"
	"encoding/json"

	"lib/processing/data"
)

type Command struct {
	Model
}

/**
 * @Function 初始化指令Model
 * @Auther Nelg
 */
func CommandModelInit(db *sql.DB) (command Command) {
	command.db = new(sql.DB)
	command.db = db
	command.table = "tz_std_command"
	return
}

/**
 * @Function 保存指令
 * @Auther Nelg
 */
func (this *Command) SaveCommand(cmd data.Data) (id int64) {
	stmt, _ := this.db.Prepare("INSERT INTO " + this.table + "(stdid, commandId, sendcontent) VALUES (?, ?, ?)")
	sendcontent, _ := json.Marshal(cmd.Body)
	res, _ := stmt.Exec(cmd.Device, cmd.Sign, string(sendcontent))
	id, _ = res.LastInsertId()
	return
}

/**
 * @Function 保存流水号
 * @Auther Nelg
 */
func (this *Command) SaveSn(id int64, sn string) (affect int64) {
	stmt, _ := this.db.Prepare("update " + this.table + " set sn=? where id=?")
	res, _ := stmt.Exec(sn, id)
	affect, _ = res.RowsAffected()
	return
}
