package service_middleware

import (
	"context"
	"time"

	"libs/common/logger"
)

// LoggingMiddleware возвращает service_middleware для логирования
func LoggingMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			logger.WithFields(logger.Field{Key: "request", Value: req}).Infof(ctx, "handling request")

			startTime := time.Now()

			resp, err := next(ctx, req)

			fields := []logger.Field{
				{Key: "elapsed", Value: time.Since(startTime)},
				{Key: "response", Value: resp},
			}

			if err != nil {
				fields = append(fields, logger.Field{Key: "error", Value: err})
				logger.WithFields(fields...).Errorf(ctx, "request failed")
			} else {
				logger.WithFields(fields...).Infof(ctx, "request completed successfully")
			}

			return resp, err
		}
	}
}
