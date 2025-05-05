package book

import (
	"context"
	"time"

	kafka_lib "libs/common/kafka"
	"libs/common/logger"

	"template/pkg/reqresp"
)

type KafkaClient interface {
	CreateBook(ctx context.Context, bookName string) error
}

type kafkaClient struct {
	Kafka *kafka_lib.Producer
}

func NewKafkaClient(producer *kafka_lib.Producer) KafkaClient {
	return &kafkaClient{Kafka: producer}
}

func (c *kafkaClient) CreateBook(ctx context.Context, bookName string) error {
	req := reqresp.SaveBookRequest{Name: bookName}

	logger.Infof(context.TODO(), "Processing book: %s", req.Name)

	time.Sleep(5 * time.Second)

	logger.Infof(context.TODO(), "Finished processing book: %s", req.Name)

	err := c.Kafka.SendMessage("", "some-topic", []byte("hello from kafka handler"), make(map[string]string))
	if err != nil {
		logger.Errorf(context.TODO(), "couldn't send the message to producer's buffer: %v", err)
	}

	return err
}
