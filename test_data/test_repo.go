package testdata

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type IUserRepo interface {
	WithTx(tx pgx.Tx) *UserRepo
	CreateUser(ctx context.Context, arg string) (string, error)
}

type UserRepo struct{}

func (u *UserRepo) WithTx(_ pgx.Tx) *UserRepo {
	return u
}
func (u *UserRepo) CreateUser(_ context.Context, _ string) (string, error) {
	return "", nil
}

type IScopeRepo interface {
	WithTx(tx pgx.Tx) *ScopeRepo
	CreateScope(ctx context.Context, arg string) (string, error)
}

type ScopeRepo struct{}

func (s *ScopeRepo) WithTx(_ pgx.Tx) *ScopeRepo {
	return s
}
func (s *ScopeRepo) CreateScope(_ context.Context, _ string) (string, error) {
	return "", nil
}
