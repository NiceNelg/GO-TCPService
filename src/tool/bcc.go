package tool

/**
 * @Function 计算校验和
 * @Auther Nelg
 * @Date 2019.05.30
 */
func BuildBCC(content []byte) (xor byte) {
	xor = 0x00
	for _, value := range content {
		xor ^= value
	}
	return
}
