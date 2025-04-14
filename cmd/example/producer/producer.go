package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"

	"libs/common/logger"

	"template/pkg/reqresp"
)

const (
	kafkaBrokers = "localhost:9092"
	topic        = "test-topic"
)

type Producer struct {
	asyncProducer  sarama.AsyncProducer
	partitionCount int32
}

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewManualPartitioner

	// Get partition count
	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	partitions, err := client.Partitions(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to get partitions: %w", err)
	}

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{
		asyncProducer:  producer,
		partitionCount: int32(len(partitions)),
	}, nil
}

// SendMessageToPartition sends a specific message to a specific partition
func (p *Producer) SendMessageToPartition(topic string, partition int32, message []byte) {

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.StringEncoder(message),
		Partition: partition,
	}
	p.asyncProducer.Input() <- msg
}

// SendMessagesToAllPartitions sends messages to all partitions simultaneously
func (p *Producer) SendMessagesToAllPartitions(topic string) {
	var wg sync.WaitGroup

	// Create a channel for each partition
	for partition := int32(0); partition < p.partitionCount; partition++ {
		wg.Add(1)
		go func(part int32) {
			defer wg.Done()

			book := &reqresp.SaveBookRequest{Name: fmt.Sprintf("book %d", part)}
			msg, err := json.Marshal(book)
			if err != nil {
				logger.Errorf(context.TODO(), "Failed to marshal message: %v", err)
				return
			}

			p.SendMessageToPartition(topic, part, msg)
		}(partition)

		//to imitate delayed sending of 3 messages, for graceful shutdown
		//if partition == 1 {
		//	time.Sleep(10 * time.Second)
		//}

	}

	wg.Wait()
}

// HandleResponses handles the async responses from Kafka
func (p *Producer) HandleResponses(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// Create maps to track message counts per partition
	successCounts := make(map[int32]int)
	errorCounts := make(map[int32]int)

	for {
		select {
		case success := <-p.asyncProducer.Successes():
			successCounts[success.Partition]++
			log.Printf("✅ Success - Partition: %d, Offset: %d, Message: %s\n",
				success.Partition, success.Offset, success.Value)

		case err := <-p.asyncProducer.Errors():
			errorCounts[err.Msg.Partition]++
			log.Printf("❌ Error - Partition: %d, Error: %v\n",
				err.Msg.Partition, err.Err)

		case <-ctx.Done():
			// Print final statistics
			log.Printf("\nFinal Statistics:")
			for partition := int32(0); partition < p.partitionCount; partition++ {
				log.Printf("Partition %d - Successes: %d, Errors: %d\n",
					partition,
					successCounts[partition],
					errorCounts[partition])
			}
			return
		}
	}
}

func main() {
	// Initialize producer
	brokers := strings.Split(kafkaBrokers, ",")
	producer, err := NewProducer(brokers)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Failed to create producer: s%"), err)
	}
	defer producer.asyncProducer.Close()

	log.Printf("Initialized producer with %d partitions\n", producer.partitionCount)

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// Start response handler
	wg.Add(1)
	go producer.HandleResponses(ctx, &wg)

	// Send messages to all partitions simultaneously
	producer.SendMessagesToAllPartitions(topic)

	log.Printf("All messages have been sent. Waiting for confirmations...")

	// Wait for signal to terminate
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Printf("Received termination signal. Shutting down...")
	cancel()
	wg.Wait()
}
