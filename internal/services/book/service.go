package book

import (
	"context"

	"libs/common/ctxconst"
	"libs/common/logger"
	"libs/common/service_middleware"

	"template/internal/app/store"
	"template/internal/usecases/book"
	"template/pkg/reqresp"
)

const (
	// Methods
	saveBook       = "book.SaveBook"
	saveBookNoResp = "book.SaveBookWithNoResponse"
)

// Service определяет интерфейс сервиса для работы с книгами
type Service interface {
	SaveBook(ctx context.Context, request reqresp.SaveBookRequest) (reqresp.SaveBookResponse, error)
	SaveBookWithNoResponse(ctx context.Context, request reqresp.SaveBookRequest) error
}

// service реализует интерфейс Service
type service struct {
	st *store.RepositoryStore
	cl *store.ClientStore
	mw *service_middleware.ServiceMiddleware
}

// New создает новый экземпляр сервиса с глобальными middleware и специфичными для метода
func New(st *store.RepositoryStore, cl *store.ClientStore, mw *service_middleware.ServiceMiddleware) Service {
	// Создаем middleware для транзакций
	txMiddleware := func(next service_middleware.HandlerFunc) service_middleware.HandlerFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			ctx, err = st.PgContext.SqlxBegin(ctx)
			if err != nil {
				return nil, err
			}

			defer func() {
				if err != nil {
					rErr := st.PgContext.SqlxRollback(ctx)
					if rErr != nil {
						logger.Errorf(ctx, "failed to rollback tx: %v", rErr)
					}
				}
			}()

			resp, err = next(ctx, req)
			if err != nil {
				return resp, err
			}

			err = st.PgContext.SqlxCommit(ctx)
			if err != nil {
				logger.Errorf(ctx, "failed to commit tx: %v", err)
			}

			return resp, err
		}
	}

	// Регистрируем дополнительные middleware для конкретных методов
	mw.UseForMethod(saveBook, txMiddleware)

	return &service{
		st: st,
		cl: cl,
		mw: mw,
	}
}

// SaveBook сохраняет книгу
func (s *service) SaveBook(ctx context.Context, req reqresp.SaveBookRequest) (reqresp.SaveBookResponse, error) {
	ctx = ctxconst.SetMethodName(ctx, saveBook)

	handler := s.mw.Wrap(saveBook, func(ctx context.Context, req interface{}) (interface{}, error) {
		typedReq := req.(reqresp.SaveBookRequest)
		return book.SaveBook(ctx, book.NewSaveBookRepository(s.st, s.cl), typedReq)
	})

	resp, err := handler(ctx, req)
	typedResp, err := service_middleware.SafeCast[reqresp.SaveBookResponse](resp, err)
	if err != nil {
		return reqresp.SaveBookResponse{}, err
	}

	return typedResp, nil
}

// SaveBookWithNoResponse сохраняет книгу без возврата результата
func (s *service) SaveBookWithNoResponse(ctx context.Context, req reqresp.SaveBookRequest) error {
	ctx = ctxconst.SetMethodName(ctx, saveBookNoResp)

	handler := s.mw.Wrap(saveBookNoResp, func(ctx context.Context, req interface{}) (interface{}, error) {
		typedReq := req.(reqresp.SaveBookRequest)

		var resp reqresp.SaveBookResponse
		resp, useCaseErr := book.SaveBook(ctx, book.NewSaveBookRepository(s.st, s.cl), typedReq)

		return resp, useCaseErr
	})

	_, err := handler(ctx, req)

	return err
}
