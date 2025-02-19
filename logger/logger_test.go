package logger

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestLoggerBasic(t *testing.T) {
	// 清理测试日志文件
	defer func() {
		os.RemoveAll("./logs")
	}()

	// 测试默认日志驱动
	t.Run("测试默认日志驱动", func(t *testing.T) {
		SetLogLevel("debug")
		Debug("这是一条调试日志")
		Info("这是一条信息日志")
		Warn("这是一条警告日志")
		Error("这是一条错误日志")

		Debugf("这是一条格式化的调试日志: %s", "test")
		Infof("这是一条格式化的信息日志: %s", "test")
		Warnf("这是一条格式化的警告日志: %s", "test")
		Errorf("这是一条格式化的错误日志: %s", "test")
	})

	// 测试日志级别
	t.Run("测试日志级别", func(t *testing.T) {
		tests := []struct {
			level string
			want  int32
		}{
			{"debug", LOG_LEVEL_DEBUG},
			{"info", LOG_LEVEL_INFO},
			{"warn", LOG_LEVEL_WARN},
			{"error", LOG_LEVEL_ERROR},
			{"invalid", LOG_LEVEL_DEBUG}, // 无效的日志级别应保持原值
		}

		for _, tt := range tests {
			SetLogLevel(tt.level)
			if logLevel != tt.want && tt.level != "invalid" {
				t.Errorf("SetLogLevel(%s) = %d; want %d", tt.level, logLevel, tt.want)
			}
		}
	})

	// 测试Zap日志驱动
	t.Run("测试Zap日志驱动", func(t *testing.T) {
		// 启动Zap日志
		StartZapLogger()
		logDriver = LOG_DRIVER_ZAP

		// 检查日志文件是否创建
		if _, err := os.Stat("./logs"); os.IsNotExist(err) {
			t.Error("日志目录未创建")
		}

		// 写入测试日志
		Debug("Zap调试日志")
		Info("Zap信息日志")
		Warn("Zap警告日志")
		Error("Zap错误日志")

		// 检查日志文件是否存在
		today := time.Now().Format("2006-01-02")
		logFile := filepath.Join("./logs", "server-"+today+".log")
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Error("日志文件未创建:", logFile)
		}
	})

	// 测试日志级别过滤
	t.Run("测试日志级别过滤", func(t *testing.T) {
		// 设置为ERROR级别
		SetLogLevel("error")

		// Debug/Info/Warn级别的日志应该被过滤
		Debug("这条调试日志不应该出现")
		Info("这条信息日志不应该出现")
		Warn("这条警告日志不应该出现")
		Error("这条错误日志应该出现")

		// 设置回DEBUG级别
		SetLogLevel("debug")
	})
}

// 性能测试
func BenchmarkLogger(b *testing.B) {
	SetLogLevel("debug")

	b.Run("默认日志驱动-Info", func(b *testing.B) {
		logDriver = LOG_DRIVER_DEFAULT
		for i := 0; i < b.N; i++ {
			Info("这是一条基准测试日志")
		}
	})

	b.Run("默认日志驱动-Infof", func(b *testing.B) {
		logDriver = LOG_DRIVER_DEFAULT
		for i := 0; i < b.N; i++ {
			Infof("这是一条基准测试日志 %d", i)
		}
	})

	b.Run("Zap日志驱动-Info", func(b *testing.B) {
		StartZapLogger()
		logDriver = LOG_DRIVER_ZAP
		for i := 0; i < b.N; i++ {
			Info("这是一条基准测试日志")
		}
	})

	b.Run("Zap日志驱动-Infof", func(b *testing.B) {
		logDriver = LOG_DRIVER_ZAP
		for i := 0; i < b.N; i++ {
			Infof("这是一条基准测试日志 %d", i)
		}
	})
}

// 并发测试
func TestConcurrentLogging(t *testing.T) {
	defer os.RemoveAll("./logs")

	SetLogLevel("debug")
	var wg sync.WaitGroup

	// 测试并发写入默认日志
	t.Run("并发测试-默认日志驱动", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					Info("并发测试日志 goroutine:", id, "次数:", j)
				}
			}(i)
		}
		wg.Wait()
	})

	// 测试并发写入Zap日志
	t.Run("并发测试-Zap日志驱动", func(t *testing.T) {
		StartZapLogger()
		logDriver = LOG_DRIVER_ZAP

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					Info("并发测试日志 goroutine:", id, "次数:", j)
				}
			}(i)
		}
		wg.Wait()
	})
}

func TestLogLevelSwitch(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error"}
	messages := map[string]string{
		"debug": "调试信息",
		"info":  "普通信息",
		"warn":  "警告信息",
		"error": "错误信息",
	}

	for _, level := range levels {
		t.Run("测试日志级别:"+level, func(t *testing.T) {
			SetLogLevel(level)

			// 记录所有级别的日志
			Debug(messages["debug"])
			Info(messages["info"])
			Warn(messages["warn"])
			Error(messages["error"])

			// 这里可以添加日志输出验证
			// 根据当前设置的级别，检查是否只输出了应该输出的日志
		})
	}
}

func TestEdgeCases(t *testing.T) {
	// 测试空日志
	t.Run("测试空日志", func(t *testing.T) {
		Info()
		Debug()
		Warn()
		Error()
	})

	// 测试大量参数
	t.Run("测试大量参数", func(t *testing.T) {
		params := make([]interface{}, 1000)
		for i := range params {
			params[i] = i
		}
		Info(params...)
	})

	// 测试特殊字符
	t.Run("测试特殊字符", func(t *testing.T) {
		specialChars := []string{
			"\n", "\t", "\r", "\000",
			"中文", "🎉", "",
		}
		for _, char := range specialChars {
			Info("特殊字符测试:", char)
		}
	})
}
