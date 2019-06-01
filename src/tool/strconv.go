package tool

func StrPad(str string, character string, length int, direction string) (newstr string) {
	l := length - len(str)
	if l > 0 {
		var tmp string
		for i := 0; i < l; i++ {
			tmp += character
		}
		if direction == "RIGHT" {
			newstr = str + tmp
		} else {
			newstr = tmp + str
		}
	}
	return
}
