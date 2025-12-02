package blockchain

import (
	"fmt"
	"strings"
	"time"

	"github.com/rx3lixir/lab1-blockchain/block"
	"github.com/rx3lixir/lab1-blockchain/models"
)

type Blockchain struct {
	Blocks []*models.Block
}

// New creates a new blockchain with genesis block
func New() *Blockchain {
	return &Blockchain{
		Blocks: []*models.Block{block.CreateGenesisBlock()},
	}
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data models.StudentRecord) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	newBlock := &models.Block{
		Index:        prevBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
		Data:         data,
		PreviousHash: prevBlock.Hash,
		Nonce:        0,
	}

	block.MineBlock(newBlock, "00")
	bc.Blocks = append(bc.Blocks, newBlock)
}

// GetBlockByIndex returns a block by its index
func (bc *Blockchain) GetBlockByIndex(index int) *models.Block {
	for _, blk := range bc.Blocks {
		if blk.Index == index {
			return blk
		}
	}
	return nil
}

// SearchByKeyword searches blocks by keyword in multiple fields
// This is the "smart search" - it searches across all text fields
func (bc *Blockchain) SearchByKeyword(keyword string) []*models.Block {
	var result []*models.Block
	keyword = strings.ToLower(strings.TrimSpace(keyword))

	if keyword == "" {
		return result
	}

	for _, blk := range bc.Blocks {
		// Build searchable string from all text fields
		searchStr := strings.ToLower(fmt.Sprintf("%s %s %s %s",
			blk.Data.FullName,
			blk.Data.Zachetka,
			blk.Data.Group,
			blk.Data.Subject,
		))

		if strings.Contains(searchStr, keyword) {
			result = append(result, blk)
		}
	}

	return result
}

// FilterByGrade returns all blocks with a specific grade
func (bc *Blockchain) FilterByGrade(grade int) []*models.Block {
	var result []*models.Block

	for _, blk := range bc.Blocks {
		if blk.Data.Grade == grade {
			result = append(result, blk)
		}
	}

	return result
}

// FilterByCourse returns all blocks from a specific course
func (bc *Blockchain) FilterByCourse(course int) []*models.Block {
	var result []*models.Block

	for _, blk := range bc.Blocks {
		if blk.Data.Course == course {
			result = append(result, blk)
		}
	}

	return result
}

// GetBlocksByTimeRange returns blocks within a time range
func (bc *Blockchain) GetBlocksByTimeRange(startTime, endTime int64) []*models.Block {
	var result []*models.Block

	for _, blk := range bc.Blocks {
		if blk.Timestamp >= startTime && blk.Timestamp <= endTime {
			result = append(result, blk)
		}
	}

	return result
}

// GetStatistics returns blockchain statistics
func (bc *Blockchain) GetStatistics() map[string]any {
	stats := make(map[string]any)

	stats["total_blocks"] = len(bc.Blocks)

	gradeCount := make(map[int]int)
	courseCount := make(map[int]int)

	for _, blk := range bc.Blocks {
		if blk.Index == 0 { // Skip genesis block
			continue
		}
		gradeCount[blk.Data.Grade]++
		courseCount[blk.Data.Course]++
	}

	stats["grades"] = gradeCount
	stats["courses"] = courseCount

	return stats
}

// ValidateChain validates the entire blockchain
func (bc *Blockchain) ValidateChain() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		// Check if hash is correct
		if currentBlock.Hash != block.CalcHash(currentBlock) {
			fmt.Printf("Block #%d has invalid hash\n", currentBlock.Index)
			return false
		}

		// Check if previous hash matches
		if currentBlock.PreviousHash != prevBlock.Hash {
			fmt.Printf("Block #%d is not linked to previous block\n", currentBlock.Index)
			return false
		}

		// Check proof-of-work
		if !strings.HasPrefix(currentBlock.Hash, "00") {
			fmt.Printf("Block #%d failed Proof-of-Work validation\n", currentBlock.Index)
			return false
		}
	}

	return true
}
