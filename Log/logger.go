package logger

/////////////////////////
// 日志类 便于日志输出 //
/////////////////////////

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	callDepth = 1
)

type LogLevel uint8

const (
	LogLevelDEBUG LogLevel = 1 << iota
	LogLevelINFO
	LogLevelWARN
	LogLevelERROR
)

type tlog struct {
	logLevel LogLevel
	logFile  *os.File
}

func (l *tlog) Write(p []byte) (int, error) {
	return l.logFile.Write(p)
}

var myLogger *tlog

func Init(model string, logLevel LogLevel) {
	out, err := os.OpenFile(fmt.Sprintf("%v.log", model), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	myLogger = &tlog{
		logFile:  out,
		logLevel: logLevel,
	}
	output := io.MultiWriter(myLogger, os.Stdout)
	log.SetOutput(output)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
}

func Debugf(format string, params ...interface{}) {
	if myLogger.logLevel&LogLevelDEBUG == 0 {
		return
	}
	log.SetPrefix("[DEBUG] ")
	log.Output(callDepth, fmt.Sprintf(format, params...))
}
func Debug(msg string) {
	if myLogger.logLevel&LogLevelDEBUG == 0 {
		return
	}
	log.SetPrefix("[DEBUG] ")
	log.Output(callDepth, msg)
}
func Infof(format string, params ...interface{}) {
	if myLogger.logLevel&LogLevelINFO == 0 {
		return
	}
	log.SetPrefix("[INFO] ")
	log.Output(callDepth, fmt.Sprintf(format, params...))
}
func Info(msg string) {
	if myLogger.logLevel&LogLevelINFO == 0 {
		return
	}
	log.SetPrefix("[INFO] ")
	log.Output(callDepth, msg)
}
func Warnf(format string, params ...interface{}) {
	if myLogger.logLevel&LogLevelWARN == 0 {
		return
	}
	log.SetPrefix("[WARN] ")
	log.Output(callDepth, fmt.Sprintf(format, params...))
}
func Warn(msg string) {
	if myLogger.logLevel&LogLevelWARN == 0 {
		return
	}
	log.SetPrefix("[WARN] ")
	log.Output(callDepth, msg)
}
func Errorf(format string, params ...interface{}) {
	if myLogger.logLevel&LogLevelERROR == 0 {
		return
	}
	log.SetPrefix("[ERROR] ")
	log.Output(callDepth, fmt.Sprintf(format, params...))
}
func Error(msg string) {
	if myLogger.logLevel&LogLevelERROR == 0 {
		return
	}
	log.SetPrefix("[ERROR] ")
	log.Output(callDepth, msg)
}
