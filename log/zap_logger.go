package log

import (
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 全局zap日志对象
var sugarLogger *zap.SugaredLogger

// 启动zap日志
func StartZapLogger() {
	dir := "./logs"
	_ = os.Mkdir(dir, os.ModePerm)

	writer, err := rotatelogs.New(
		dir+"/server-%Y-%m-%d.log",                 // 日志文件名
		rotatelogs.WithLinkName("logs/server.log"), // 日志文件链接名，指向的最新日志文件
		rotatelogs.WithMaxAge(time.Hour*24*30),     // 日志文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24),  // 日志文件轮转时间
	)

	if err != nil {
		fmt.Println("zap logger init error: ", err)
		panic(err)
	}

	config := zap.NewProductionEncoderConfig()       // 创建生产环境的zap日志配置
	config.EncodeTime = zapcore.ISO8601TimeEncoder   // 设置时间编码格式
	config.EncodeLevel = zapcore.CapitalLevelEncoder // 设置日志等级编码格式

	core := zapcore.NewCore(
		// zapcore.NewJSONEncoder(config),    // 创建zap日志编码器
		zapcore.NewConsoleEncoder(config), // 创建zap日志编码器
		zapcore.AddSync(writer),           // 设置日志写入器
		zap.InfoLevel,                     // 设置日志等级
	)

	// 创建zap日志对象
	logger := zap.New(
		core,                 // 设置zap日志核心
		zap.AddCaller(),      // 添加调用者信息
		zap.AddCallerSkip(1)) // 添加调用者跳过信息

	zap.ReplaceGlobals(logger) // 替换zap全局日志对象

	// 创建zap日志对象的糖化对象
	sugarLogger = logger.Sugar()
}
