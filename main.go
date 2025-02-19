package main

import (
	"engo/log"
	"errors"
	"time"
)

func init() {
	log.StartZapLogger()
	log.SetLogDriver(log.LOG_DRIVER_ZAP)
	log.SetLogLevel("debug")
}

func main() {
	log.Debug("服务器启动， 开始初始化。。。")
	log.Info("当前时间： ", time.Now().Format("2006-01-02 15:04:05"))

	// 格式化日志示例
	serverPort := 8080
	log.Infof("HTTP服务器启动于端口: %d", serverPort)

	// 模拟一些警告场景
	memoryUsage := 85.5
	if memoryUsage > 80 {
		log.Warnf("内存使用率较高: %.2f%%", memoryUsage)
	}

	// 模拟错误处理
	if err := simulateError(); err != nil {
		log.Error("操作失败:", err)
	}

	// 模拟业务流程
	processOrder("ORDER123")

	// 性能日志示例
	startTime := time.Now()
	simulateWork()
	log.Debugf("操作耗时: %v", time.Since(startTime))
}

// 模拟订单处理
func processOrder(orderID string) {
	log.Info("开始处理订单:", orderID)

	// 模拟订单处理步骤
	steps := []string{"验证订单", "检查库存", "处理支付", "更新订单状态"}
	for i, step := range steps {
		log.Debugf("步骤 %d: %s", i+1, step)
		simulateWork()

		// 模拟随机警告
		if i == 2 {
			log.Warn("支付处理较慢，请检查支付网关状态")
		}
	}

	log.Info("订单处理完成:", orderID)
}

// 模拟错误
func simulateError() error {
	return errors.New("模拟错误")
}

// 模拟工作负载
func simulateWork() {
	time.Sleep(time.Millisecond * 100)
}
