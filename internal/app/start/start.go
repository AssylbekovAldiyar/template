package start

import (
	"context"
	"time"

	"template/internal/deliveries/book/kafka"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"libs/common/grace"
	kafka_lib "libs/common/kafka"
	"libs/common/logger"

	"template/internal/deliveries/book/http/v1"
	"template/internal/services/book"
)

// todo: Добавть Health :8081
// e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
// e.GET("/ready", ...)
// e.GET("/live", ...)

const (
	kafkaTestTopic = "test-topic"
)

//	@title			Template API
//	@version		1.0
//	@description	Template service API documentation

// @host		localhost
// @BasePath	/v1
// @schemes	http https
func HTTP(
	srv book.Service,
	timeoutSeconds int,
	httpAddr string,
	errs chan<- error,
) grace.Service {
	ctx := context.Background()
	lg := logger.WithFields(logger.Field{Key: "_start_type", Value: "http"})

	e := echo.New()
	e.HidePort = true
	e.HideBanner = true
	e.Use(middleware.RequestID())
	e.Logger.SetLevel(log.OFF)

	lg.Infof(ctx, "starting server on %s", httpAddr)
	go func() {
		errs <- e.Start(httpAddr)
	}()

	v1.Register(e, v1.NewV1(srv, timeoutSeconds))

	return grace.NewService("http", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			lg.Errorf(ctx, "error while shutting down http server: %v", err)
		}
	})
}

func KafkaConsumer(
	srv book.Service,
	timeoutSeconds int,
	consumer *kafka_lib.Consumer,
	producer *kafka_lib.Producer,
	errs chan<- error,
) grace.Service {
	ctx := context.Background()
	lg := logger.WithFields(logger.Field{Key: "_start_type", Value: "kafka"})

	bookHandler := kafka.NewBookHandler(srv, producer, timeoutSeconds)
	bookHandler.RegisterHandler(consumer)

	lg.Infof(ctx, "starting server")
	go func() {
		errs <- consumer.Start(ctx)
	}()

	return grace.NewService("kafka consumer", func() {
		if err := consumer.Close(); err != nil {
			lg.Errorf(ctx, "error while shutting down http server: %v", err)
		}
	})
}
