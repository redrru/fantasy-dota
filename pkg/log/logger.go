package log

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	traceID = "trace_id"
	spanID  = "span_id"
)

var (
	zapLogger *logger
	once      sync.Once
)

func init() {
	_ = GetLogger()
}

type Logger interface {
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	Sync() error
}

func GetLogger() Logger {
	once.Do(func() {
		var err error
		zapLogger, err = newLogger()
		if err != nil {
			panic(err)
		}
	})

	return zapLogger
}

type logger struct {
	zapLogger *zap.Logger
}

func newLogger() (*logger, error) {
	newZapLogger, err := zap.NewDevelopmentConfig().Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &logger{zapLogger: newZapLogger}, nil
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, withTracingFields(ctx, fields...)...)
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, withTracingFields(ctx, fields...)...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, withTracingFields(ctx, fields...)...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Error(msg, withTracingFields(ctx, fields...)...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.zapLogger.Fatal(msg, withTracingFields(ctx, fields...)...)
}

func (l *logger) Sync() error {
	return l.zapLogger.Sync()
}

func (l *logger) With(fields ...zap.Field) Logger {
	clone := l.clone()
	clone.zapLogger = clone.zapLogger.With(fields...)

	return clone
}

func (l *logger) clone() *logger {
	cp := *l
	return &cp
}

func withTracingFields(ctx context.Context, fields ...zap.Field) []zap.Field {
	ctxFields := tracingFieldsFromCtx(ctx)

	result := make([]zap.Field, 0, len(fields)+len(ctxFields))
	result = append(result, ctxFields...)
	result = append(result, fields...)

	return result
}

func tracingFieldsFromCtx(ctx context.Context) []zap.Field {
	if ctx == nil {
		return nil
	}

	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return nil
	}

	return []zap.Field{
		zap.String(traceID, spanContext.TraceID().String()),
		zap.String(spanID, spanContext.SpanID().String()),
	}
}
