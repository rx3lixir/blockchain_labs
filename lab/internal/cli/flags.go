package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/rx3lixir/lab_bc/internal/blockchain"
	"github.com/rx3lixir/lab_bc/internal/fork"
	"github.com/rx3lixir/lab_bc/internal/resolve"
	"github.com/rx3lixir/lab_bc/internal/storage"
)

const ConfigFile = "fork_config.json"

func Run() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	chainName := os.Args[1]

	listFlag := flag.Bool("list", false, "List all blocks")
	validateFlag := flag.Bool("validate", false, "Validate blockchain(s)")
	searchFlag := flag.String("search", "", "Search keyword")
	addFlag := flag.Bool("add", false, "Add new record")
	forkFlag := flag.String("fork", "", "Create fork from current chain")
	resolveFlag := flag.String("resolve", "", "Resolve fork conflict with another chain")

	name := flag.String("name", "", "Student name")
	course := flag.Int("course", 0, "Course number")
	group := flag.String("group", "", "Group name")
	zachetka := flag.String("zachetka", "", "Zachetka number")
	subject := flag.String("subject", "", "Subject name")
	grade := flag.Int("grade", 0, "Grade (2-5)")

	flag.CommandLine.Parse(os.Args[2:])

	forkMgr, err := fork.NewManager(ConfigFile)
	if err != nil {
		return fmt.Errorf("failed to initialize fork manager: %w", err)
	}

	if err := forkMgr.RegisterChain(chainName); err != nil {
		return fmt.Errorf("failed to register chain: %w", err)
	}

	chainFile, err := forkMgr.GetChainFile(chainName)
	if err != nil {
		return err
	}

	store := storage.NewJSONStorage(chainFile)
	app, err := NewApp(store)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	switch {
	case *listFlag:
		return app.CmdList()

	case *validateFlag != false:
		if flag.NArg() > 0 {
			otherChain := flag.Arg(0)
			resolveMgr := resolve.NewManager(forkMgr)
			return resolveMgr.Validate(chainName, otherChain)
		}
		return app.CmdValidate()

	case *searchFlag != "":
		return app.CmdSearch(*searchFlag)

	case *forkFlag != "":
		return forkMgr.CreateFork(chainName, *forkFlag)

	case *resolveFlag != "":
		resolveMgr := resolve.NewManager(forkMgr)
		return resolveMgr.Resolve(chainName, *resolveFlag)

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
		printUsage()
		return nil
	}
}

func printUsage() {
	fmt.Println("Usage: bc <chain_name> <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  -list                    List all blocks in chain")
	fmt.Println("  -validate [other_chain]  Validate chain(s)")
	fmt.Println("  -search <keyword>        Search for keyword")
	fmt.Println("  -add                     Add new record")
	fmt.Println("  -fork <target_name>      Create fork from current chain")
	fmt.Println("  -resolve <other_chain>   Resolve fork conflict")
	fmt.Println()
	fmt.Println("Options for -add:")
	fmt.Println("  -name <string>      Student full name")
	fmt.Println("  -course <int>       Course number")
	fmt.Println("  -group <string>     Group name")
	fmt.Println("  -zachetka <string>  Zachetka number")
	fmt.Println("  -subject <string>   Subject name")
	fmt.Println("  -grade <int>        Grade (2-5)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  bc main -add -name \"Иванов И.И.\" -grade 5 -course 5 -group \"5.507M\" -zachetka \"202434\" -subject \"Математика\"")
	fmt.Println("  bc main -fork branch_a")
	fmt.Println("  bc branch_a -add -name \"Петров П.П.\" -grade 4 -course 5 -group \"5.507M\" -zachetka \"202435\" -subject \"Физика\"")
	fmt.Println("  bc main -validate branch_a")
	fmt.Println("  bc main -resolve branch_a")
}
