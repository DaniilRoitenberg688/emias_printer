package logger

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type loggerKeyType string
type requestIdType string

var (
	loggerKey    loggerKeyType = "logger"
	requestIdKey requestIdType = "requestId"
)

type Logger struct {
	l *zap.Logger
}

func NewLogger(ctx context.Context) (context.Context, error) {
	z, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	l := &Logger{
		l: z,
	}
	return context.WithValue(ctx, loggerKey, l), nil
}

func GetLoggerFromContext(ctx context.Context) *Logger {
	return ctx.Value(loggerKey).(*Logger)
}

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdKey, requestId)
}

func (l *Logger) Info(ctx context.Context, msg string, params ...zap.Field) {
	if ctx.Value(requestIdKey) != nil {
		params = append(params, zap.String(string(requestIdKey), ctx.Value(requestIdKey).(string)))
	}
	l.l.Info(msg, params...)
}

func (l *Logger) Debug(ctx context.Context, msg string, params ...zap.Field) {
	if ctx.Value(requestIdKey) != nil {
		params = append(params, zap.String(string(requestIdKey), ctx.Value(requestIdKey).(string)))
	}
	l.l.Debug(msg, params...)
}

func (l *Logger) Fatal(ctx context.Context, msg string, params ...zap.Field) {
	if ctx.Value(requestIdKey) != nil {
		params = append(params, zap.String(string(requestIdKey), ctx.Value(requestIdKey).(string)))
	}
	l.l.Fatal(msg, params...)
}


func (l *Logger) Warn(ctx context.Context, msg string, params ...zap.Field) {
	if ctx.Value(requestIdKey) != nil {
		params = append(params, zap.String(string(requestIdKey), ctx.Value(requestIdKey).(string)))
	}
	l.l.Warn(msg, params...)
}

func LoggerMidleware() func(ctx context.Context, next http.Handler) http.Handler {
	return func(ctx context.Context, next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Context().Value(requestIdKey)
			if id == nil {
				ctx = WithRequestId(ctx, uuid.NewString())
			} else {
				ctx = WithRequestId(ctx, id.(string))
			}
			GetLoggerFromContext(ctx).Info(ctx, "incoming request", zap.String("route", r.RequestURI), zap.String("method", r.Method))
			newReq := r.WithContext(ctx)
			next.ServeHTTP(w, newReq)
		})
	}
}

