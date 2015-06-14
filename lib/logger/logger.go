package logger

import (
	"log"
	"os"
	"time"
)

func PutLog(msg, dir, logname string) {
	filename := logname + ".log"
	path := dir + "/" + time.Now().Format("20060102")
	os.MkdirAll(path, 0777)
	flagMode := log.Ldate | log.Ltime
	if logname == "error" {
		flagMode = log.Ldate | log.Ltime | log.Llongfile
	}

	logfile, err := os.OpenFile(path+"/"+filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Printf("%s\r\n", msg)
		os.Exit(-1)
	}
	defer logfile.Close()
	logger := log.New(logfile, "", flagMode)
	logger.Println(msg)
}
