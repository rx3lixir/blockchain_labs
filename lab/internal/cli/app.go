package cli

import (
	"github.com/rx3lixir/lab_bc/internal/blockchain"
	"github.com/rx3lixir/lab_bc/internal/storage"
)

type App struct {
	bc      *blockchain.Blockchain
	storage storage.Storage
}

func NewApp(storage storage.Storage) (*App, error) {
	bc, err := storage.Load()
	if err != nil {
		return nil, err
	}

	if bc == nil {
		bc = blockchain.NewBlockchain(nil)
		if err := storage.Save(bc); err != nil {
			return nil, err
		}
	}

	return &App{
		bc:      bc,
		storage: storage,
	}, nil
}
