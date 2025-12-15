package cli

import (
	"fmt"
	"time"

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

// AddRecord now returns mining duration
func (a *App) AddRecord(record domain.StudentRecord) (time.Duration, error) {
	miningTime, err := a.blockchain.AddBlock(record)
	if err != nil {
		return 0, err
	}

	if err := a.storage.Save(a.blockchain); err != nil {
		return 0, err
	}

	return miningTime, nil
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
