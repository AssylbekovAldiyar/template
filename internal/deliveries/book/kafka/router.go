package kafka

import (
	kafka_lib "libs/common/kafka"
)

func (bh *BookHandler) RegisterHandler(consumer *kafka_lib.Consumer) {
	consumer.Handlers["test-topic"] = bh.SaveBook
}
