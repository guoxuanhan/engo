package logger

import "fmt"

// 日志等级
var logLevel int32 = LOG_LEVEL_DEBUG

// 日志驱动
var logDriver int32 = LOG_DRIVER_DEFAULT

// 定义日志驱动方式
const (
	LOG_DRIVER_DEFAULT = iota // 默认日志驱动
	LOG_DRIVER_ZAP            // zap日志驱动
)

// 定义日志等级
const (
	LOG_LEVEL_DEBUG = iota // 调试等级
	LOG_LEVEL_INFO         // 信息等级
	LOG_LEVEL_WARN         // 警告等级
	LOG_LEVEL_ERROR        // 错误等级
)

// 设置日志驱动方式
func SetLogDriver(driver int32) {
	switch driver {
	case LOG_DRIVER_DEFAULT:
		logDriver = LOG_DRIVER_DEFAULT
	case LOG_DRIVER_ZAP:
		logDriver = LOG_DRIVER_ZAP
	default:
		println("invalid log driver: ", driver)
	}
}

// 设置日志等级
func SetLogLevel(level string) {
	switch level {
	case "debug":
		logLevel = LOG_LEVEL_DEBUG
	case "info":
		logLevel = LOG_LEVEL_INFO
	case "warn":
		logLevel = LOG_LEVEL_WARN
	case "error":
		logLevel = LOG_LEVEL_ERROR
	default:
		println("invalid log level: ", level)
	}
}

// 错误日志
func Error(v ...any) {
	switch logDriver {
	case LOG_DRIVER_ZAP:
		sugarLogger.Error(v...)
	default:
		add(formatErr(v...))
	}
}

// 警告日志
func Warn(v ...any) {
	if logLevel > LOG_LEVEL_WARN {
		return
	}

	switch logDriver {
	case LOG_DRIVER_ZAP:
		sugarLogger.Warn(v...)
	default:
		add(formatWarn(v...))
	}
}

// 信息日志
func Info(v ...any) {
	if logLevel > LOG_LEVEL_INFO {
		return
	}

	switch logDriver {
	case LOG_DRIVER_ZAP:
		sugarLogger.Info(v...)
	default:
		add(formatInfo(v...))
	}
}

// 调试日志
func Debug(v ...any) {
	if logLevel > LOG_LEVEL_DEBUG {
		return
	}

	switch logDriver {
	case LOG_DRIVER_ZAP:
		sugarLogger.Debug(v...)
	default:
		add(formatDebug(v...))
	}
}

// 格式化的错误日志
func Errorf(format string, v ...any) {
	Error(fmt.Sprintf(format, v...))
}

// 格式化的警告日志
func Warnf(format string, v ...any) {
	Warn(fmt.Sprintf(format, v...))
}

// 格式化的信息日志
func Infof(format string, v ...any) {
	Info(fmt.Sprintf(format, v...))
}

// 格式化的调试日志
func Debugf(format string, v ...any) {
	Debug(fmt.Sprintf(format, v...))
}
