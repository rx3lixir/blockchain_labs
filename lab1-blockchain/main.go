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

	// Search flags
	searchFlag := flag.String("search", "", "Search blocks (by index, name, group, subject or grade)")
	filterGrade := flag.String("filter-grade", "", "Filter by grade (2-5)")

	// Add flags
	addFlag := flag.Bool("add", false, "Add new record")
	course := flag.Int("course", 0, "Course number")
	group := flag.String("group", "", "Student full name")
	name := flag.String("name", "", "Student full name")
	recordBook := flag.String("record", "", "Record book number")
	subject := flag.String("subject", "", "Subject name")
	grade := flag.Int("grade", 0, "Grade (2-5)")

	flag.Parse()

	// Load or create blockchain
	bc, err := storage.LoadBlockchain()
	if err != nil {
		fmt.Printf("Failed to load blockchain: %v\n", err)
		os.Exit(1)
	}

	if bc == nil {
		fmt.Println("Creating new blockchain...")
		bc = blockchain.New()
		storage.SaveBlockchain(bc)
	}

	switch {
	case *addFlag:
		handleAdd(bc, *course, *group, *name, *subject, *recordBook, *grade)
	case *listFlag:
		ui.ShowAllBlocks(bc.Blocks)
	case *validateFlag:
		handleValidate(bc)
	case *searchFlag != "":
		handleSearch(bc, *searchFlag, *filterGrade)
	default:
		showUsage()
	}
}

func handleAdd(
	bc *blockchain.Blockchain,
	course int,
	group, name, subject, recordBook string,
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

	studentRecord := models.StudentRecord{
		Course:   course,
		Group:    group,
		FullName: name,
		Zachetka: recordBook,
		Subject:  subject,
		Grade:    grade,
	}

	ui.ShowInfo("Mining new block...")

	bc.AddBlock(studentRecord)
	storage.SaveBlockchain(bc)

	ui.ShowSuccess("Record added successfully")
}

func handleValidate(bc *blockchain.Blockchain) {
	ui.ShowInfo("Validating blockchain integrity...")
	if bc.ValidateChain() {
		ui.ShowSuccess("Blockchain is valid!")
	} else {
		ui.ShowError("Blockchain validation failed!")
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
		// Search by keyword
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

func showUsage() {
	usage := `

USAGE:
  blockchain-app [command] [options]

COMMANDS:
  -list                          Show all blocks in the blockchain
  -validate                      Validate blockchain integrity
  -search <query>                Search blocks by index, name, group, or subject
  -search <query> -filter-grade <2-5>  Filter search results by grade
  
  -add [options]                 Add new student record
    Required options:
      -course <number>             Course number
      -group <name>                Group name
      -name <full name>            Student full name
      -record <number>             Record book number
      -subject <name>              Subject name
      -grade <2-5>                 Grade (2-5)

EXAMPLES:
  ./blockchain-app -list
  ./blockchain-app -validate
  ./blockchain-app -search "Иванов"
  ./blockchain-app -search "5.507M" -filter-grade 2
  ./blockchain-app -search 1
  ./blockchain-app -add -course 5 -group "5.507M" -name "Иванов И.И." \
                   -record "202434" -subject "Математика" -grade 5
`
	fmt.Println(usage)
}
