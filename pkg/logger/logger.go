// pkg/logger/logger.go
package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Fatal *log.Logger
)

func init() {
	// 创建日志文件
	logPath := "app.log"
	absPath, err := filepath.Abs(logPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path for log file: %v", err)
	}

	logFile, err := os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file at %s: %v", absPath, err)
	}

	// 打印日志文件位置
	log.Printf("Log file created at: %s", absPath)
	Info.Printf("Application logs will be saved to: %s", absPath)

	// 初始化日志器，同时输出到文件和终端
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.Lshortfile)

	// 添加文件输出
	Info.SetOutput(io.MultiWriter(os.Stdout, logFile))
	Error.SetOutput(io.MultiWriter(os.Stderr, logFile))
	Fatal.SetOutput(io.MultiWriter(os.Stderr, logFile))
}

func LogFatal(format string, v ...interface{}) {
	Fatal.Fatalf(format, v...)
}

func LogInfo(format string, v ...interface{}) {
	Info.Printf(format, v...)
}

func LogError(format string, v ...interface{}) {
	Error.Printf(format, v...)
}

// Warn
func Warn(format string, v ...interface{}) {
	Error.Printf(format, v...)
}
