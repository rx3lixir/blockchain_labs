package storage

import (
	"encoding/json"
	"os"

	"github.com/rx3lixir/lab_bc/internal/blockchain"
)

type JSONStorage struct {
	filename string
}

func NewJSONStorage(filename string) *JSONStorage {
	return &JSONStorage{filename: filename}
}

func (s *JSONStorage) Save(bc *blockchain.Blockchain) error {
	data := struct {
		Blocks []*blockchain.Block `json:"blocks"`
	}{
		Blocks: bc.Blocks(),
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filename, bytes, 0o644)
}

func (s *JSONStorage) Load() (*blockchain.Blockchain, error) {
	if !s.Exists() {
		return nil, nil
	}

	bytes, err := os.ReadFile(s.filename)
	if err != nil {
		return nil, err
	}

	var data struct {
		Blocks []*blockchain.Block
	}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	return blockchain.NewBlockchain(data.Blocks), nil
}

func (s *JSONStorage) Exists() bool {
	_, err := os.Stat(s.filename)
	return err == nil
}
