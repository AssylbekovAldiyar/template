package pg

import (
	"context"

	"libs/common/dbctx"

	_ "github.com/lib/pq"

	"template/internal/repositories/book"
	"template/pkg/domain"
)

type repository struct {
	db dbctx.DBContext
}

func NewPgRepository(db dbctx.DBContext) book.Repository {
	return &repository{db: db}
}

func (r *repository) Save(ctx context.Context, book *domain.Book) error {
	// примеры использованиея
	//sqlx.SelectContext(ctx, r.db.Sqlx(ctx), &result, query, args...)
	//sqlx.NamedExecContext(ctx, r.db.Sqlx(ctx), query, args...)
	//sqlx.GetContext(ctx, r.db.Sqlx(ctx), &result, query, args...)
	//r.db.Sqlx(ctx).QueryxContext(ctx, query, args...)
	//r.db.Sqlx(ctx).ExecContext(ctx, query, args...)
	//r.db.Sqlx(ctx).QueryRowxContext(ctx, query, args...)

	// ошибки оборачиваются в слове репозиториев или в слое юзкейсов
	return nil
}
