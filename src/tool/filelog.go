package tool

import (
	"log"
	"os"
	"time"
)

/**
 * @Function 写入文件日志
 * @Auther Nelg
 * @Date 2019.06.01
 */
func WriteLog(fileName string, content string, prefix string)  {
	filePath := "../log/" + time.Now().Format("20060102") + fileName;

	logFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		log.Fatalln(err.Error())
	}

	loger := log.New(logFile, prefix, log.Ldate|log.Ltime|log.Lshortfile)
	loger.Println(content);

	err = logFile.Close()
}