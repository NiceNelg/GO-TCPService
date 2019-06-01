package tool_protocol808

/**
 * @Function 数据转义
 * @Auther Nelg
 * @Date 2019.05.30
 */
func Escape(cmd []byte) (data []byte) {
	temp := make([]byte, 0)
	/*转义*/
	cmdMaxIndex := len(cmd) - 1
	for index, value := range cmd {
		if index <= 0 || index >= cmdMaxIndex {
			temp = append(temp, value)
			continue
		}
		if value == 0x7e {
			temp = append(temp, 0x7d, 0x02)
		} else if value == 0x7d {
			temp = append(temp, 0x7d, 0x01)
		} else {
			temp = append(temp, value)
		}
	}
	data = temp
	return

}

/**
 * @Function 数据反转义
 * @Auther Nelg
 * @Date 2019.05.30
 */
func ReverseEscape(cmd []byte) (data []byte) {
	if len(cmd) <= 0 {
		return
	}
	temp := make([]byte, 0)
	/*反转义*/
	cmdMaxIndex := len(cmd) - 1
	for index := 0; index <= cmdMaxIndex; {
		if index <= 0 || index >= cmdMaxIndex {
			temp = append(temp, cmd[index])
			index++
			continue
		}
		if cmd[index] == 0x7d && index+1 < cmdMaxIndex && cmd[index+1] == 0x02 {
			temp = append(temp, 0x7e)
			index += 2
		} else if cmd[index] == 0x7d && index+1 < cmdMaxIndex && cmd[index+1] == 0x01 {
			temp = append(temp, 0x7d)
			index += 2
		} else {
			temp = append(temp, cmd[index])
			index++
		}
	}
	data = temp
	return
}
