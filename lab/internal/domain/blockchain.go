package domain

import (
	"fmt"
	"strings"
	"time"
)

type Blockchain struct {
	blocks     []*Block
	forkConfig *ForkConfig
}

func NewBlockchain(config *ForkConfig) *Blockchain {
	if config == nil {
		config = &ForkConfig{
			SoftForkHeight: 5,
			HardForkHeight: 10,
			DifficultyOld:  "00",
			DifficultyNew:  "000",
		}
	}

	genesis := &Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Data:         StudentRecord{FullName: "GENESIS", Zachetka: "000000", Group: "GENESIS", Subject: "GENESIS"},
		PreviousHash: "0",
		ForkVersion:  ForkOriginal,
	}

	miner := NewMiner(&SHA384Hasher{})
	miner.Mine(genesis, "00")

	return &Blockchain{
		blocks:     []*Block{genesis},
		forkConfig: config,
	}
}

// LoadBlockchain creates a blockchain from existing blocks (for loading from storage)
func LoadBlockchain(blocks []*Block, config *ForkConfig) *Blockchain {
	if config == nil {
		config = &ForkConfig{
			SoftForkHeight: 5,
			HardForkHeight: 10,
			DifficultyOld:  "00",
			DifficultyNew:  "000",
		}
	}

	return &Blockchain{
		blocks:     blocks,
		forkConfig: config,
	}
}

// AddBlock now returns the mining duration along with any error
func (bc *Blockchain) AddBlock(data StudentRecord) (time.Duration, error) {
	if len(bc.blocks) == 0 {
		return 0, fmt.Errorf("blockchain has no blocks (corrupted)")
	}

	prevBlock := bc.blocks[len(bc.blocks)-1]
	nextIndex := prevBlock.Index + 1

	if nextIndex >= bc.forkConfig.SoftForkHeight && strings.TrimSpace(data.Teacher) == "" {
		return 0, fmt.Errorf("teacher is required. soft fork started")
	}

	forkVersion := bc.determineForkVersion(nextIndex)
	difficulty := bc.getDifficulty(nextIndex)
	hasher := bc.getHasher(nextIndex)

	newBlock := &Block{
		Index:        nextIndex,
		Timestamp:    time.Now().Unix(),
		Data:         data,
		PreviousHash: prevBlock.Hash,
		ForkVersion:  forkVersion,
	}

	miner := NewMiner(hasher)
	miningTime := miner.Mine(newBlock, difficulty)

	bc.blocks = append(bc.blocks, newBlock)
	return miningTime, nil
}

func (bc *Blockchain) Blocks() []*Block {
	return bc.blocks
}

func (bc *Blockchain) ForkConfig() *ForkConfig {
	return bc.forkConfig
}

func (bc *Blockchain) GetBlock(index int) *Block {
	for _, block := range bc.blocks {
		if block.Index == index {
			return block
		}
	}
	return nil
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
				"%s %s %s %s %s",
				block.Data.FullName, block.Data.Zachetka, block.Data.Group,
				block.Data.Subject, block.Data.Teacher,
			),
		)

		if strings.Contains(searchStr, keyword) {
			results = append(results, block)
		}
	}
	return results
}

func (bc *Blockchain) FilterByGrade(grade int) []*Block {
	var results []*Block

	for _, block := range bc.blocks {
		if block.Data.Grade == grade {
			results = append(results, block)
		}
	}

	return results
}

func (bc *Blockchain) Validate() error {
	for i := 1; i < len(bc.blocks); i++ {
		current := bc.blocks[i]
		prev := bc.blocks[i-1]

		// Check hash integrity (hard fork)
		hasher := bc.getHasher(current.Index)
		recalculated := hasher.Hash(current)
		if current.Hash != recalculated {
			return fmt.Errorf("block %d: invalid hash", current.Index)
		}

		// Check chain linkage
		if current.PreviousHash != prev.Hash {
			return fmt.Errorf("block #%d: broken chain link", current.Index)
		}

		// Check proof-of-work
		difficulty := bc.getDifficulty(current.Index)
		if !strings.HasPrefix(current.Hash, difficulty) {
			return fmt.Errorf("block %d: invalid proof-of-work", current.Index)
		}
	}
	return nil
}

func (bc *Blockchain) determineForkVersion(index int) int {
	if index >= bc.forkConfig.HardForkHeight {
		return ForkHard
	}
	if index >= bc.forkConfig.SoftForkHeight {
		return ForkSoft
	}
	return ForkOriginal
}

func (bc *Blockchain) getDifficulty(index int) string {
	if index >= bc.forkConfig.SoftForkHeight {
		return bc.forkConfig.DifficultyNew
	}
	return bc.forkConfig.DifficultyOld
}

func (bc *Blockchain) getHasher(index int) Hasher {
	if index >= bc.forkConfig.HardForkHeight {
		return &SHA512Hasher{}
	}
	return &SHA384Hasher{}
}
