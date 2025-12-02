package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/rx3lixir/lab1-blockchain/blockchain"
	"github.com/rx3lixir/lab1-blockchain/models"
)

const blockchainCSV = "blockchain_hardfork.csv"

// SaveBlockchainCSV сохраняет блокчейн в CSV формате (Hard Fork)
func SaveBlockchainCSV(bc *blockchain.Blockchain) error {
	file, err := os.Create(blockchainCSV)
	if err != nil {
		return fmt.Errorf("error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Заголовок CSV
	header := []string{
		"Index",
		"Timestamp",
		"FullName",
		"Zachetka",
		"Group",
		"Subject",
		"Course",
		"Grade",
		"Teacher",
		"PreviousHash",
		"Hash",
		"Nonce",
		"ForkVersion",
	}

	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing CSV header: %v", err)
	}

	// Записываем блоки
	for _, block := range bc.Blocks {
		record := []string{
			strconv.Itoa(block.Index),
			strconv.FormatInt(block.Timestamp, 10),
			block.Data.FullName,
			block.Data.Zachetka,
			block.Data.Group,
			block.Data.Subject,
			strconv.Itoa(block.Data.Course),
			strconv.Itoa(block.Data.Grade),
			block.Data.Teacher,
			block.PreviousHash,
			block.Hash,
			strconv.Itoa(block.Nonce),
			strconv.Itoa(block.ForkVersion),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error writing CSV record: %v", err)
		}
	}

	fmt.Println("✅ Blockchain saved to CSV:", blockchainCSV)
	return nil
}

// LoadBlockchainCSV загружает блокчейн из CSV формата
func LoadBlockchainCSV() (*blockchain.Blockchain, error) {
	file, err := os.Open(blockchainCSV)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Файл не существует - это нормально
		}
		return nil, fmt.Errorf("error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	if len(records) < 2 { // Должен быть хотя бы заголовок + genesis
		return nil, fmt.Errorf("CSV file is empty or corrupted")
	}

	bc := &blockchain.Blockchain{
		Blocks:     make([]*models.Block, 0),
		ForkConfig: blockchain.DefaultForkConfig(),
	}

	// Пропускаем заголовок (records[0])
	for i := 1; i < len(records); i++ {
		record := records[i]

		if len(record) < 13 {
			return nil, fmt.Errorf("invalid CSV record at line %d", i+1)
		}

		index, _ := strconv.Atoi(record[0])
		timestamp, _ := strconv.ParseInt(record[1], 10, 64)
		course, _ := strconv.Atoi(record[6])
		grade, _ := strconv.Atoi(record[7])
		nonce, _ := strconv.Atoi(record[11])
		forkVersion, _ := strconv.Atoi(record[12])

		block := &models.Block{
			Index:     index,
			Timestamp: timestamp,
			Data: models.StudentRecord{
				FullName: record[2],
				Zachetka: record[3],
				Group:    record[4],
				Subject:  record[5],
				Course:   course,
				Grade:    grade,
				Teacher:  record[8],
			},
			PreviousHash: record[9],
			Hash:         record[10],
			Nonce:        nonce,
			ForkVersion:  forkVersion,
		}

		bc.Blocks = append(bc.Blocks, block)
	}

	fmt.Println("✅ Blockchain loaded from CSV:", blockchainCSV)
	return bc, nil
}

// CSVExists проверяет существование CSV файла
func CSVExists() bool {
	_, err := os.Stat(blockchainCSV)
	return !os.IsNotExist(err)
}

// ExportToCSV экспортирует существующий JSON блокчейн в CSV
func ExportToCSV(bc *blockchain.Blockchain) error {
	fmt.Println("\n🔄 Exporting blockchain to CSV format (Hard Fork)...")
	return SaveBlockchainCSV(bc)
}

// CreateHardForkBranch создаёт отдельную ветку для Hard Fork
func CreateHardForkBranch(originalBC *blockchain.Blockchain, forkHeight int) (*blockchain.Blockchain, error) {
	if forkHeight >= len(originalBC.Blocks) {
		return nil, fmt.Errorf("fork height %d exceeds blockchain length %d", forkHeight, len(originalBC.Blocks))
	}

	// Копируем блоки до точки форка
	hardForkBC := &blockchain.Blockchain{
		Blocks:     make([]*models.Block, forkHeight),
		ForkConfig: blockchain.DefaultForkConfig(),
	}

	// Копируем только блоки до Hard Fork
	for i := 0; i < forkHeight; i++ {
		// Создаём копию блока
		originalBlock := originalBC.Blocks[i]
		blockCopy := &models.Block{
			Index:        originalBlock.Index,
			Timestamp:    originalBlock.Timestamp,
			Data:         originalBlock.Data,
			PreviousHash: originalBlock.PreviousHash,
			Hash:         originalBlock.Hash,
			Nonce:        originalBlock.Nonce,
			ForkVersion:  originalBlock.ForkVersion,
		}
		hardForkBC.Blocks[i] = blockCopy
	}

	fmt.Printf("\n⚠️  Hard Fork branch created at block #%d\n", forkHeight)
	fmt.Println("   This branch is incompatible with the original chain!")
	fmt.Println("   From now on, new blocks will use:")
	fmt.Println("   - SHA3-512 hashing")
	fmt.Println("   - CSV storage format")
	fmt.Println()

	return hardForkBC, nil
}
