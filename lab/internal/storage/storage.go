package storage

import "github.com/rx3lixir/lab_bc/internal/domain"

type Storage interface {
	Save(bc *domain.Blockchain) error
	Load() (*domain.Blockchain, error)
	Exists() bool
}
