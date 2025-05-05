package book

import (
	"context"
	"libs/apperror"

	"template/internal/app/store"
	"template/pkg/domain"
	"template/pkg/reqresp"
)

type SaveBookRepository interface {
	SaveBook(ctx context.Context, book *domain.Book) error
}

func SaveBook(
	ctx context.Context,
	repo SaveBookRepository,
	req reqresp.SaveBookRequest,
) (reqresp.SaveBookResponse, error) {
	book := domain.Book{
		ID:   0,
		Name: req.Name,
	}

	err := repo.SaveBook(ctx, &book)
	if err != nil {
		return reqresp.SaveBookResponse{}, err
	}

	return reqresp.SaveBookResponse{
		Name: book.Name,
	}, nil
}

type saveBookRepositoryFacade struct {
	st      *store.RepositoryStore
	clients *store.ClientStore
}

func NewSaveBookRepository(resources *store.RepositoryStore, clients *store.ClientStore) SaveBookRepository {
	return &saveBookRepositoryFacade{st: resources, clients: clients}
}

func (f *saveBookRepositoryFacade) SaveBook(ctx context.Context, book *domain.Book) error {
	// здесь можно инкапсулировать несколько запросов в разные репозитории
	err := f.st.BookRepository.Save(ctx, book)
	if err != nil {
		return err
	}

	err = f.clients.BookKafkaClient.CreateBook(ctx, book.Name)

	if err != nil {
		return apperror.AsError(err)
	}

	return nil
}
