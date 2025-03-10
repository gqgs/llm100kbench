package service

import (
	"github.com/gqgs/llminvestbench/pkg/repository"
	"github.com/gqgs/llminvestbench/pkg/storage"
)

type Service interface {
	Exec(fn func(s repository.Repository) error) error
	ExecTx(fn func(s repository.Repository) error) error
}

type service struct {
	storage storage.Storage
}

func New(storage storage.Storage) *service {
	return &service{
		storage: storage,
	}
}

func (s *service) Exec(fn func(s repository.Repository) error) error {
	return fn(repository.New(s.storage.DB()))
}

func (s *service) ExecTx(fn func(s repository.Repository) error) error {
	tx, err := s.storage.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = fn(repository.New(tx)); err != nil {
		return err
	}

	return tx.Commit()
}
