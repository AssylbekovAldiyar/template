package service_middleware

import (
	"context"
	"errors"

	"libs/common/logger"
)

// WaitContextCancel возвращает middleware для обработки таймаутов и паник
func WaitContextCancel() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			respCh := make(chan struct {
				resp interface{}
				err  error
			}, 1)

			go func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Errorf(ctx, "panic recovered: %v", r)
						respCh <- struct {
							resp interface{}
							err  error
						}{
							resp: nil,
							err:  errors.New("panic recovered"),
						}
					}
				}()

				resp, err := next(ctx, req)
				respCh <- struct {
					resp interface{}
					err  error
				}{
					resp: resp,
					err:  err,
				}
			}()

			// Ожидание результата или отмены контекста
			select {
			case <-ctx.Done():
				return nil, errors.New("reached timeout")
			case result := <-respCh:
				return result.resp, result.err
			}
		}
	}
}
