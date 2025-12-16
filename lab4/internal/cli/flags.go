package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rx3lixir/lab_bc/internal/blockchain"
)

func Run() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	listFlag := flag.Bool("list", false, "List all blocks")
	validateFlag := flag.Bool("validate", false, "Validate blockchain")
	addFlag := flag.Bool("add", false, "Add new transaction(s)")

	// Merkle команды
	merkleBuildFlag := flag.Int("merkle-build", -1, "Build Merkle tree for block")
	merkleProofFlag := flag.String("merkle-proof", "", "Get Merkle proof (format: block,tx)")
	merkleVerifyFlag := flag.String("merkle-verify", "", "Verify transaction (format: block,tx)")

	// Данные для транзакций
	names := flag.String("names", "", "Student names (comma-separated)")
	courses := flag.String("courses", "", "Course numbers (comma-separated)")
	groups := flag.String("groups", "", "Group names (comma-separated)")
	zachetkas := flag.String("zachetkas", "", "Zachetka numbers (comma-separated)")
	subjects := flag.String("subjects", "", "Subject names (comma-separated)")
	grades := flag.String("grades", "", "Grades (comma-separated, 2-5)")

	flag.CommandLine.Parse(os.Args[1:])

	app, err := NewApp("blockchain.json")
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	switch {
	case *listFlag:
		return app.CmdList()

	case *validateFlag:
		return app.CmdValidate()

	case *merkleBuildFlag >= 0:
		return app.CmdMerkleBuild(*merkleBuildFlag)

	case *merkleProofFlag != "":
		var blockIdx, txIdx int
		if _, err := fmt.Sscanf(*merkleProofFlag, "%d,%d", &blockIdx, &txIdx); err != nil {
			return fmt.Errorf("invalid format, use: block,tx (e.g., 1,0)")
		}
		return app.CmdMerkleProof(blockIdx, txIdx)

	case *merkleVerifyFlag != "":
		var blockIdx, txIdx int
		if _, err := fmt.Sscanf(*merkleVerifyFlag, "%d,%d", &blockIdx, &txIdx); err != nil {
			return fmt.Errorf("invalid format, use: block,tx (e.g., 1,0)")
		}
		return app.CmdMerkleVerify(blockIdx, txIdx)

	case *addFlag:
		transactions, err := parseTransactions(*names, *zachetkas, *groups, *subjects, *courses, *grades)
		if err != nil {
			return err
		}
		return app.CmdAdd(transactions)

	default:
		printUsage()
		return nil
	}
}

func parseTransactions(names, zachetkas, groups, subjects, courses, grades string) ([]blockchain.StudentRecord, error) {
	nameList := splitAndTrim(names)
	zachetkaList := splitAndTrim(zachetkas)
	groupList := splitAndTrim(groups)
	subjectList := splitAndTrim(subjects)
	courseList := splitAndTrim(courses)
	gradeList := splitAndTrim(grades)

	if len(nameList) == 0 {
		return nil, fmt.Errorf("at least one name is required")
	}

	count := len(nameList)
	transactions := make([]blockchain.StudentRecord, count)

	for i := 0; i < count; i++ {
		transactions[i] = blockchain.StudentRecord{
			FullName: nameList[i],
			Zachetka: getOrDefault(zachetkaList, i, "000000"),
			Group:    getOrDefault(groupList, i, "Unknown"),
			Subject:  getOrDefault(subjectList, i, "Unknown"),
			Course:   parseIntOrDefault(getOrDefault(courseList, i, "0")),
			Grade:    parseIntOrDefault(getOrDefault(gradeList, i, "0")),
		}
	}

	return transactions, nil
}

func splitAndTrim(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func getOrDefault(list []string, index int, defaultVal string) string {
	if index < len(list) {
		return list[index]
	}
	return defaultVal
}

func parseIntOrDefault(s string) int {
	var val int
	fmt.Sscanf(s, "%d", &val)
	return val
}

func printUsage() {
	fmt.Println("Blockchain with Merkle Tree - Lab Work")
	fmt.Println("Usage: bc <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  -list                        List all blocks")
	fmt.Println("  -validate                    Validate blockchain integrity")
	fmt.Println("  -add                         Add new transaction(s) to blockchain")
	fmt.Println()
	fmt.Println("Merkle Tree Commands (main lab focus):")
	fmt.Println("  -merkle-build <block>        Build and display Merkle tree for block")
	fmt.Println("  -merkle-proof <block,tx>     Get Merkle proof for transaction (SPV)")
	fmt.Println("  -merkle-verify <block,tx>    Verify transaction using Merkle proof")
	fmt.Println()
	fmt.Println("Options for -add (comma-separated for multiple transactions):")
	fmt.Println("  -names <string,...>      Student full names")
	fmt.Println("  -courses <int,...>       Course numbers")
	fmt.Println("  -groups <string,...>     Group names")
	fmt.Println("  -zachetkas <string,...>  Zachetka numbers")
	fmt.Println("  -subjects <string,...>   Subject names")
	fmt.Println("  -grades <int,...>        Grades (2-5)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  # Add block with multiple transactions")
	fmt.Println("  bc -add -names \"Иванов И.И.,Петров П.П.,Сидоров С.С.\" \\")
	fmt.Println("     -grades \"5,4,3\" -courses \"5,5,5\" \\")
	fmt.Println("     -groups \"5.507M,5.507M,5.507M\" \\")
	fmt.Println("     -zachetkas \"202434,202435,202436\" \\")
	fmt.Println("     -subjects \"Математика,Физика,Химия\"")
	fmt.Println()
	fmt.Println("  # Build Merkle tree visualization")
	fmt.Println("  bc -merkle-build 1")
	fmt.Println()
	fmt.Println("  # Get SPV proof for transaction")
	fmt.Println("  bc -merkle-proof 1,0")
	fmt.Println()
	fmt.Println("  # Verify transaction with SPV")
	fmt.Println("  bc -merkle-verify 1,0")
}
