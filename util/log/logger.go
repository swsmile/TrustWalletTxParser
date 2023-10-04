package log

import (
	"fmt"
	"log"
	"os"
)

const (
	logFile = "log/data.log"
)

var (
	warningLogger *log.Logger
	infoLogger    *log.Logger
	errorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// TODO to fix the wrong file name and line number
	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Warning(format string, a ...any) {
	warningLogger.Println(fmt.Sprintf(format, a...))
}

func Info(format string, a ...any) {
	infoLogger.Println(fmt.Sprintf(format, a...))
}

func Error(format string, a ...any) {
	errorLogger.Println(fmt.Sprintf(format, a...))
}
