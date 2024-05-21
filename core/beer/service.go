package beer

import (
	_ "github.com/mattn/go-sqlite3"
)

type UseCase interface {
	GetAll() ([]*Beer, error)
	Get(ID int64) (*Beer, error)
	Store(beer *Beer) error
	Update(beer *Beer) error
	Remove(ID int64) error
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetAll() ([]*Beer, error) {
	return nil, nil
}

func (s *Service) Get(ID int) (*Beer, error) {
	return nil, nil
}

func (s *Service) Store(beer *Beer) error {
	return nil
}

func (s *Service) Update(beer *Beer) error {
	return nil
}

func (s *Service) Remove(ID int) error {
	return nil
}
