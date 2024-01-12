package init

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	DebugLog *log.Logger
	InfoLog  *log.Logger
	logf     *os.File
)

func Log_init() {
	T := time.Now()
	_, err := os.Stat("src/errorlog")
	if err != nil {
		os.Mkdir("src/errorlog", 0777)
		os.Chmod("src/errorlog", 0777)
	}
	logfile := fmt.Sprintf("src/errorlog/%04d-%02d-%02d.log", T.Year(), T.Month(), T.Day())
	log.SetFlags(log.Ltime)
	logf, _ = os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	log.SetOutput(logf)
	InfoLog = log.New(logf, "[Info] - ", log.Ltime)
	DebugLog = log.New(logf, "[Debug] - ", log.Ltime)

}
