package storage

import (
	"encoding/json"
	"os"

	"github.com/rx3lixir/lab_bc/internal/domain"
)

type JSONStorage struct {
	filename string
}

func NewJSONStorage(filename string) *JSONStorage {
	return &JSONStorage{filename: filename}
}

func (s *JSONStorage) Save(bc *domain.Blockchain) error {
	data := struct {
		Blocks     []*domain.Block
		ForkConfig *domain.ForkConfig
	}{
		Blocks:     bc.Blocks(),
		ForkConfig: bc.ForkConfig(),
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, bytes, 0o644)
}

func (s *JSONStorage) Load() (*domain.Blockchain, error) {
	if !s.Exists() {
		return nil, nil
	}

	bytes, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}

	var data struct {
		Blocks     []*domain.Block
		ForkConfig *domain.ForkConfig
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	// Reconstruct blockchain
	bc := domain.NewBlockchain(data.ForkConfig)

	// Replace genesis with loaded blocks
	bc = &domain.Blockchain{}

	// Dobavi't suda SetBlocks method or use reflection
	// Poka shto tak

	return bc, nil
}

func (s *JSONStorage) Exists() bool {
	_, err := os.Stat(s.filename)
	return err == nil
}
