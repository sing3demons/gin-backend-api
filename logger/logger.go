package logger

import (
	"fmt"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerFactory interface {
	NewLogger() (*zap.Logger, func(), error)
}

type loggerFactory struct {
	logPath string
}

func NewLoggerFactory(logPath string) LoggerFactory {
	os.MkdirAll(logPath, 0755)
	return &loggerFactory{
		logPath: logPath,
	}
}

func (lf *loggerFactory) NewLogger() (*zap.Logger, func(), error) {
	now := time.Now()
	logfile := path.Join(lf.logPath, fmt.Sprintf("%s.log", now.Format("2006-01-02")))

	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)
	pe.EncodeTime = zapcore.RFC3339TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	core := zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(file), highPriority),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), highPriority),
		)
	})

	log, _ := zap.NewProduction(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel))

	close := func() {
		log.Sync()
		file.Close()
	}

	return log, close, nil
}
