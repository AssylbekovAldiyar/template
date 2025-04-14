package book

import (
	"context"

	"template/pkg/domain"
)

type Repository interface {
	Save(ctx context.Context, book *domain.Book) error
}

type KafkaRepository interface {
	SendEvent(ctx context.Context, book *domain.Book) error
}
