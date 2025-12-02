package block

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/rx3lixir/lab1-blockchain/models"
)

// CalcHash - вычисляет SHA3-384 хэш для блока (старый алгоритм)
func CalcHash(block *models.Block) string {
	return CalcHashWithAlgorithm(block, "sha3-384")
}

// CalcHashWithAlgorithm - вычисляет хэш с указанным алгоритмом
func CalcHashWithAlgorithm(block *models.Block, algorithm string) string {
	// Формируем строку для хэширования
	// Включаем Teacher если он есть (для Soft Fork совместимости)
	record := fmt.Sprintf(
		"%d%d%d%s%s%s%s%d%s%s%d",
		block.Index,
		block.Timestamp,
		block.Data.Grade,
		block.Data.FullName,
		block.Data.Zachetka,
		block.Data.Group,
		block.Data.Subject,
		block.Data.Course,
		block.Data.Teacher, // Новое поле - если пустое, просто не влияет на хэш старых блоков
		block.PreviousHash,
		block.Nonce,
	)

	var hash []byte

	switch algorithm {
	case "sha3-512":
		// Hard Fork: SHA3-512 (64 байта)
		h := sha512.Sum512([]byte(record))
		hash = h[:]
	case "sha3-384":
		fallthrough
	default:
		// Оригинал и Soft Fork: SHA3-384 (48 байт)
		h := sha512.Sum384([]byte(record))
		hash = h[:]
	}

	return hex.EncodeToString(hash)
}

// MineBlock - реализация Proof-of-Work с алгоритмом по умолчанию
func MineBlock(block *models.Block, difficulty string) {
	MineBlockWithAlgorithm(block, difficulty, "sha3-384")
}

// MineBlockWithAlgorithm - реализация Proof-of-Work с выбором алгоритма
func MineBlockWithAlgorithm(block *models.Block, difficulty string, algorithm string) {
	block.Nonce = 0
	startTime := time.Now()

	for {
		block.Hash = CalcHashWithAlgorithm(block, algorithm)

		if strings.HasPrefix(block.Hash, difficulty) {
			elapsed := time.Since(startTime)
			hashDisplay := block.Hash
			if len(hashDisplay) > 16 {
				hashDisplay = hashDisplay[:16] + "..."
			}

			fmt.Printf("✅ Block #%d mined! Nonce: %d, Hash: %s (took %v)\n",
				block.Index, block.Nonce, hashDisplay, elapsed.Round(time.Millisecond))
			break
		}
		block.Nonce++

		// Показываем прогресс каждые 10000 попыток для сложности 000
		if block.Nonce%10000 == 0 && len(difficulty) >= 3 {
			fmt.Printf("   Mining block #%d... nonce: %d\n", block.Index, block.Nonce)
		}
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
		Teacher:  "", // Пустое поле для genesis блока
	}

	block := &models.Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Data:         genesisData,
		PreviousHash: "0",
		Nonce:        0,
		ForkVersion:  models.ForkVersionOriginal,
	}

	MineBlock(block, "00") // Genesis всегда с оригинальной сложностью
	return block
}

func PrintBlock(block *models.Block) {
	if block == nil {
		fmt.Println("Блок не найден")
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 70))

	// Показываем версию форка
	forkLabel := ""
	switch block.ForkVersion {
	case models.ForkVersionOriginal:
		forkLabel = " [Original]"
	case models.ForkVersionSoft:
		forkLabel = " [Soft Fork]"
	case models.ForkVersionHard:
		forkLabel = " [Hard Fork]"
	}

	fmt.Printf("БЛОК #%d%s\n", block.Index, forkLabel)
	fmt.Println(strings.Repeat("=", 70))

	fmt.Printf("Время создания:    %s\n", time.Unix(block.Timestamp, 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("ФИО:               %s\n", block.Data.FullName)
	fmt.Printf("№ зачётки:         %s\n", block.Data.Zachetka)
	fmt.Printf("Курс/Группа:       %d / %s\n", block.Data.Course, block.Data.Group)
	fmt.Printf("Дисциплина:        %s\n", block.Data.Subject)

	// Показываем Teacher если есть (Soft Fork+)
	if block.Data.Teacher != "" {
		fmt.Printf("Преподаватель:     %s\n", block.Data.Teacher)
	} else if block.ForkVersion >= models.ForkVersionSoft {
		fmt.Printf("Преподаватель:     (не указан)\n")
	}

	fmt.Printf("Оценка:            %d\n", block.Data.Grade)

	prevHashDisplay := block.PreviousHash
	if len(prevHashDisplay) > 20 {
		prevHashDisplay = prevHashDisplay[:20] + "..."
	}

	hashDisplay := block.Hash
	hashLen := len(block.Hash)
	if hashLen > 20 {
		hashDisplay = hashDisplay[:20] + "..."
	}

	// Показываем длину хэша для Hard Fork
	hashInfo := hashDisplay
	if block.ForkVersion == models.ForkVersionHard {
		hashInfo = fmt.Sprintf("%s [SHA3-512, %d chars]", hashDisplay, hashLen)
	} else {
		hashInfo = fmt.Sprintf("%s [SHA3-384, %d chars]", hashDisplay, hashLen)
	}

	fmt.Printf("Previous Hash:     %s\n", prevHashDisplay)
	fmt.Printf("Hash:              %s\n", hashInfo)
	fmt.Printf("Nonce:             %d\n", block.Nonce)
	fmt.Println(strings.Repeat("=", 70))
}
