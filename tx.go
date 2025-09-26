package sstx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ssturch/sstx/test_data"
)

type TxManager struct {
	pool testdata.IPgxPool
}

type pgxTx struct {
	pgx.Tx
}
type txUnit[T comparable] struct {
	i *T
}

func (u *txUnit[T]) Exec() *T {
	return u.i
}

type repoTx[T comparable] interface {
	WithTx(pgx.Tx) *T
}

// StopTxByErr - rollbacks or commits the transaction depending on the error
func (u pgxTx) StopTxByErr(ctx context.Context, err error) error {
	if err != nil {
		return u.Rollback(ctx)
	}
	return u.Commit(ctx)
}

// Commit - commits the transaction
func (u pgxTx) Commit(ctx context.Context) error {
	return u.Tx.Commit(ctx)
}

// Rollback - rollback the transaction
func (u pgxTx) Rollback(ctx context.Context) error {
	return u.Tx.Rollback(ctx)
}

// New - init new TxManager
func New(pool testdata.IPgxPool) *TxManager {
	return &TxManager{
		pool: pool,
	}
}

// GetTx - begin new transaction by TxManager
func (m *TxManager) GetTx(ctx context.Context) (pgxTx, error) {
	tx, err := m.pool.Begin(ctx)
	if err != nil {
		return pgxTx{}, err
	}
	return pgxTx{tx}, nil
}

// GetTxWithOpt - begin new transaction with options by TxManager
func (m *TxManager) GetTxWithOpt(ctx context.Context, opt pgx.TxOptions) (pgxTx, error) {
	tx, err := m.pool.BeginTx(ctx, opt)
	if err != nil {
		return pgxTx{}, err
	}
	return pgxTx{tx}, nil
}

// NewUnit - Create new transactional unit by interface
func NewUnit[T comparable](tx pgxTx, repo repoTx[T]) *txUnit[T] {
	return &txUnit[T]{i: repo.WithTx(tx)}
}
