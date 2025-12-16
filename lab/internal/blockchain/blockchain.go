package blockchain

import (
	"fmt"
	"strings"
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
			Data: StudentRecord{
				FullName: "GENESIS",
				Zachetka: "000000",
				Group:    "GENESIS",
				Subject:  "GENESIS",
			},
			PreviousHash: "0",
		}
		miner := NewMiner()
		miner.Mine(genesis, "00")

		blocks = []*Block{genesis}
	}

	return &Blockchain{blocks: blocks}
}

func (bc *Blockchain) AddBlock(data StudentRecord) (time.Duration, error) {
	if len(bc.blocks) == 0 {
		return 0, fmt.Errorf("blockchain has no blocks (corrupted)")
	}

	prevBlock := bc.blocks[len(bc.blocks)-1]
	nextIndex := prevBlock.Index + 1

	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	newBlock := &Block{
		Index:        nextIndex,
		Timestamp:    time.Now().Unix(),
		Data:         data,
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

func (bc *Blockchain) Search(keyword string) []*Block {
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return nil
	}

	var results []*Block

	for _, block := range bc.blocks {
		searchStr := strings.ToLower(
			fmt.Sprintf(
				"%s %s %s %s",
				block.Data.FullName, block.Data.Zachetka,
				block.Data.Group, block.Data.Subject,
			),
		)

		if strings.Contains(searchStr, keyword) {
			results = append(results, block)
		}
	}
	return results
}

func (bc *Blockchain) Validate() error {
	for i := 1; i < len(bc.blocks); i++ {
		current := bc.blocks[i]
		prev := bc.blocks[i-1]

		recalculated := CalculateHash(current)

		if current.Hash != recalculated {
			return fmt.Errorf("block %d: invalid hash", current.Index)
		}

		// Check chain linkage
		if current.PreviousHash != prev.Hash {
			return fmt.Errorf("block #%d: broken chain link", current.Index)
		}

		if !strings.HasPrefix(current.Hash, Difficulty) {
			return fmt.Errorf("block %d: invalid proof-of-work", current.Index)
		}
	}
	return nil
}
