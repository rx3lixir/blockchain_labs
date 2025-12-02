package blockchain

import (
	"fmt"
	"strings"
	"time"

	"github.com/rx3lixir/lab1-blockchain/block"
	"github.com/rx3lixir/lab1-blockchain/models"
)

type Blockchain struct {
	Blocks     []*models.Block
	ForkConfig *ForkConfig
}

// ForkConfig хранит конфигурацию форков
type ForkConfig struct {
	SoftForkHeight int    // Блок активации Soft Fork
	HardForkHeight int    // Блок активации Hard Fork
	DifficultyOld  string // "00"
	DifficultyNew  string // "000"
	HashAlgorithm  string // "sha3-384" или "sha3-512"
}

// DefaultForkConfig возвращает стандартную конфигурацию
func DefaultForkConfig() *ForkConfig {
	return &ForkConfig{
		SoftForkHeight: 5,
		HardForkHeight: 10,
		DifficultyOld:  "00",
		DifficultyNew:  "000",
		HashAlgorithm:  "sha3-384",
	}
}

// New creates a new blockchain with genesis block
func New() *Blockchain {
	return &Blockchain{
		Blocks:     []*models.Block{block.CreateGenesisBlock()},
		ForkConfig: DefaultForkConfig(),
	}
}

// GetDifficulty возвращает сложность для следующего блока
func (bc *Blockchain) GetDifficulty() string {
	nextIndex := len(bc.Blocks)

	if nextIndex >= bc.ForkConfig.SoftForkHeight {
		return bc.ForkConfig.DifficultyNew
	}
	return bc.ForkConfig.DifficultyOld
}

// GetHashAlgorithm возвращает алгоритм хэширования для следующего блока
func (bc *Blockchain) GetHashAlgorithm() string {
	nextIndex := len(bc.Blocks)

	if nextIndex >= bc.ForkConfig.HardForkHeight {
		return "sha3-512"
	}
	return "sha3-384"
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(data models.StudentRecord) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	nextIndex := prevBlock.Index + 1

	// Определяем версию форка для нового блока
	forkVersion := models.ForkVersionOriginal
	if nextIndex >= bc.ForkConfig.HardForkHeight {
		forkVersion = models.ForkVersionHard
	} else if nextIndex >= bc.ForkConfig.SoftForkHeight {
		forkVersion = models.ForkVersionSoft
	}

	// Для Soft Fork - если поле Teacher пустое, предупреждаем
	if forkVersion >= models.ForkVersionSoft && data.Teacher == "" {
		fmt.Println("⚠️  Warning: Teacher field is empty (Soft Fork active)")
	}

	newBlock := &models.Block{
		Index:        nextIndex,
		Timestamp:    time.Now().Unix(),
		Data:         data,
		PreviousHash: prevBlock.Hash,
		Nonce:        0,
		ForkVersion:  forkVersion,
	}

	// Получаем текущие параметры
	difficulty := bc.GetDifficulty()
	hashAlg := bc.GetHashAlgorithm()

	fmt.Printf("Mining block #%d (Fork v%d, Difficulty: %s, Hash: %s)\n",
		nextIndex, forkVersion, difficulty, hashAlg)

	// Майним блок с нужным алгоритмом
	block.MineBlockWithAlgorithm(newBlock, difficulty, hashAlg)
	bc.Blocks = append(bc.Blocks, newBlock)

	// Выводим информацию о форке если он только активировался
	if nextIndex == bc.ForkConfig.SoftForkHeight {
		fmt.Println("\n🔄 SOFT FORK ACTIVATED!")
		fmt.Println("   → Difficulty: 00 → 000")
		fmt.Println("   → New field: Teacher (optional for compatibility)")
		fmt.Println()
	}

	if nextIndex == bc.ForkConfig.HardForkHeight {
		fmt.Println("\n⚠️  HARD FORK ACTIVATED!")
		fmt.Println("   → Hash: SHA3-384 → SHA3-512")
		fmt.Println("   → Storage: JSON → CSV")
		fmt.Println("   → NOT backwards compatible!")
		fmt.Println()
	}
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
func (bc *Blockchain) SearchByKeyword(keyword string) []*models.Block {
	var result []*models.Block
	keyword = strings.ToLower(strings.TrimSpace(keyword))

	if keyword == "" {
		return result
	}

	for _, blk := range bc.Blocks {
		// Build searchable string from all text fields including Teacher
		searchStr := strings.ToLower(fmt.Sprintf("%s %s %s %s %s",
			blk.Data.FullName,
			blk.Data.Zachetka,
			blk.Data.Group,
			blk.Data.Subject,
			blk.Data.Teacher,
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

// ValidateChain validates the entire blockchain with fork awareness
func (bc *Blockchain) ValidateChain() bool {
	fmt.Println("\n=== Validating blockchain with fork rules ===")

	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		// Определяем ожидаемую сложность для этого блока
		expectedDifficulty := bc.ForkConfig.DifficultyOld
		if currentBlock.Index >= bc.ForkConfig.SoftForkHeight {
			expectedDifficulty = bc.ForkConfig.DifficultyNew
		}

		// Определяем ожидаемый алгоритм хэширования
		expectedAlg := "sha3-384"
		if currentBlock.Index >= bc.ForkConfig.HardForkHeight {
			expectedAlg = "sha3-512"
		}

		fmt.Printf("Block #%d: expected difficulty=%s, algorithm=%s\n",
			currentBlock.Index, expectedDifficulty, expectedAlg)

		// Проверяем хэш с правильным алгоритмом
		recalculated := block.CalcHashWithAlgorithm(currentBlock, expectedAlg)
		if currentBlock.Hash != recalculated {
			fmt.Printf("❌ Block #%d has invalid hash\n", currentBlock.Index)
			return false
		}

		// Проверяем связь с предыдущим блоком
		if currentBlock.PreviousHash != prevBlock.Hash {
			fmt.Printf("❌ Block #%d is not linked to previous block\n", currentBlock.Index)
			return false
		}

		// Проверяем proof-of-work (сложность)
		if !strings.HasPrefix(currentBlock.Hash, expectedDifficulty) {
			fmt.Printf("❌ Block #%d failed Proof-of-Work validation (expected: %s)\n",
				currentBlock.Index, expectedDifficulty)
			return false
		}

		// Проверка Soft Fork: с блока #5 должно быть поле Teacher (но может быть пустым)
		if currentBlock.Index >= bc.ForkConfig.SoftForkHeight {
			// Это просто информация, не ошибка (обратная совместимость)
			if currentBlock.Data.Teacher == "" {
				fmt.Printf("ℹ️  Block #%d: Teacher field is empty (allowed for compatibility)\n",
					currentBlock.Index)
			}
		}

		fmt.Printf("✅ Block #%d valid\n", currentBlock.Index)
	}

	fmt.Println("=== Validation complete ===\n")
	return true
}

// PrintForkStatus выводит текущий статус форков
func (bc *Blockchain) PrintForkStatus() {
	currentIndex := 0
	if len(bc.Blocks) > 0 {
		currentIndex = bc.Blocks[len(bc.Blocks)-1].Index
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("FORK STATUS")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Current blockchain height: %d blocks\n", len(bc.Blocks))
	fmt.Printf("Next block will be: #%d\n", currentIndex+1)
	fmt.Println()

	// Soft Fork Status
	if currentIndex+1 >= bc.ForkConfig.SoftForkHeight {
		fmt.Println("✅ SOFT FORK ACTIVE (since block #5)")
		fmt.Println("   Changes:")
		fmt.Println("   • Difficulty increased: 00 → 000")
		fmt.Println("   • New field: Teacher")
		fmt.Println("   • Old blocks remain valid (backwards compatible)")
	} else {
		fmt.Printf("⏳ Soft Fork will activate at block #%d\n", bc.ForkConfig.SoftForkHeight)
		fmt.Printf("   Blocks until activation: %d\n", bc.ForkConfig.SoftForkHeight-(currentIndex+1))
	}

	fmt.Println()

	// Hard Fork Status
	if currentIndex+1 >= bc.ForkConfig.HardForkHeight {
		fmt.Println("⚠️  HARD FORK ACTIVE (since block #10)")
		fmt.Println("   Changes:")
		fmt.Println("   • Hash algorithm: SHA3-384 → SHA3-512")
		fmt.Println("   • Storage format: JSON → CSV")
		fmt.Println("   • NOT backwards compatible!")
	} else {
		fmt.Printf("⏳ Hard Fork will activate at block #%d\n", bc.ForkConfig.HardForkHeight)
		fmt.Printf("   Blocks until activation: %d\n", bc.ForkConfig.HardForkHeight-(currentIndex+1))
	}

	fmt.Println(strings.Repeat("=", 60) + "\n")
}

// GetStatistics returns blockchain statistics
func (bc *Blockchain) GetStatistics() map[string]any {
	stats := make(map[string]any)

	stats["total_blocks"] = len(bc.Blocks)

	gradeCount := make(map[int]int)
	courseCount := make(map[int]int)
	forkVersions := make(map[int]int)

	for _, blk := range bc.Blocks {
		if blk.Index == 0 {
			continue
		}
		gradeCount[blk.Data.Grade]++
		courseCount[blk.Data.Course]++
		forkVersions[blk.ForkVersion]++
	}

	stats["grades"] = gradeCount
	stats["courses"] = courseCount
	stats["fork_versions"] = forkVersions

	return stats
}
