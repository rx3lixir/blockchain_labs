package cli

import (
	"flag"
	"fmt"

	"github.com/rx3lixir/lab_bc/internal/blockchain"
	"github.com/rx3lixir/lab_bc/internal/storage"
)

func Run() error {
	// Команды
	listFlag := flag.Bool("list", false, "List all blocks")
	validateFlag := flag.Bool("validate", false, "Validate blockchain")
	searchFlag := flag.String("search", "", "Search keyword")
	addFlag := flag.Bool("add", false, "Add new record")

	// Параметры для добавления записи
	name := flag.String("name", "", "Student name")
	course := flag.Int("course", 0, "Course number")
	group := flag.String("group", "", "Group name")
	zachetka := flag.String("zachetka", "", "Zachetka number")
	subject := flag.String("subject", "", "Subject name")
	grade := flag.Int("grade", 0, "Grade (2-5)")

	flag.Parse()

	store := storage.NewJSONStorage("blockchain.json")
	app, err := NewApp(store)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	switch {
	case *listFlag:
		return app.CmdList()

	case *validateFlag:
		return app.CmdValidate()

	case *searchFlag != "":
		return app.CmdSearch(*searchFlag)

	case *addFlag:
		record := blockchain.StudentRecord{
			FullName: *name,
			Zachetka: *zachetka,
			Group:    *group,
			Subject:  *subject,
			Course:   *course,
			Grade:    *grade,
		}
		return app.CmdAdd(record)

	default:
		flag.Usage()
		return nil
	}
}
