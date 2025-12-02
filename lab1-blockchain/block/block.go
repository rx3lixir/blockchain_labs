package block

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/rx3lixir/lab1-blockchain/models"
)

// CalculateHash - вычисляет SHA3-384 хэш для блока
func CalcHash(block *models.Block) string {
	record := fmt.Sprintf(
		"%d%d%d%s%s%s%s%d%s%d",
		block.Index,
		block.Timestamp,
		block.Data.Grade,
		block.Data.FullName,
		block.Data.Zachetka,
		block.Data.Group,
		block.Data.Subject,
		block.Data.Course,
		block.PreviousHash,
		block.Nonce,
	)

	hash := sha512.Sum384([]byte(record))

	return hex.EncodeToString(hash[:])
}

// MineBlock - реализация Proof-of-Work (поиск nonce)
// Ищем такое число nonce, чтобы хэш начинался с "00"
func MineBlock(block *models.Block, difficulty string) {
	block.Nonce = 0

	for {
		block.Hash = CalcHash(block)

		if strings.HasPrefix(block.Hash, difficulty) {
			hashDisplay := block.Hash
			if len(hashDisplay) > 20 {
				hashDisplay = hashDisplay[:20]
			}
			fmt.Printf("Блок #%d намайнен! Nonce: %d, Hash: %s\n", block.Index, block.Nonce, hashDisplay)
			break
		}
		block.Nonce++
	}
}

func CreateGenesisBlock() *models.Block {
	genesisData := models.StudentRecord{
		Course:   0,
		Group:    "GENESIS",
		FullName: "GENESIS",
		Zachetka: "000000",
		Subject:  "GENESIS",
		Grade:    0,
	}

	block := &models.Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Data:         genesisData,
		PreviousHash: "0",
		Nonce:        0,
	}

	MineBlock(block, "00")
	return block
}

func PrintBlock(block *models.Block) {
	if block == nil {
		fmt.Println("Блок не найден")
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("БЛОК #%d\n", block.Index)
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("Время создания:    %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("ФИО:               %s\n", block.Data.FullName)
	fmt.Printf("№ зачётки:         %s\n", block.Data.Zachetka)
	fmt.Printf("Курс/Группа:       %d / %s\n", block.Data.Course, block.Data.Group)
	fmt.Printf("Дисциплина:        %s\n", block.Data.Subject)
	fmt.Printf("Оценка:            %d\n", block.Data.Grade)

	prevHashDisplay := block.PreviousHash
	if len(prevHashDisplay) > 20 {
		prevHashDisplay = prevHashDisplay[:20] + "..."
	}

	hashDisplay := block.Hash
	if len(hashDisplay) > 20 {
		hashDisplay = hashDisplay[:20] + "..."
	}

	fmt.Printf("Previous Hash:     %s\n", prevHashDisplay)
	fmt.Printf("Hash:              %s\n", hashDisplay)
	fmt.Printf("Nonce:             %d\n", block.Nonce)
	fmt.Println(strings.Repeat("=", 70))
}
