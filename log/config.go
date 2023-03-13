package log

// Configuration for logging
type Config struct {
	FilePrefix            string // 日志前缀
	ConsoleLoggingEnabled bool   // 开启控制台输出
	FileLoggingEnabled    bool   // 开启文件存储
	Directory             string // 日志文件保存目录
	MaxSize               int    // 最大日志大小 MB
	MaxBackups            int    // 最大保留文件数量
	MaxAge                int    // 文件最大保留天数
	ShowFileLine          bool   // 显示代码文件行数
	LogLevel              int8   // 日志等级
}

var DefaultConfig = &Config{
	FilePrefix:            "debug",
	ConsoleLoggingEnabled: true,
	FileLoggingEnabled:    false,
	Directory:             "./",
	MaxSize:               0,
	MaxBackups:            0,
	MaxAge:                0,
	ShowFileLine:          true,
	LogLevel:              -1,
}
