package service_middleware

import (
	"context"

	"libs/common/logger"
	"libs/common/validator"
)

// ValidationMiddleware возвращает service_middleware для валидации
func ValidationMiddleware(validator *validator.Validator) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if _, err := validator.Validate(req); err != nil {
				logger.WithFields(logger.Field{Key: "error", Value: err}).Errorf(ctx, "validation failed")

				return nil, err
			}

			return next(ctx, req)
		}
	}
}
