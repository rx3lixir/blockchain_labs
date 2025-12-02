package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rx3lixir/lab1-blockchain/blockchain"
)

const blockchainFile = "blockchain.json"

func SaveBlockchain(bc *blockchain.Blockchain) error {
	data, err := json.MarshalIndent(bc, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %v", err)
	}

	err = os.WriteFile(blockchainFile, data, 0o644)
	if err != nil {
		return fmt.Errorf("ошибка записи файла: %v", err)
	}

	fmt.Println("Блокчейн сохранен в...", blockchainFile)
	return nil
}

func LoadBlockchain() (*blockchain.Blockchain, error) {
	data, err := os.ReadFile(blockchainFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	var bc blockchain.Blockchain
	err = json.Unmarshal(data, &bc)
	if err != nil {
		return nil, fmt.Errorf("ошибка десериализации: %v", err)
	}

	fmt.Println("Блокчейн загружен из", blockchainFile)
	return &bc, nil
}

func BlockchainExists() bool {
	_, err := os.Stat(blockchainFile)
	return !os.IsNotExist(err)
}
