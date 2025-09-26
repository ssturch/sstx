package sstx

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/ssturch/sstx/test_data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var ErrTest = errors.New("test error")

func TestNew(t *testing.T) {
	t.Run("should create a new instance of ITx", func(t *testing.T) {
		mpool := testdata.NewMockIPgxPool(t)
		txManager := New(mpool)
		assert.NotNil(t, txManager)
	})
}
func TestTxManager_GetTx(t *testing.T) {
	t.Run("should get a tx from the pool", func(t *testing.T) {
		ctx := context.Background()
		mpool := testdata.NewMockIPgxPool(t)
		mtx := testdata.NewMockTx(t)
		txManager := New(mpool)

		require.NotNil(t, txManager)

		want := pgxTx{mtx}
		mock.InOrder(
			mpool.EXPECT().Begin(ctx).Return(mtx, nil).Once(),
		)
		got, err := txManager.GetTx(ctx)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("get a tx error", func(t *testing.T) {
		ctx := context.Background()
		mpool := testdata.NewMockIPgxPool(t)
		txManager := New(mpool)

		require.NotNil(t, txManager)

		want := pgxTx{}
		mock.InOrder(
			mpool.EXPECT().Begin(ctx).Return(nil, ErrTest).Once(),
		)
		got, err := txManager.GetTx(ctx)
		assert.Error(t, err)
		assert.Equal(t, ErrTest, err)
		assert.Equal(t, want, got)
	})
}
func TestTxManager_GetTxWithOpt(t *testing.T) {
	t.Run("should get a tx from the pool", func(t *testing.T) {
		ctx := context.Background()
		mpool := testdata.NewMockIPgxPool(t)
		mtx := testdata.NewMockTx(t)
		txManager := New(mpool)

		require.NotNil(t, txManager)

		want := pgxTx{mtx}
		opt := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

		mock.InOrder(
			mpool.EXPECT().BeginTx(ctx, opt).Return(mtx, nil).Once(),
		)

		got, err := txManager.GetTxWithOpt(ctx, opt)
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
	t.Run("get a tx error", func(t *testing.T) {
		ctx := context.Background()
		mpool := testdata.NewMockIPgxPool(t)
		txManager := New(mpool)

		require.NotNil(t, txManager)

		want := pgxTx{}
		opt := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}

		mock.InOrder(
			mpool.EXPECT().BeginTx(ctx, opt).Return(nil, ErrTest).Once(),
		)
		got, err := txManager.GetTxWithOpt(ctx, opt)
		assert.Error(t, err)
		assert.Equal(t, ErrTest, err)
		assert.Equal(t, want, got)
	})
}

func TestNewUnit(t *testing.T) {
	t.Run("should create a new repo unit", func(t *testing.T) {
		mtx := testdata.NewMockTx(t)
		repo := new(testdata.UserRepo)
		repoWithTx := repo.WithTx(mtx)

		want := &txUnit[testdata.UserRepo]{i: repoWithTx}
		got := NewUnit[testdata.UserRepo](pgxTx{mtx}, repo)

		assert.Equal(t, want, got)
	})
}
func TestTxManager_StopTxByErr(t *testing.T) {
	t.Run("stop and commit", func(t *testing.T) {
		ctx := context.Background()
		mpool := testdata.NewMockIPgxPool(t)
		mtx := testdata.NewMockTx(t)

		txManager := New(mpool)
		require.NotNil(t, txManager)

		mock.InOrder(
			mpool.EXPECT().Begin(ctx).Return(mtx, nil).Once(),
			mtx.EXPECT().Commit(ctx).Return(nil).Once(),
		)

		tx, err := txManager.GetTx(ctx)
		require.NoError(t, err)

		err = tx.StopTxByErr(ctx, nil)
		require.NoError(t, err)
	})
	t.Run("stop and rollback", func(t *testing.T) {
		ctx := context.Background()
		mpool := testdata.NewMockIPgxPool(t)
		mtx := testdata.NewMockTx(t)

		txManager := New(mpool)
		require.NotNil(t, txManager)

		mock.InOrder(
			mpool.EXPECT().Begin(ctx).Return(mtx, nil).Once(),
			mtx.EXPECT().Rollback(ctx).Return(nil).Once(),
		)

		tx, err := txManager.GetTx(ctx)
		require.NoError(t, err)

		err = tx.StopTxByErr(ctx, ErrTest)
		require.NoError(t, err)
	})
}
