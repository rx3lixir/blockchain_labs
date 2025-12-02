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

func New() *Blockchain {
	return &Blockchain{
		Blocks: []*models.Block{block.CreateGenesisBlock()},
	}
}

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

func (bc *Blockchain) GetBlockByIndex(index int) *models.Block {
	for _, blk := range bc.Blocks {
		if blk.Index == index {
			return blk
		}
	}
	return nil
}

func (bc *Blockchain) GetBlocksByTimeRange(startTime, endTime int64) []*models.Block {
	var result []*models.Block

	for _, blk := range bc.Blocks {
		if blk.Timestamp >= startTime && blk.Timestamp <= endTime {
			result = append(result, blk)
		}
	}

	return result
}

func (bc *Blockchain) SearchByKeyword(keyword string) []*models.Block {
	var result []*models.Block
	keyword = strings.ToLower(keyword)

	for _, blk := range bc.Blocks {
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

func (bc *Blockchain) FindStudentsWithDebts() []*models.Block {
	var result []*models.Block

	for _, blk := range bc.Blocks {
		if blk.Data.Grade == 2 {
			result = append(result, blk)
		}
	}

	return result
}

func (bc *Blockchain) ValidateChain() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		if currentBlock.Hash != block.CalcHash(currentBlock) {
			fmt.Printf("Блок #%d имеет неверный хэш\n", currentBlock.Index)
			return false
		}

		if currentBlock.PreviousHash != prevBlock.Hash {
			fmt.Printf("Блок #%d не связа с предыдущим\n", currentBlock.Index)
			return false
		}

		if !strings.HasPrefix(currentBlock.Hash, "00") {
			fmt.Printf("Блок #%d не прошел Proof-of-Work\n", currentBlock.Index)
			return false
		}
	}

	fmt.Println("Блокчейн прошел валидацию")
	return true
}

func (bc *Blockchain) PrintAllBlocks() {
	fmt.Printf("\nБЛОКЧЕЙН ЖУРНАЛ УСПЕВАЕМОСТИ (всего блоков: %d)\n", len(bc.Blocks))
	for _, blk := range bc.Blocks {
		block.PrintBlock(blk)
	}
}
