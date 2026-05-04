package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 封装的 zap 日志实例
type Logger struct {
	*zap.Logger
}

// logEncoder 自定义编码器，支持时间分片和前缀
type logEncoder struct {
	zapcore.Encoder
	logDir      string
	outLogger   *lumberjack.Logger
	errLogger   *lumberjack.Logger
	currentDate string
	prefix      string
}

const (
	BlueColor   = "\033[34m"
	YellowColor = "\033[33m"
	RedColor    = "\033[31m"
	ResetColor  = "\033[0m"
)

// myEncodeLevel 自定义日志级别颜色（开发环境）
func myEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.InfoLevel:
		enc.AppendString(BlueColor + "INFO" + ResetColor)
	case zapcore.WarnLevel:
		enc.AppendString(YellowColor + "WARN" + ResetColor)
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		enc.AppendString(RedColor + "ERROR" + ResetColor)
	default:
		enc.AppendString(level.String())
	}
}

// openLogFiles 初始化日志文件
func (e *logEncoder) openLogFiles(date string) error {
	logDir := filepath.Join(e.logDir, date)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// 普通日志文件
	if e.outLogger != nil {
		e.outLogger.Close()
	}
	e.outLogger = &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "out.log"),
		MaxSize:    100, // MB
		MaxBackups: 30,  // 保留30个备份
		MaxAge:     90,  // 保留90天
		Compress:   true,
	}

	// 错误日志文件
	if e.errLogger != nil {
		e.errLogger.Close()
	}
	e.errLogger = &lumberjack.Logger{
		Filename:   filepath.Join(logDir, "err.log"),
		MaxSize:    100,
		MaxBackups: 30,
		MaxAge:     90,
		Compress:   true,
	}

	e.currentDate = date
	return nil
}

func (e *logEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 检查日期并更新文件句柄
	now := time.Now().Format("2006-01-02")
	if e.currentDate != now || e.outLogger == nil {
		if err := e.openLogFiles(now); err != nil {
			fmt.Fprintf(os.Stderr, "failed to open log files: %v\n", err)
		}
	}

	// 调用原始 EncodeEntry
	buff, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}
	data := buff.String()
	buff.Reset()
	buff.AppendString(e.prefix + data)

	// 写入普通日志文件
	if e.outLogger != nil {
		if _, err := e.outLogger.Write([]byte(buff.String() + "\n")); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to out.log: %v\n", err)
		}
	}

	// 写入错误日志文件（仅 ERROR 级别及以上）
	if entry.Level >= zapcore.ErrorLevel && e.errLogger != nil {
		if _, err := e.errLogger.Write([]byte(buff.String() + "\n")); err != nil {
			fmt.Fprintf(os.Stderr, "failed to write to err.log: %v\n", err)
		}
	}

	return buff, nil
}

// NewLogger 初始化日志库
// env: "production" 或其他（如 "development"）
// logDir: 日志目录（如 "logs"）
// prefix: 日志前缀（如 "[myApp] "）
func NewLogger(env, logDir, prefix string) (*Logger, error) {
	var cfg zap.Config
	if env == "production" {
		// 生产环境：JSON 格式
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.TimeKey = "timestamp"
	} else {
		// 开发环境：彩色控制台
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = myEncodeLevel
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// 创建自定义编码器
	encoder := &logEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		logDir:  logDir,
		prefix:  prefix,
	}

	// 初始化文件句柄
	if err := encoder.openLogFiles(time.Now().Format("2006-01-02")); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize log files: %v\n", err)
	}

	// 创建 Core
	var core zapcore.Core
	if env == "production" {
		// 生产环境：仅文件输出
		outCore := zapcore.NewCore(encoder, zapcore.AddSync(encoder.outLogger), zapcore.InfoLevel)
		errCore := zapcore.NewCore(encoder, zapcore.AddSync(encoder.errLogger), zapcore.ErrorLevel)
		core = zapcore.NewTee(outCore, errCore)
	} else {
		// 开发环境：文件+控制台
		outCore := zapcore.NewCore(encoder, zapcore.AddSync(encoder.outLogger), zapcore.InfoLevel)
		errCore := zapcore.NewCore(encoder, zapcore.AddSync(encoder.errLogger), zapcore.ErrorLevel)
		consoleCore := zapcore.NewCore(
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		)
		core = zapcore.NewTee(outCore, errCore, consoleCore)
	}

	// 创建 Logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)

	logger.Info("配置与日志初始化完成")
	return &Logger{logger}, nil
}

// Close 关闭日志，同步缓冲区
func (l *Logger) Close() error {
	if err := l.Logger.Sync(); err != nil {
		return fmt.Errorf("failed to sync logger: %w", err)
	}
	return nil
}

// WithFields 添加公共字段
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return &Logger{l.Logger.With(zapFields...)}
}

func convertToFields(args []interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case zap.Field:
			fields = append(fields, v)
		case error:
			fields = append(fields, zap.Error(v))
		case string:
			// 默认当作 message 字段处理，也可以选择忽略或其它逻辑
			fields = append(fields, zap.String("msg", v))
		case int:
			fields = append(fields, zap.Int("value", v))
		// 你也可以继续支持更多类型
		default:
			fields = append(fields, zap.Any("value", v))
		}
	}
	return fields
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	fields := convertToFields(args)
	l.Logger.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	fields := convertToFields(args)
	l.Logger.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
	fields := convertToFields(args)
	l.Logger.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	fields := convertToFields(args)
	l.Logger.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	fields := convertToFields(args)
	l.Logger.Fatal(msg, fields...)
}
