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
		return fmt.Errorf("error making marshal indent: %v", err)
	}

	err = os.WriteFile(blockchainFile, data, 0o644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	fmt.Println("✅ Blockchain saved to:", blockchainFile)
	return nil
}

func LoadBlockchain() (*blockchain.Blockchain, error) {
	data, err := os.ReadFile(blockchainFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var bc blockchain.Blockchain
	err = json.Unmarshal(data, &bc)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling: %v", err)
	}

	// Если ForkConfig не был загружен (старый формат), создаём дефолтный
	if bc.ForkConfig == nil {
		bc.ForkConfig = blockchain.DefaultForkConfig()
	}

	fmt.Println("✅ Blockchain loaded from:", blockchainFile)
	return &bc, nil
}

func BlockchainExists() bool {
	_, err := os.Stat(blockchainFile)
	return !os.IsNotExist(err)
}
