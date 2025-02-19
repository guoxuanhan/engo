package logger

import (
	"fmt"
	"log"
	"sync/atomic"
)

// 全局日志消息队列
var msgQueue chan string

// 丢弃的日志消息数量
var dropped int32 = 0

// 日志消息队列缓冲区长度
const bufferLen = 1024

// 定义日志颜色输出格式
const (
	debugFormator = "\033[1;35m[DEBUG] %v \033[0m\n"
	infoFormator  = "\033[32m[INFO] %v \033[0m\n"
	warnFormator  = "\033[1;33m[WARN] %v \033[0m\n"
	errorFormator = "\033[1;4;31m[ERROR] %v \033[0m\n"
)

// 启动日志消息队列
func Start() {
	msgQueue = make(chan string, bufferLen)
	go worker()
}

// 停止日志消息队列
func Stop() {
	if logDriver == LOG_DRIVER_ZAP {
		_ = sugarLogger.Sync()
	}
}

// 添加日志消息
func add(msg string) {
	if msgQueue == nil {
		log.Print(msg)
		return
	}

	if len(msgQueue) >= bufferLen {
		atomic.AddInt32(&dropped, 1)
		return
	}

	if dropped > 0 {
		println("dropped ", dropped)
		dropped = 0
	}

	msgQueue <- msg
}

// 日志消息队列处理
func worker() {
	for msg := range msgQueue {
		log.Print(msg)
	}
}

// 颜色格式化调试日志
func formatDebug(v ...any) string {
	return fmt.Sprintf(debugFormator, fmt.Sprint(v...))
}

// 颜色格式化信息日志
func formatInfo(v ...any) string {
	return fmt.Sprintf(infoFormator, fmt.Sprint(v...))
}

// 颜色格式化警告日志
func formatWarn(v ...any) string {
	return fmt.Sprintf(warnFormator, fmt.Sprint(v...))
}

// 格式化错误日志
func formatErr(v ...any) string {
	return fmt.Sprintf(errorFormator, fmt.Sprintln(v...))
}
