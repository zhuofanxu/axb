package logger

import (
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Manager struct {
	mu      sync.RWMutex
	loggers map[string]*zap.Logger
	config  *Config
}

func NewManager(config *Config) *Manager {
	return &Manager{
		loggers: make(map[string]*zap.Logger),
		config:  config,
	}
}

// GetLogger 获取指定模块的logger
func (m *Manager) GetLogger(module string) *zap.Logger {
	m.mu.RLock()
	if logger, exists := m.loggers[module]; exists {
		m.mu.RUnlock()
		return logger
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	// 双重检查, 避免多个线程同时创建同一个模块的logger
	if logger, exists := m.loggers[module]; exists {
		return logger
	}

	logger := m.createModuleLogger(module)
	m.loggers[module] = logger

	return logger
}

func (m *Manager) createModuleLogger(module string) *zap.Logger {
	moduleConfig, exists := m.config.Modules[module]
	if !exists {
		moduleConfig = ModuleLogConfig{
			Level:   m.config.Level,
			Enabled: true,
			SubDir:  "web",
		}
	}

	if !moduleConfig.Enabled {
		return zap.NewNop()
	}

	logDir := fmt.Sprintf("%s/%s", m.config.BaseDir, moduleConfig.SubDir)

	cores := []zapcore.Core{
		getEncoderCore(fmt.Sprintf("%s/info.log", logDir),
			zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
				return lev == zap.InfoLevel
			}), m.config),
		getEncoderCore(fmt.Sprintf("%s/warn.log", logDir),
			zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
				return lev == zap.WarnLevel
			}), m.config),
		getEncoderCore(fmt.Sprintf("%s/error.log", logDir),
			zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
				return lev >= zap.ErrorLevel
			}), m.config),
	}

	if moduleConfig.Level == "debug" {
		cores = append(cores, getEncoderCore(
			fmt.Sprintf("%s/debug.log", logDir),
			zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
				return lev == zap.DebugLevel
			}), m.config))
	}

	if m.config.LogInConsole {
		cores = append(cores, getConsoleCore(
			zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
				return lev >= zap.DebugLevel
			}), m.config))
	}

	logger := zap.New(zapcore.NewTee(cores...))
	if m.config.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	if m.config.LogInConsole {
		logger = logger.Named(fmt.Sprintf("[%s]", strings.ToUpper(module)))
	}

	return logger
	//return logger.With(zap.String("module", module))
}
