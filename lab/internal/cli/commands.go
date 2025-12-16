package cli

import (
	"fmt"

	"github.com/rx3lixir/lab_bc/internal/blockchain"
)

func (a *App) CmdList() error {
	blocks := a.bc.Blocks()
	fmt.Printf("Total blocks: %d\n\n", len(blocks))

	for _, block := range blocks {
		PrintBlock(block)
	}
	return nil
}

func (a *App) CmdValidate() error {
	if err := a.bc.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	fmt.Println("✓ Blockchain is valid")
	return nil
}

func (a *App) CmdSearch(query string) error {
	results := a.bc.Search(query)
	fmt.Printf("Found %d results\n\n", len(results))

	for _, block := range results {
		PrintBlock(block)
	}
	return nil
}

func (a *App) CmdAdd(record blockchain.StudentRecord) error {
	fmt.Println("Mining block...")

	miningTime, err := a.bc.AddBlock(record)
	if err != nil {
		return fmt.Errorf("failed to add block: %w", err)
	}

	if err := a.storage.Save(a.bc); err != nil {
		return fmt.Errorf("failed to save blockchain: %w", err)
	}

	fmt.Printf("✓ Block mined successfully in %v\n", miningTime)
	return nil
}

func PrintBlock(b *blockchain.Block) {
	fmt.Printf("========== Block #%d ==========\n", b.Index)
	fmt.Printf("Timestamp:    %d\n", b.Timestamp)
	fmt.Printf("ID:           %s\n", b.Data.ID)
	fmt.Printf("Name:         %s\n", b.Data.FullName)
	fmt.Printf("Zachetka:     %s\n", b.Data.Zachetka)
	fmt.Printf("Group:        %s\n", b.Data.Group)
	fmt.Printf("Subject:      %s\n", b.Data.Subject)
	fmt.Printf("Course:       %d\n", b.Data.Course)
	fmt.Printf("Grade:        %d\n", b.Data.Grade)
	fmt.Printf("Hash:         %s...\n", b.Hash)
	fmt.Printf("PreviousHash: %s...\n", b.PreviousHash)
	fmt.Printf("Nonce:        %d\n", b.Nonce)
	fmt.Println()
}
