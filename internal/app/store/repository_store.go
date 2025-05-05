package store

import (
	"libs/common/dbctx"

	"template/internal/app/connections"
	"template/internal/repositories/book"
	bookpg "template/internal/repositories/book/pg"
)

type RepositoryStore struct {
	PgContext dbctx.DBContext

	BookRepository book.Repository
}

func NewRepositoryStore(conns *connections.Connections) *RepositoryStore {
	st := &RepositoryStore{
		PgContext: dbctx.New(conns.DB),
	}

	st.BookRepository = bookpg.NewPgRepository(st.PgContext)

	return st
}
