package app

import (
	"context"

	"libs/common/grace"
	"libs/common/logger"
	"libs/common/metrics"
	"libs/common/service_middleware"
	"libs/common/validator"

	"template/internal/app/config"
	"template/internal/app/connections"
	"template/internal/app/start"
	"template/internal/app/store"
	"template/internal/services/book"
)

func Run(filenames ...string) {
	cfg, err := config.New(filenames...)
	if err != nil {
		panic(err)
	}

	// Configure logger
	err = logger.Configure(cfg.Logger.Level, cfg.Logger.Format)
	if err != nil {
		panic(err)
	}

	// Buffer size: number of writers that start
	// Setting more than needed just in case
	errs := make(chan error, 50)

	conns, err := connections.New(cfg)
	if err != nil {
		logger.Fatalf(context.Background(), "can't create connections: %s", err)
	}

	// Инициализация middleware
	mw := setupMiddleware()

	st := store.NewRepositoryStore(conns)

	clients := store.NewClientStore(conns)

	bookService := book.New(st, clients, mw)

	listeners := make([]grace.Service, 0)
	listeners = append(listeners, start.HTTP(bookService, cfg.HTTP.RequestTimeoutSeconds, cfg.HTTP.Addr, errs))
	listeners = append(listeners,
		start.KafkaConsumer(bookService, cfg.Kafka.TimeOutSeconds, conns.Consumer, conns.Producer, errs))

	graceful := grace.KillThemSoftly(listeners...)
	graceful.Shutdown(errs, logger.L(), conns)
}

func setupMiddleware() *service_middleware.ServiceMiddleware {
	// Инициализация метрик
	requestCounter := metrics.NewRequestCounter("api", "service")
	errorCounter := metrics.NewErrorCounter("api", "service")
	latencyHistogram := metrics.NewRequestLatency("api", "service")

	// Инициализация валидатора
	validator, err := validator.New()
	if err != nil {
		panic(err)
	}

	// Создание цепочки service_middleware
	return service_middleware.NewServiceMiddleware(
		service_middleware.MetricsMiddleware(requestCounter, errorCounter, latencyHistogram),
		service_middleware.LoggingMiddleware(),
		service_middleware.ValidationMiddleware(validator),
		service_middleware.WaitContextCancel(),
	)
}
