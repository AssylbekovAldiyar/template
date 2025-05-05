package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"libs/common/ctxconst"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const JSONFormat = "json"

var _global = &Logger{
	zapLogger: zap.NewNop(),
}

type Logger struct {
	zapLogger *zap.Logger
}

type Field struct {
	Key   string
	Value interface{}
}

func Configure(level int, format string) error {
	zapLevel := zapcore.Level(level)
	if !isValidLevel(zapLevel) {
		return errors.New("unknown logger level")
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(time.RFC3339Nano))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var encoder zapcore.Encoder
	switch format {
	case JSONFormat:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)

	logger := zap.New(core, zap.AddCaller())

	_global.zapLogger = logger

	return nil
}

func L() *Logger {
	return _global
}

func isValidLevel(level zapcore.Level) bool {
	return level >= zapcore.DebugLevel && level <= zapcore.FatalLevel
}

func fieldsToZapFields(fields []Field) []zap.Field {
	var zapFields []zap.Field
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}
	return zapFields
}

func (l *Logger) readCtx(ctx context.Context) []zap.Field {
	var fields []zap.Field
	if ctx != nil {
		if ctxconst.GetRequestID(ctx) != nil {
			fields = append(fields, zap.Any("request_id", ctxconst.GetRequestID(ctx)))
		}
		if ctxconst.GetUserID(ctx) != nil {
			fields = append(fields, zap.Any("user_id", ctxconst.GetUserID(ctx)))
		}
		if ctxconst.GetUserPhoneNumber(ctx) != nil {
			fields = append(fields, zap.Any("user_phone_number", ctxconst.GetUserPhoneNumber(ctx)))
		}
	}
	return fields
}

func (l *Logger) Infof(ctx context.Context, msg string, args ...interface{}) {
	l.zapLogger.Info(fmt.Sprintf(msg, args...), l.readCtx(ctx)...)
}

func (l *Logger) Errorf(ctx context.Context, msg string, args ...interface{}) {
	l.zapLogger.Error(fmt.Sprintf(msg, args...), l.readCtx(ctx)...)
}

func (l *Logger) Debugf(ctx context.Context, msg string, args ...interface{}) {
	l.zapLogger.Debug(fmt.Sprintf(msg, args...), l.readCtx(ctx)...)
}

func (l *Logger) Warnf(ctx context.Context, msg string, args ...interface{}) {
	l.zapLogger.Warn(fmt.Sprintf(msg, args...), l.readCtx(ctx)...)
}

func (l *Logger) Fatalf(ctx context.Context, msg string, args ...interface{}) {
	l.zapLogger.Fatal(fmt.Sprintf(msg, args...), l.readCtx(ctx)...)
}

func (l *Logger) Panicf(ctx context.Context, msg string, args ...interface{}) {
	l.zapLogger.Panic(fmt.Sprintf(msg, args...), l.readCtx(ctx)...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.zapLogger.Sugar().Infof(format, args...)
}

func Infof(ctx context.Context, msg string, args ...interface{}) {
	_global.Infof(ctx, msg, args...)
}

func Errorf(ctx context.Context, msg string, args ...interface{}) {
	_global.Errorf(ctx, msg, args...)
}

func Debugf(ctx context.Context, msg string, args ...interface{}) {
	_global.Debugf(ctx, msg, args...)
}

func Warnf(ctx context.Context, msg string, args ...interface{}) {
	_global.Warnf(ctx, msg, args...)
}

func Fatalf(ctx context.Context, msg string, args ...interface{}) {
	_global.Fatalf(ctx, msg, args...)
}

func Panicf(ctx context.Context, msg string, args ...interface{}) {
	_global.Panicf(ctx, msg, args...)
}

func Printf(format string, args ...interface{}) {
	_global.Printf(format, args...)
}

func WithFields(fields ...Field) *Logger {
	newZapLogger := _global.zapLogger.With(fieldsToZapFields(fields)...)
	return &Logger{zapLogger: newZapLogger}
}
