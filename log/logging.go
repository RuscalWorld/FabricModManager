package log

import (
	"github.com/gookit/color"
	"log"
	"os"
)

var (
	InfoLogger    *log.Logger
	FineLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	FatalLogger   *log.Logger
)

func SetupLogging() {
	InfoLogger = log.New(os.Stdout, color.LightBlue.Sprintf("[INFO] "), 0)
	FineLogger = log.New(os.Stdout, color.Green.Sprintf("[FINE] "), 0)
	WarningLogger = log.New(os.Stdout, color.LightYellow.Sprintf("[WARN] "), 0)
	ErrorLogger = log.New(os.Stdout, color.LightRed.Sprintf("[ERROR] "), 0)
	FatalLogger = log.New(os.Stdout, color.Red.Sprintf("[FATAL] "), 0)
}

func Info(a ...interface{}) {
	InfoLogger.Println(a...)
}

func Fine(a ...interface{}) {
	FineLogger.Println(a...)
}

func Warn(a ...interface{}) {
	WarningLogger.Println(a...)
}

func Error(a ...interface{}) {
	ErrorLogger.Println(a...)
}

func Fatal(a ...interface{}) {
	FatalLogger.Fatalln(a...)
}
