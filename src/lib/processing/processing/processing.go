package processing

import (
	handleinterface "core/service/handle-interface"
	"encoding/json"
	"lib/processing/analysis"
)

/**
 * @Function 截取数据包（不同协议需要更换方法）
 * @Auther Nelg
 * @Date 2019.05.30
 */
func Cutpack(readData []byte, buffer []byte) (dataArray [][]byte, unCompleted []byte) {
	dataArray, unCompleted = analysis.Protocol808_Cutpack(readData, buffer)
	return
}

/**
 * @Function 解析数据包（不同协议需要更换方法），返回客户端编号（用于跟连接套接字建立联系）、包内容（只能为字符串，用于存入redis）
 * @Auther Nelg
 * @Date 2019.05.30
 */
func Resolvepack(cmd []byte) (content string, clientNumber string, err error) {
	data, err := analysis.Protocol808_Resolvepack(cmd)
	if err != nil {
		return "", "", err
	}
	clientNumber = data.ClientNumber
	jsonData, err := json.MarshalIndent(data, "", "")
	if err != nil {
		return "", "", err
	}
	content = string(jsonData)
	return
}

/**
 * @Function 获取业务对象（不同协议需要更换方法）
 * @Auther Nelg
 * @Date 2019.05.30
 */
func GetBusinessObject(content string) (unit handleinterface.HandleInterface, err error) {
	unit, err = analysis.Protocol808_GetBusinessObject(content)
	return
}

/**
 * @Function 重组发送数据包（不同协议需要更换方法）
 * @Auther Nelg
 * @Date 2019.05.31
 */
func RebuildSendPackage(content string) (pkg []byte, err error) {
	pkg, err = analysis.Protocol808_RebuildSendPackage(content)
	return
}
