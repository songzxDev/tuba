package log

import (
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

// Config is log config .
type Config struct {
	LogPath string
	AppName string
	Debug   bool
}

// Init initialize a log config .
func Init(c *Config) {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	})
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	fatalLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.FatalLevel
	})
	var cores []zapcore.Core
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	if c.Debug {
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleDebugging, highPriority))
	}
	cores = append(cores, zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority))
	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(getWriter(c.LogPath+c.AppName+"_debug.log")), debugLevel))
	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(getWriter(c.LogPath+c.AppName+"_info.log")), infoLevel))
	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(getWriter(c.LogPath+c.AppName+"_warn.log")), warnLevel))
	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(getWriter(c.LogPath+c.AppName+"_error.log")), errorLevel))
	cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(getWriter(c.LogPath+c.AppName+"_fatal.log")), fatalLevel))
	logger = zap.New(zapcore.NewTee(cores...), zap.AddCaller()).Sugar()
}

// TimeEncoder time encoder .
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(time.RFC3339Nano))
}

// Debugf log
func Debugf(msg string, args ...interface{}) {
	logger.Debugf(msg, args...)
}

// Debug log
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Infof log
func Infof(msg string, args ...interface{}) {
	logger.Infof(msg, args...)
}

// Info log
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Errorf log
func Errorf(msg string, args ...interface{}) {
	logger.Errorf(msg, args...)
}

// Error log
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Warnf log
func Warnf(msg string, args ...interface{}) {
	logger.Warnf(msg, args...)
}

// Warn log
func Warn(msg string, args ...interface{}) {
	logger.Warn(args...)
}

// Fatalf send log fatalf
func Fatalf(msg string, args ...interface{}) {
	logger.Fatalf(msg, args...)
}

// Fatal send log fatal
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"_%Y%m%d%H.log",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return hook
}
