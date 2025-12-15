package storage

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/rx3lixir/lab_bc/internal/domain"
)

type CSVStorage struct {
	filename string
}

func NewCSVStorage(filename string) *CSVStorage {
	return &CSVStorage{filename: filename}
}

func (s *CSVStorage) Save(bc *domain.Blockchain) error {
	file, err := os.Create(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Добавляем строку с ForkConfig в начало файла (можно как комментарий или отдельно, но проще в хедер)
	header := []string{
		"Index", "Timestamp", "FullName", "Zachetka", "Group",
		"Subject", "Course", "Grade", "Teacher", "PreviousHash", "Hash", "Nonce", "ForkVersion",
		"SoftForkHeight", "HardForkHeight", "DifficultyOld", "DifficultyNew", // добавляем config
	}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, block := range bc.Blocks() {
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
			strconv.Itoa(bc.ForkConfig().SoftForkHeight),
			strconv.Itoa(bc.ForkConfig().HardForkHeight),
			bc.ForkConfig().DifficultyOld,
			bc.ForkConfig().DifficultyNew,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

func (s *CSVStorage) Load() (*domain.Blockchain, error) {
	if !s.Exists() {
		return nil, nil
	}

	file, err := os.Open(s.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	// Читаем config из первой строки (последние 4 поля)
	config := &domain.ForkConfig{
		SoftForkHeight: 5,
		HardForkHeight: 10,
		DifficultyOld:  "00",
		DifficultyNew:  "000",
	}

	if len(records) > 1 {
		lastRow := records[len(records)-1] // последняя строка содержит config (или первую после хедера)
		if len(lastRow) >= 17 {
			if v, err := strconv.Atoi(lastRow[13]); err == nil {
				config.SoftForkHeight = v
			}
			if v, err := strconv.Atoi(lastRow[14]); err == nil {
				config.HardForkHeight = v
			}
			config.DifficultyOld = lastRow[15]
			config.DifficultyNew = lastRow[16]
		}
	}

	var blocks []*domain.Block

	for i, r := range records {
		if i == 0 { // пропуск хедера
			continue
		}
		index, _ := strconv.Atoi(r[0])
		timestamp, _ := strconv.ParseInt(r[1], 10, 64)
		course, _ := strconv.Atoi(r[6])
		grade, _ := strconv.Atoi(r[7])
		nonce, _ := strconv.Atoi(r[11])
		forkVersion, _ := strconv.Atoi(r[12])

		blocks = append(blocks, &domain.Block{
			Index:     index,
			Timestamp: timestamp,
			Data: domain.StudentRecord{
				FullName: r[2],
				Zachetka: r[3],
				Group:    r[4],
				Subject:  r[5],
				Course:   course,
				Grade:    grade,
				Teacher:  r[8],
			},
			PreviousHash: r[9],
			Hash:         r[10],
			Nonce:        nonce,
			ForkVersion:  forkVersion,
		})
	}

	return domain.LoadBlockchain(blocks, config), nil
}

func (s *CSVStorage) Exists() bool {
	_, err := os.Stat(s.filename)
	return err == nil
}
