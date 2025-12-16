package storage

import "github.com/rx3lixir/lab_bc/internal/blockchain"

type Storage interface {
	Save(bc *blockchain.Blockchain) error
	Load() (*blockchain.Blockchain, error)
	Exists() bool
}
