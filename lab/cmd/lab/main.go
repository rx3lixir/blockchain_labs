package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rx3lixir/lab_bc/internal/cli"
	"github.com/rx3lixir/lab_bc/internal/domain"
	"github.com/rx3lixir/lab_bc/internal/storage"
)

func main() {
	listFlag := flag.Bool("list", false, "List all blocks")
	validateFlag := flag.Bool("validate", false, "Validate blockchain")
	searchFlag := flag.String("search", "", "Search keyword")
	addFlag := flag.Bool("add", false, "Add new record")

	// For add flag
	name := flag.String("name", "", "Student name")
	course := flag.Int("course", 0, "Course number")
	group := flag.String("group", "", "Group name")
	zachetka := flag.String("zachetka", "", "Zachetka number")
	subject := flag.String("subject", "", "Subject name")
	teacher := flag.String("teacher", "", "Teacher name (required for blocks after Soft Fork)")
	grade := flag.Int("grade", 0, "Grade (2-5)")

	useCSV := flag.Bool("csv", false, "Use CSV storage")

	flag.Parse()

	var s storage.Storage
	if *useCSV {
		s = storage.NewCSVStorage("blockchain.csv")
	} else {
		s = storage.NewJSONStorage("blockchain.json")
	}

	app, err := cli.NewApp(s)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	switch {
	case *listFlag:
		app.ListBlocks()

	case *validateFlag:
		if err := app.Validate(); err != nil {
			fmt.Printf("Validation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Blockchain is valid")

	case *searchFlag != "":
		app.Search(*searchFlag)

	case *addFlag:
		record := domain.StudentRecord{
			FullName: *name,
			Zachetka: *zachetka,
			Group:    *group,
			Subject:  *subject,
			Course:   *course,
			Grade:    *grade,
			Teacher:  *teacher,
		}
		if err := app.AddRecord(record); err != nil {
			fmt.Printf("Error adding record: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Record added successfully")
	default:
		flag.Usage()
	}
}
