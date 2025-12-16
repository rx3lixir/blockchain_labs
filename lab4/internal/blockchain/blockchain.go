package blockchain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Blockchain struct {
	blocks []*Block
}

func NewBlockchain(blocks []*Block) *Blockchain {
	if len(blocks) == 0 {
		genesis := &Block{
			Index:     0,
			Timestamp: time.Now().Unix(),
			Transactions: []StudentRecord{
				{
					FullName: "GENESIS",
					Zachetka: "000000",
					Group:    "GENESIS",
					Subject:  "GENESIS",
				},
			},
			PreviousHash: "0",
			MerkleRoot:   "",
		}
		miner := NewMiner()
		miner.Mine(genesis, "00")

		blocks = []*Block{genesis}
	}

	return &Blockchain{blocks: blocks}
}

func (bc *Blockchain) AddBlock(transactions []StudentRecord) (time.Duration, error) {
	if len(bc.blocks) == 0 {
		return 0, fmt.Errorf("blockchain has no blocks (corrupted)")
	}

	if len(transactions) == 0 {
		return 0, fmt.Errorf("cannot add block without transactions")
	}

	prevBlock := bc.blocks[len(bc.blocks)-1]
	nextIndex := prevBlock.Index + 1

	for i := range transactions {
		if transactions[i].ID == "" {
			transactions[i].ID = uuid.New().String()
		}
	}

	newBlock := &Block{
		Index:        nextIndex,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PreviousHash: prevBlock.Hash,
	}

	miner := NewMiner()
	mineTime := miner.Mine(newBlock, Difficulty)

	bc.blocks = append(bc.blocks, newBlock)
	return mineTime, nil
}

func (bc *Blockchain) Blocks() []*Block {
	return bc.blocks
}

func (bc *Blockchain) Length() int {
	return len(bc.blocks)
}

func (bc *Blockchain) GetBlock(index int) (*Block, error) {
	if index < 0 || index >= len(bc.blocks) {
		return nil, fmt.Errorf("block index %d out of range", index)
	}
	return bc.blocks[index], nil
}

func (bc *Blockchain) Validate() error {
	for i := 1; i < len(bc.blocks); i++ {
		current := bc.blocks[i]
		prev := bc.blocks[i-1]

		// Проверяем хеш блока
		recalculated := CalculateHash(current)
		if current.Hash != recalculated {
			return fmt.Errorf("block %d: invalid hash", current.Index)
		}

		// Проверяем связь с предыдущим блоком
		if current.PreviousHash != prev.Hash {
			return fmt.Errorf("block #%d: broken chain link", current.Index)
		}

		// Проверяем proof-of-work
		if current.Hash[:2] != Difficulty {
			return fmt.Errorf("block %d: invalid proof-of-work", current.Index)
		}

		// Проверяем MerkleRoot
		calculatedRoot := CalculateMerkleRoot(current.Transactions)
		if current.MerkleRoot != calculatedRoot {
			return fmt.Errorf("block %d: invalid merkle root", current.Index)
		}
	}
	return nil
}
