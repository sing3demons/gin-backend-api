package logger

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
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

	log, _ := zap.NewProduction(core, zap.AddStacktrace(zap.WarnLevel))

	close := func() {
		log.Sync()
		file.Close()
	}

	return log, close, nil
}

type Config struct {
	TimeFormat string
	UTC        bool
	SkipPaths  []string
	TraceID    bool // optionally log Open Telemetry TraceID
	Context    func(c *gin.Context) []zapcore.Field
}

func ZapLogger(logger *zap.Logger) gin.HandlerFunc {
	conf := &Config{}
	conf.UTC = true
	conf.TimeFormat = time.RFC3339
	skipPaths := make(map[string]bool, len(conf.SkipPaths))
	for _, path := range conf.SkipPaths {
		skipPaths[path] = true
	}
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		if _, ok := skipPaths[path]; !ok {
			end := time.Now()
			latency := end.Sub(start)
			if conf.UTC {
				end = end.UTC()
			}

			fields := []zapcore.Field{
				zap.Int("status", c.Writer.Status()),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("ip", c.ClientIP()),
				zap.String("user-agent", c.Request.UserAgent()),
				zap.Duration("latency", latency),
				zap.String("X-Request-Id", c.Writer.Header().Get("X-Request-Id")),
				zap.String("Response-Json", c.Writer.Header().Get("Response-Json")),
			}
			if conf.TimeFormat != "" {
				fields = append(fields, zap.String("time", end.Format(conf.TimeFormat)))
			}

			if conf.Context != nil {
				fields = append(fields, conf.Context(c)...)
			}

			if len(c.Errors) > 0 {
				// Append error field if this is an erroneous request.
				for _, e := range c.Errors.Errors() {
					logger.Error(e, fields...)
				}
			} else {
				logger.Info(path, fields...)
			}
		}
	}
}
