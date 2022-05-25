package db

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"

	"github.com/redrru/fantasy-dota/pkg/log"
)

const (
	traceErrStr = "[DB] %s %s [%.3fms] [rows:%v] %s"
	traceStr    = "[DB] %s [%.3fms] [rows:%v] %s"
	logStr      = "[DB] %s"
)

type gormLogger struct {
	logger log.Logger
}

func newGormLogger() logger.Interface {
	return &gormLogger{logger: log.GetLogger()}
}

func (l *gormLogger) LogMode(logLevel logger.LogLevel) logger.Interface {
	return newGormLogger()
}

func (l *gormLogger) Info(ctx context.Context, format string, args ...interface{}) {
	l.logger.Debug(ctx, fmt.Sprintf(logStr, fmt.Sprintf(format, args...)))
}

func (l *gormLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	l.logger.Warn(ctx, fmt.Sprintf(logStr, fmt.Sprintf(format, args...)))
}

func (l *gormLogger) Error(ctx context.Context, format string, args ...interface{}) {
	l.logger.Error(ctx, fmt.Sprintf(logStr, fmt.Sprintf(format, args...)))
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()

	duration := float64(time.Since(begin).Nanoseconds()) / 1e6

	if err != nil {
		l.logger.Debug(ctx, fmt.Sprintf(traceErrStr, utils.FileWithLineNum(), err, duration, rows, sql))
	} else {
		l.logger.Debug(ctx, fmt.Sprintf(traceStr, utils.FileWithLineNum(), duration, rows, sql))
	}
}
