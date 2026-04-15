package logger

type Config struct {
	Level         string                     `mapstructure:"level" json:"level" yaml:"level"`                           // 级别
	Format        string                     `mapstructure:"format" json:"format" yaml:"format"`                        // 输出
	Prefix        string                     `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                        // 日志前缀
	BaseDir       string                     `mapstructure:"base-dir" json:"baseDir"  yaml:"base-dir"`                  // 日志文件夹
	ShowLine      bool                       `mapstructure:"show-line" json:"showLine" yaml:"show-line"`                // 显示行
	EncodeLevel   string                     `mapstructure:"encode-level" json:"encodeLevel" yaml:"encode-level"`       // 编码级
	StacktraceKey string                     `mapstructure:"stacktrace-key" json:"stacktraceKey" yaml:"stacktrace-key"` // 栈名
	LogInConsole  bool                       `mapstructure:"log-in-console" json:"logInConsole" yaml:"log-in-console"`  // 输出控制台
	Modules       map[string]ModuleLogConfig `mapstructure:"modules" json:"modules" yaml:"modules"`                     // 模块日志配置
}

type ModuleLogConfig struct {
	Level   string `mapstructure:"level" json:"level" yaml:"level"`
	Enabled bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	SubDir  string `mapstructure:"sub-dir" json:"subDir" yaml:"subDir"` // 对应 sub-dir
}
