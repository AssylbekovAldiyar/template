package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"

	"libs/common/logger"
)

// Handler interface for processing messages
type HandleFunc func(msg *sarama.ConsumerMessage) error

type Consumer struct {
	Ready    chan bool
	Handlers map[string]HandleFunc
	client   sarama.ConsumerGroup
	admin    sarama.ClusterAdmin
	topicCfg topicConfig
}

func NewConsumer(brokers []string, groupID string, partNum int) (*Consumer, error) {
	if groupID == "" {
		return nil, fmt.Errorf("groupID is empty")
	}

	config := sarama.NewConfig()

	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		fmt.Errorf("failed to create admin client: %v", err)
	}

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("creating consumer group client: %w", err)
	}

	return &Consumer{
		Ready:    make(chan bool),
		Handlers: make(map[string]HandleFunc),
		client:   client,
		admin:    admin,
		topicCfg: topicConfig{
			partitionNumber:   partNum,
			replicationFactor: min(len(brokers), 3),
		},
	}, nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				logger.Debugf(context.Background(), "message channel was closed")
				return nil
			}
			if handler, okk := c.Handlers[message.Topic]; okk {
				if err := handler(message); err != nil {
					logger.Errorf(context.Background(), "Error handling message: %v", err)
				}
			}
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *Consumer) EnsureTopicExists(topics []string) error {
	existingTopics, err := c.admin.ListTopics()
	if err != nil {
		return fmt.Errorf("failed to list topics: %v", err)
	}

	for _, topic := range topics {
		if _, exists := existingTopics[topic]; !exists {
			topicDetail := &sarama.TopicDetail{
				NumPartitions:     int32(c.topicCfg.partitionNumber),
				ReplicationFactor: int16(c.topicCfg.replicationFactor),
			}

			err = c.admin.CreateTopic(topic, topicDetail, false)
			if err != nil {
				return fmt.Errorf("failed to create topic %s: %v", topic, err)
			}
		}
	}
	return nil
}

func (c *Consumer) Start(ctx context.Context) error {

	topics := make([]string, 0, len(c.Handlers))
	for k := range c.Handlers {
		topics = append(topics, k)
	}

	if err := c.EnsureTopicExists(topics); err != nil {
		return err
	}

	errs := make(chan error, 1)
	go func() {
		for {
			if err := c.client.Consume(ctx, topics, c); err != nil {
				errs <- err
			}
			if ctx.Err() != nil {
				errs <- ctx.Err()
			}
			c.Ready = make(chan bool)
		}
	}()

	// if consumer will not be ready in 10 seconds -> timeout
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-c.Ready:
		logger.Debugf(ctx, "Kafka consumer up and running")
	case <-ctxTimeout.Done():
		return ctxTimeout.Err()
	}

	return <-errs
}

func (c *Consumer) Close() error {
	if err := c.client.Close(); err != nil {
		return err
	}
	return c.admin.Close()
}
