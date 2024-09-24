package logger

import (
	"log"
	"os"
)

var Logger *log.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)

// func InitLogger() {

// 	Logger =
// }

// Funciones para distintos niveles de log
func Info(message string) {
	Logger.SetPrefix("INFO: ")
	Logger.Println(message)
}

func Warn(message string) {
	Logger.SetPrefix("WARN: ")
	Logger.Println(message)
}

func Error(message string) {
	Logger.SetPrefix("ERROR: ")
	Logger.Println(message)
}
