package dbctx

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type contextKey string

const (
	sqlxTxKey contextKey = "dbctxSqlxTxKey"
)

var ContextTransactionEmptyErr = errors.New("context transaction is empty")

type SqlxBuilder interface {
	sqlx.Ext
	sqlx.ExtContext
}

type DBContext interface {
	Sqlx(ctx context.Context) SqlxBuilder
	SqlxBegin(ctx context.Context) (context.Context, error)
	SqlxCommit(ctx context.Context) error
	SqlxRollback(ctx context.Context) error
}

func New(db *sqlx.DB) DBContext {
	return &dbContext{
		dbSqlx: db,
	}
}

type dbContext struct {
	dbSqlx *sqlx.DB
}

func (d *dbContext) Sqlx(ctx context.Context) SqlxBuilder {
	if tx, ok := ctx.Value(sqlxTxKey).(*sqlx.Tx); ok {
		return tx
	}

	return d.dbSqlx
}

func (d *dbContext) SqlxBegin(ctx context.Context) (context.Context, error) {
	tx, err := d.dbSqlx.BeginTxx(ctx, nil)
	if err != nil {
		return ctx, err
	}

	return context.WithValue(ctx, sqlxTxKey, tx), nil
}

func (d *dbContext) SqlxCommit(ctx context.Context) error {
	if tx, ok := ctx.Value(sqlxTxKey).(*sqlx.Tx); ok {
		return tx.Commit()
	}

	return ContextTransactionEmptyErr
}

func (d *dbContext) SqlxRollback(ctx context.Context) error {
	if tx, ok := ctx.Value(sqlxTxKey).(*sqlx.Tx); ok {
		return tx.Rollback()
	}

	return ContextTransactionEmptyErr
}
