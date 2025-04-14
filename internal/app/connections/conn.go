package connections

import (
	"context"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	kafka_lib "libs/common/kafka"
	"libs/common/logger"

	"template/internal/app/config"
)

type Connections struct {
	DB *sqlx.DB

	HTTPClient *http.Client

	Producer *kafka_lib.Producer // нужно закрывать в Close()
	Consumer *kafka_lib.Consumer // не нужно закрывать в Close()
}

func (c *Connections) Close() {
	if c.DB != nil {
		_ = c.DB.Close()
	}
	if c.Producer != nil {
		_ = c.Producer.Shutdown(context.Background())
	}
}

func New(cfg *config.Config) (*Connections, error) {
	httpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    cfg.HTTPClient.MaxIdleConnections,
			MaxConnsPerHost: cfg.HTTPClient.MaxConnsPerHost,
			IdleConnTimeout: time.Duration(cfg.HTTPClient.IdleConnTimeout) * time.Millisecond,
		},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Duration(cfg.HTTPClient.Timeout) * time.Millisecond,
	}

	producer, err := kafka_lib.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.PartitionNumber)
	if err != nil {
		logger.Fatalf(context.TODO(), "can't create kafka producer: %v", err)
	}

	consumer, err := kafka_lib.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.ConsumerGroup, cfg.Kafka.PartitionNumber)
	if err != nil {
		logger.Fatalf(context.TODO(), "can't create kafka consumer: %v", err)
	}

	return &Connections{
		HTTPClient: &httpClient,
		Producer:   producer,
		Consumer:   consumer,
	}, nil

}
