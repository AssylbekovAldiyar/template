package kafka

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/IBM/sarama"

	"libs/common/logger"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
	wg            sync.WaitGroup
	stopChan      chan struct{}
	closed        atomic.Bool
	admin         sarama.ClusterAdmin
	topics        map[string]struct{}
	mx            sync.RWMutex
	topicCfg      topicConfig
}

func NewProducer(brokers []string, partNum int) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// Create admin client
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return nil, err
	}

	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		admin.Close()
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	producer := &Producer{
		asyncProducer: asyncProducer,
		stopChan:      make(chan struct{}),
		closed:        atomic.Bool{},
		admin:         admin,
		topics:        make(map[string]struct{}),
		topicCfg: topicConfig{
			partitionNumber:   partNum,
			replicationFactor: min(len(brokers), 3),
		},
		mx: sync.RWMutex{},
	}

	// Start monitoring routines
	producer.wg.Add(2)
	go producer.handleSuccess()
	go producer.handleError()

	return producer, nil
}

// handleSuccess monitors for successfully delivered messages
func (p *Producer) handleSuccess() {
	defer p.wg.Done()
	for {
		select {
		case success, ok := <-p.asyncProducer.Successes():
			if !ok {
				return
			}
			logger.Debugf(context.Background(), "Message delivered successfully - Topic: %s, Partition: %d, Offset: %d\n",
				success.Topic, success.Partition, success.Offset)
		case <-p.stopChan:
			return
		}
	}
}

// handleError monitors for message delivery errors
func (p *Producer) handleError() {
	defer p.wg.Done()
	for {
		select {
		case err, ok := <-p.asyncProducer.Errors():
			if !ok {
				return
			}
			logger.Errorf(context.Background(), "Failed to deliver message: %v\n", err)
		case <-p.stopChan:
			return
		}
	}
}

// SendMessage sends a message to Kafka asynchronously
func (p *Producer) SendMessage(key, topic string, value []byte, headers map[string]string) error {
	if p.closed.Load() {
		return fmt.Errorf("producer is closed")
	}

	err := p.EnsureTopicExists(topic)
	if err != nil {
		return err
	}

	saramaHeaders := make([]sarama.RecordHeader, 0, len(headers))
	for k, v := range headers {
		saramaHeaders = append(saramaHeaders, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.StringEncoder(key),
		Value:   sarama.ByteEncoder(value),
		Headers: saramaHeaders,
	}

	select {
	case p.asyncProducer.Input() <- msg:
		return nil
	case <-p.stopChan:
		return fmt.Errorf("producer is shutting down")
	}
}

func (p *Producer) EnsureTopicExists(topic string) error {

	p.mx.RLock()
	_, ok := p.topics[topic]
	p.mx.RUnlock()

	if !ok {
		topics, err := p.admin.ListTopics()
		if err != nil {
			return err
		}

		if _, okk := topics[topic]; !okk {
			topicDetail := &sarama.TopicDetail{
				NumPartitions:     int32(p.topicCfg.partitionNumber),
				ReplicationFactor: int16(p.topicCfg.replicationFactor),
			}

			err = p.admin.CreateTopic(topic, topicDetail, false)
			if err != nil {
				return err
			}
		}

		p.mx.Lock()
		p.topics[topic] = struct{}{}
		p.mx.Unlock()
	}

	return nil
}

//(TODO) В принципе, этот shutdown не критичен. Тут он нужен для того, чтобы подождать завершение горутин handleError & handleSuccess

// Shutdown gracefully shuts down the producer
func (p *Producer) Shutdown(ctx context.Context) error {
	if p.closed.Load() {
		return nil
	}
	p.closed.Store(true)

	// Signal monitoring routines to stop
	close(p.stopChan)

	// Create a channel to signal when shutdown is complete
	done := make(chan struct{})
	go func() {
		// Wait for monitoring routines to finish
		p.wg.Wait()
		// Close the async producer
		p.asyncProducer.Close()
		p.admin.Close()
		close(done)
	}()

	// Wait for shutdown to complete or context to timeout
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())
	}
}
