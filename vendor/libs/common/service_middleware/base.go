package service_middleware

import (
	"context"
)

// HandlerFunc представляет собой обобщенную функцию-обработчик
type HandlerFunc func(ctx context.Context, req interface{}) (interface{}, error)

// Middleware представляет собой функцию-декоратор для HandlerFunc
type Middleware func(next HandlerFunc) HandlerFunc

// MiddlewareOption представляет опцию middleware для конкретного метода
type MiddlewareOption struct {
	Middleware Middleware
	IsRequired bool
}

// ServiceMiddleware предоставляет базовую структуру для сервисов с middleware
type ServiceMiddleware struct {
	middlewares       map[string][]MiddlewareOption
	globalMiddlewares []Middleware
}

// NewServiceMiddleware создает новый ServiceMiddleware
func NewServiceMiddleware(globalMiddlewares ...Middleware) *ServiceMiddleware {
	return &ServiceMiddleware{
		middlewares:       make(map[string][]MiddlewareOption),
		globalMiddlewares: globalMiddlewares,
	}
}

// UseForMethod добавляет middleware для конкретного метода
func (m *ServiceMiddleware) UseForMethod(methodName string, middleware Middleware) {
	if m.middlewares[methodName] == nil {
		m.middlewares[methodName] = make([]MiddlewareOption, 0)
	}
	m.middlewares[methodName] = append(m.middlewares[methodName], MiddlewareOption{
		Middleware: middleware,
		IsRequired: true,
	})
}

// Wrap оборачивает метод сервиса в middleware chain
func (m *ServiceMiddleware) Wrap(methodName string, handler HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		// Создаем цепочку middleware для конкретного метода
		chain := make([]Middleware, 0)

		// Добавляем глобальные middleware
		chain = append(chain, m.globalMiddlewares...)

		// Добавляем middleware специфичные для метода
		if methodOptions, exists := m.middlewares[methodName]; exists {
			for _, option := range methodOptions {
				if option.IsRequired {
					chain = append(chain, option.Middleware)
				}
			}
		}

		// Применяем цепочку middleware
		wrapped := handler
		for i := len(chain) - 1; i >= 0; i-- {
			wrapped = chain[i](wrapped)
		}

		return wrapped(ctx, req)
	}
}

// Optional возвращает версию middleware, которая может быть включена/выключена
func Optional(middleware Middleware) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return next // По умолчанию middleware отключен
	}
}
