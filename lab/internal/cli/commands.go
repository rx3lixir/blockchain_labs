package cli

import (
	"fmt"

	"github.com/rx3lixir/lab_bc/internal/domain"
	"github.com/rx3lixir/lab_bc/internal/storage"
)

type App struct {
	blockchain *domain.Blockchain
	storage    storage.Storage
}

func NewApp(storage storage.Storage) (*App, error) {
	bc, err := storage.Load()
	if err != nil {
		return nil, err
	}

	if bc == nil {
		bc = domain.NewBlockchain(nil)
		if err := storage.Save(bc); err != nil {
			return nil, err
		}
	}

	return &App{
		blockchain: bc,
		storage:    storage,
	}, nil
}

func (a *App) AddRecord(record domain.StudentRecord) error {
	if err := a.blockchain.AddBlock(record); err != nil {
		return err
	}
	return a.storage.Save(a.blockchain)
}

func (a *App) ListBlocks() {
	for _, block := range a.blockchain.Blocks() {
		printBlock(block)
	}
}

func (a *App) Search(query string) {
	results := a.blockchain.Search(query)
	fmt.Printf("Found %d results\n", len(results))
	for _, block := range results {
		printBlock(block)
	}
}

func (a *App) Validate() error {
	return a.blockchain.Validate()
}

func printBlock(b *domain.Block) {
	fmt.Printf("\n----- Block %d -----\n", b.Index)
	fmt.Printf("Name: %s\n", b.Data.FullName)
	fmt.Printf("Grade: %d\n", b.Data.Grade)
	fmt.Printf("Hash: %s...\n", b.Hash[:16])
	fmt.Printf("Nonce: %d\n", b.Nonce)
}
