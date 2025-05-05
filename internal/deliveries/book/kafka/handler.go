package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"libs/common/ctxconst"

	"template/internal/services/book"

	"github.com/IBM/sarama"

	kafka_lib "libs/common/kafka"
	"libs/common/logger"

	"template/pkg/reqresp"
)

type BookHandler struct {
	// Add any dependencies here (services, repositories, etc.)
	Producer *kafka_lib.Producer
	service  book.Service
	timeout  time.Duration
}

func NewBookHandler(service book.Service, producer *kafka_lib.Producer, timeoutSeconds int) *BookHandler {

	if timeoutSeconds <= 0 {
		timeoutSeconds = 60
	}

	return &BookHandler{
		service:  service,
		Producer: producer,
		timeout:  time.Second * time.Duration(timeoutSeconds),
	}
}

func (bh *BookHandler) SaveBook(msg *sarama.ConsumerMessage) error {

	ctx, cancel := bh.context()
	defer cancel()

	var bookReq reqresp.SaveBookRequest
	if err := json.Unmarshal(msg.Value, &bookReq); err != nil {
		return fmt.Errorf("failed to unmarshal book request: %w", err)
	}

	logger.Infof(ctx, "Processing book: %s", bookReq.Name)

	time.Sleep(15 * time.Second)

	logger.Infof(ctx, "Finished processing book: %s", bookReq.Name)

	err := bh.Producer.SendMessage("", "some-topic", []byte("hello from kafka handler"), make(map[string]string))
	if err != nil {
		logger.Errorf(ctx, "couldn't send the message to producer's buffer: %v", err)
	}

	return nil
}

func (bh *BookHandler) context() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	// todo: добавить мета инфу (реквест айди, пользователь и тд)
	ctx = ctxconst.SetRequestID(ctx, "test-request-id")
	ctx = ctxconst.SetUserID(ctx, "test-user-id")

	return context.WithTimeout(ctx, bh.timeout)
}
