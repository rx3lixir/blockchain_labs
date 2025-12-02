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

	header := []string{
		"Index", "Timestamp", "FullName", "Zachetka", "Group",
		"Subject", "Course", "Grade", "Teacher", "PreviousHash", "Hash", "Nonce", "ForkVersion",
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

	if len(records) < 2 {
		return nil, nil
	}

	return nil, nil
}

func (s *CSVStorage) Exists() bool {
	_, err := os.Stat(s.filename)
	return err == nil
}
