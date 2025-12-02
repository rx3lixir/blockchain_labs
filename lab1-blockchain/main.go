package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rx3lixir/lab1-blockchain/block"
	"github.com/rx3lixir/lab1-blockchain/blockchain"
	"github.com/rx3lixir/lab1-blockchain/models"
	"github.com/rx3lixir/lab1-blockchain/storage"
)

func main() {
	// Все флаги
	addFlag := flag.Bool("add", false, "Добавить запись")
	listFlag := flag.Bool("list", false, "Показать все блоки")
	validateFlag := flag.Bool("validate", false, "Проверить целостность")
	debtsFlag := flag.Bool("debts", false, "Найти студентов с долгами")
	searchIndex := flag.Int("index", -1, "Поиск по индексу блока")
	searchKeyword := flag.String("keyword", "", "Поиск по ключевому слову")

	// Флаги для добавления записи
	course := flag.Int("course", 0, "Курс")
	group := flag.String("group", "", "Группа")
	name := flag.String("name", "", "ФИО студента")
	recordBook := flag.String("record", "", "Номер зачётки")
	subject := flag.String("subject", "", "Дисциплина")
	grade := flag.Int("grade", 0, "Оценка (2-5)")

	flag.Parse()

	// Загружает или создаём блокчейн
	bc, err := storage.LoadBlockchain()
	if err != nil {
		fmt.Printf("Ошибка загрузки: %v\n", err)
		os.Exit(1)
	}

	if bc == nil {
		fmt.Println("Создаём новый блокчейн...")
		bc = blockchain.New()
		storage.SaveBlockchain(bc)
	}

	// Обрабатываем команды
	if *addFlag {
		if *course == 0 || *group == "" || *name == "" || *recordBook == "" || *subject == "" || *grade == 0 {
			fmt.Println("Для добавления нужны все флаги: -course -group -name -record -subject -grade")
			flag.PrintDefaults()
			return
		}

		if *grade < 2 || *grade > 5 {
			fmt.Println("Оценка должна быть от 2 до 5")
			return
		}

		studentRecord := models.StudentRecord{
			Course:   *course,
			Group:    *group,
			FullName: *name,
			Zachetka: *recordBook,
			Subject:  *subject,
			Grade:    *grade,
		}

		fmt.Println("\nМайним новый блок...")

		bc.AddBlock(studentRecord)
		storage.SaveBlockchain(bc)

		fmt.Println("Запись добавлена!")

		// make list
	} else if *listFlag {
		bc.PrintAllBlocks()

		// make validate
	} else if *validateFlag {
		fmt.Println("\nПроверяем целостность блокчейна...")
		bc.ValidateChain()

		// make debts
	} else if *debtsFlag {
		fmt.Println("\nСТУДЕНТЫ С ДОЛГАМИ:")
		debts := bc.FindStudentsWithDebts()
		if len(debts) == 0 {
			fmt.Println("Студентов с долгами нет!")
			return
		}

		for _, blk := range debts {
			fmt.Printf("• %s (№%s) - %s - оценка: %d\n",
				blk.Data.FullName,
				blk.Data.Zachetka,
				blk.Data.Subject,
				blk.Data.Grade,
			)
		}

		// make search [index]
	} else if *searchIndex >= 0 {
		fmt.Printf("\nПоиск блока по индексу#%d...\n", *searchIndex)
		blk := bc.GetBlockByIndex(*searchIndex)
		block.PrintBlock(blk)

		// make search [keyword]
	} else if *searchKeyword != "" {
		fmt.Printf("\nПоиск по ключевому слову '%s'...\n", *searchKeyword)
		results := bc.SearchByKeyword(*searchKeyword)
		fmt.Printf("Найдено: %d\n", len(results))
		for _, blk := range results {
			block.PrintBlock(blk)
		}
	}
}
