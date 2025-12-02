package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/rx3lixir/lab1-blockchain/blockchain"
	"github.com/rx3lixir/lab1-blockchain/models"
	"github.com/rx3lixir/lab1-blockchain/storage"
	"github.com/rx3lixir/lab1-blockchain/ui"
)

func main() {
	// Define command flags
	listFlag := flag.Bool("list", false, "Show all blocks")
	validateFlag := flag.Bool("validate", false, "Validate blockchain integrity")
	forkStatusFlag := flag.Bool("fork-status", false, "Show fork status")

	// Search flags
	searchFlag := flag.String("search", "", "Search blocks (by index, name, group, subject, teacher or grade)")
	filterGrade := flag.String("filter-grade", "", "Filter by grade (2-5)")

	// Add flags
	addFlag := flag.Bool("add", false, "Add new record")
	course := flag.Int("course", 0, "Course number")
	group := flag.String("group", "", "Group name")
	name := flag.String("name", "", "Student full name")
	recordBook := flag.String("record", "", "Record book number")
	subject := flag.String("subject", "", "Subject name")
	teacher := flag.String("teacher", "", "Teacher name (required for blocks after Soft Fork)")
	grade := flag.Int("grade", 0, "Grade (2-5)")

	// Fork commands
	exportCSV := flag.Bool("export-csv", false, "Export blockchain to CSV (Hard Fork format)")
	loadCSV := flag.Bool("load-csv", false, "Load blockchain from CSV")
	createHardFork := flag.Bool("create-hardfork", false, "Create Hard Fork branch at block 10")

	flag.Parse()

	// Определяем какой формат использовать
	var bc *blockchain.Blockchain
	var err error

	// Если указана работа с CSV или есть CSV файл
	if *loadCSV || storage.CSVExists() {
		bc, err = storage.LoadBlockchainCSV()
		if err != nil {
			fmt.Printf("❌ Failed to load CSV blockchain: %v\n", err)
			os.Exit(1)
		}
		if bc != nil {
			fmt.Println("📁 Using CSV format (Hard Fork branch)")
		}
	}

	// Если не загрузили из CSV, пробуем JSON
	if bc == nil {
		bc, err = storage.LoadBlockchain()
		if err != nil {
			fmt.Printf("❌ Failed to load blockchain: %v\n", err)
			os.Exit(1)
		}
	}

	// Если ничего нет - создаём новый
	if bc == nil {
		fmt.Println("🔨 Creating new blockchain...")
		bc = blockchain.New()
		storage.SaveBlockchain(bc)
	}

	switch {
	case *createHardFork:
		handleCreateHardFork(bc)
	case *exportCSV:
		handleExportCSV(bc)
	case *addFlag:
		handleAdd(bc, *course, *group, *name, *subject, *recordBook, *teacher, *grade)
	case *listFlag:
		ui.ShowAllBlocks(bc.Blocks)
	case *validateFlag:
		handleValidate(bc)
	case *forkStatusFlag:
		bc.PrintForkStatus()
	case *searchFlag != "":
		handleSearch(bc, *searchFlag, *filterGrade)
	default:
		showUsage()
	}
}

func handleAdd(
	bc *blockchain.Blockchain,
	course int,
	group, name, subject, recordBook, teacher string,
	grade int,
) {
	if course == 0 || group == "" || name == "" || recordBook == "" || subject == "" || grade == 0 {
		ui.ShowError("All fields required: -course -group -name -record -subject -grade")
		flag.PrintDefaults()
		return
	}

	if grade < 2 || grade > 5 {
		ui.ShowError("Grade must be between 2 and 5")
		return
	}

	// Проверяем нужно ли поле Teacher
	nextIndex := len(bc.Blocks)
	if nextIndex >= bc.ForkConfig.SoftForkHeight && teacher == "" {
		ui.ShowWarning("Warning: Teacher field is empty after Soft Fork activation")
		ui.ShowInfo("Consider providing -teacher flag for complete data")
	}

	studentRecord := models.StudentRecord{
		Course:   course,
		Group:    group,
		FullName: name,
		Zachetka: recordBook,
		Subject:  subject,
		Grade:    grade,
		Teacher:  teacher,
	}

	ui.ShowInfo("Mining new block...")
	bc.AddBlock(studentRecord)

	// Сохраняем в правильном формате
	if nextIndex+1 >= bc.ForkConfig.HardForkHeight || storage.CSVExists() {
		// Hard Fork активен - сохраняем в CSV
		if err := storage.SaveBlockchainCSV(bc); err != nil {
			ui.ShowError(fmt.Sprintf("Failed to save CSV: %v", err))
			return
		}
	} else {
		// До Hard Fork - сохраняем в JSON
		if err := storage.SaveBlockchain(bc); err != nil {
			ui.ShowError(fmt.Sprintf("Failed to save: %v", err))
			return
		}
	}

	ui.ShowSuccess("Record added successfully")
}

func handleValidate(bc *blockchain.Blockchain) {
	ui.ShowInfo("Validating blockchain integrity with fork rules...")
	if bc.ValidateChain() {
		ui.ShowSuccess("✅ Blockchain is valid!")
	} else {
		ui.ShowError("❌ Blockchain validation failed!")
	}
}

func handleSearch(bc *blockchain.Blockchain, query, gradeFilter string) {
	var results []*models.Block

	// Try to parse as index first
	if index, err := strconv.Atoi(query); err == nil {
		if block := bc.GetBlockByIndex(index); block != nil {
			results = []*models.Block{block}
		}
	} else {
		// Search by keyword (includes Teacher field)
		results = bc.SearchByKeyword(query)
	}

	// Apply grade filter if specified
	if gradeFilter != "" {
		if gradeNum, err := strconv.Atoi(gradeFilter); err == nil {
			filtered := []*models.Block{}
			for _, block := range results {
				if block.Data.Grade == gradeNum {
					filtered = append(filtered, block)
				}
			}
			results = filtered
		}
	}

	if len(results) == 0 {
		ui.ShowWarning(fmt.Sprintf("No blocks found for query: %s", query))
		return
	}

	ui.ShowSearchResults(results, query)
}

func handleExportCSV(bc *blockchain.Blockchain) {
	ui.ShowInfo("Exporting blockchain to CSV format (Hard Fork)...")
	if err := storage.ExportToCSV(bc); err != nil {
		ui.ShowError(fmt.Sprintf("Export failed: %v", err))
		return
	}
	ui.ShowSuccess("Blockchain exported to CSV successfully!")
	ui.ShowInfo("You can now use -load-csv to work with CSV format")
}

func handleCreateHardFork(bc *blockchain.Blockchain) {
	forkHeight := bc.ForkConfig.HardForkHeight

	if len(bc.Blocks) < forkHeight {
		ui.ShowError(fmt.Sprintf("Cannot create Hard Fork: need at least %d blocks, have %d", forkHeight, len(bc.Blocks)))
		ui.ShowInfo(fmt.Sprintf("Add %d more blocks first", forkHeight-len(bc.Blocks)))
		return
	}

	ui.ShowInfo(fmt.Sprintf("Creating Hard Fork branch at block #%d...", forkHeight))

	// Создаём ветку Hard Fork
	hardForkBC, err := storage.CreateHardForkBranch(bc, forkHeight)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to create Hard Fork: %v", err))
		return
	}

	// Сохраняем в CSV
	if err := storage.SaveBlockchainCSV(hardForkBC); err != nil {
		ui.ShowError(fmt.Sprintf("Failed to save Hard Fork branch: %v", err))
		return
	}

	ui.ShowSuccess("Hard Fork branch created and saved to CSV!")
	fmt.Println()
	fmt.Println("📊 Fork Summary:")
	fmt.Println("   • Original chain: blockchain.json (JSON, SHA3-384)")
	fmt.Println("   • Hard Fork chain: blockchain_hardfork.csv (CSV, SHA3-512)")
	fmt.Println()
	fmt.Println("💡 Usage:")
	fmt.Println("   • Continue original: ./bin/bc-km -add ...")
	fmt.Println("   • Continue Hard Fork: ./bin/bc-km -load-csv -add ...")
	fmt.Println()
}

func showUsage() {
	usage := `
BLOCKCHAIN WITH FORKS

USAGE:
  blockchain-app [command] [options]

📋 BASIC COMMANDS:
  -list                          Show all blocks in the blockchain
  -validate                      Validate blockchain integrity with fork rules
  -fork-status                   Show current fork activation status
  -search <query>                Search blocks by index, name, group, subject, or teacher
  -search <query> -filter-grade <2-5>  Filter search results by grade
  
➕ ADD RECORD:
  -add [options]                 Add new student record
    Required options:
      -course <number>             Course number
      -group <name>                Group name
      -name <full name>            Student full name
      -record <number>             Record book number
      -subject <name>              Subject name
      -grade <2-5>                 Grade (2-5)
      -teacher <name>              Teacher name (REQUIRED after Soft Fork #5)

🔄 FORK COMMANDS:
  -export-csv                    Export blockchain to CSV format (Hard Fork)
  -load-csv                      Load and work with CSV blockchain
  -create-hardfork               Create Hard Fork branch at block #10

📚 FORK INFORMATION:
  
  🟢 SOFT FORK (Block #5):
     • Difficulty: 00 → 000 (harder mining)
     • New field: Teacher (optional for compatibility)
     • Old blocks remain valid (backwards compatible)
  
  🔴 HARD FORK (Block #10):
     • Hash algorithm: SHA3-384 → SHA3-512
     • Storage format: JSON → CSV
     • Creates separate incompatible branch
     • NOT backwards compatible!

📖 EXAMPLES:

  # View fork status
  ./bc-km -fork-status

  # Add record before Soft Fork (blocks 0-4)
  ./bc-km -add -course 5 -group "5.507M" -name "Петров П.П." \\
          -record "202435" -subject "Физика" -grade 4

  # Add record after Soft Fork (blocks 5+) - Teacher required!
  ./bc-km -add -course 5 -group "5.507M" -name "Сидоров С.С." \\
          -record "202436" -subject "Химия" -grade 5 \\
          -teacher "Иванов И.И."

  # Create Hard Fork and export to CSV
  ./bc-km -create-hardfork

  # Work with Hard Fork branch (CSV)
  ./bc-km -load-csv -list
  ./bc-km -load-csv -add [options...]

  # Search by teacher name
  ./bc-km -search "Иванов"

  # Validate with fork rules
  ./bc-km -validate
`
	fmt.Println(usage)
}
