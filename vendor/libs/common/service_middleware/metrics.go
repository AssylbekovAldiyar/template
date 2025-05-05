package service_middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/metrics"

	"libs/common/ctxconst"
)

// MetricsMiddleware возвращает service_middleware для метрик
func MetricsMiddleware(requestCount, requestError metrics.Counter, requestLatency metrics.Histogram) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			methodName := ctxconst.GetMethodName(ctx)

			defer func(begin time.Time) {
				requestCount.With("method", methodName).Add(1)
				requestLatency.With("method", methodName).Observe(time.Since(begin).Seconds())
			}(time.Now())

			resp, err := next(ctx, req)
			if err != nil {
				requestError.With("method", methodName).Add(1)
			}

			return resp, err
		}
	}
}
