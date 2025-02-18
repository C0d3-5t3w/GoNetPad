package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func Init(logFile string) error {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	InfoLogger = log.New(file, "INFO: ", log.LstdFlags|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.LstdFlags|log.Lshortfile)

	return nil
}

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lshortfile)
}
