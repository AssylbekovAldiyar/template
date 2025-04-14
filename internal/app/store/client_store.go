package store

import (
	"template/internal/app/connections"
	"template/internal/clients/book"
)

type ClientStore struct {
	BookHTTPClient  book.HTTPClient
	BookKafkaClient book.KafkaClient
}

func NewClientStore(conns *connections.Connections) *ClientStore {
	httpClient := book.NewHTTPClient(conns.HTTPClient)
	kafkaClient := book.NewKafkaClient(conns.Producer)

	st := &ClientStore{
		BookHTTPClient:  httpClient,
		BookKafkaClient: kafkaClient,
	}

	return st
}
