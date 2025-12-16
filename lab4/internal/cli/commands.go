package cli

import (
	"fmt"

	"github.com/rx3lixir/lab_bc/internal/blockchain"
	"github.com/rx3lixir/lab_bc/internal/merkle"
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
	fmt.Println("✓ All Merkle roots are correct")
	return nil
}

func (a *App) CmdAdd(transactions []blockchain.StudentRecord) error {
	fmt.Printf("Mining block with %d transaction(s)...\n", len(transactions))

	miningTime, err := a.bc.AddBlock(transactions)
	if err != nil {
		return fmt.Errorf("failed to add block: %w", err)
	}

	if err := a.storage.Save(a.bc); err != nil {
		return fmt.Errorf("failed to save blockchain: %w", err)
	}

	fmt.Printf("✓ Block mined successfully in %v\n", miningTime)
	return nil
}

func (a *App) CmdMerkleBuild(blockIndex int) error {
	block, err := a.bc.GetBlock(blockIndex)
	if err != nil {
		return err
	}

	if len(block.Transactions) == 0 {
		return fmt.Errorf("block #%d has no transactions", blockIndex)
	}

	tree := merkle.BuildTree(block.Transactions)
	if tree == nil {
		return fmt.Errorf("failed to build merkle tree")
	}

	fmt.Printf("=== Merkle Tree for Block #%d ===\n", blockIndex)
	fmt.Printf("Transactions: %d\n", len(block.Transactions))
	fmt.Printf("Merkle Root: %s\n\n", tree.GetRoot())

	tree.PrintTree()

	return nil
}

func (a *App) CmdMerkleProof(blockIndex, txIndex int) error {
	block, err := a.bc.GetBlock(blockIndex)
	if err != nil {
		return err
	}

	if txIndex < 0 || txIndex >= len(block.Transactions) {
		return fmt.Errorf("transaction index %d out of range (0-%d)", txIndex, len(block.Transactions)-1)
	}

	tree := merkle.BuildTree(block.Transactions)
	if tree == nil {
		return fmt.Errorf("failed to build merkle tree")
	}

	proof, err := tree.GetProof(txIndex)
	if err != nil {
		return fmt.Errorf("failed to get proof: %w", err)
	}

	tx := block.Transactions[txIndex]
	txHash := merkle.HashTransaction(&tx)

	fmt.Printf("=== Merkle Proof (SPV) ===\n")
	fmt.Printf("Block: #%d\n", blockIndex)
	fmt.Printf("Transaction #%d: %s (%s)\n", txIndex, tx.FullName, tx.Subject)
	fmt.Printf("TX Hash: %s...\n", txHash[:32])
	fmt.Printf("Merkle Root: %s...\n\n", block.MerkleRoot[:32])
	fmt.Printf("Proof Path (%d hashes):\n", len(proof.Hashes))
	for i, h := range proof.Hashes {
		fmt.Printf("  [%d] %s...\n", i, h[:32])
	}

	return nil
}

func (a *App) CmdMerkleVerify(blockIndex, txIndex int) error {
	block, err := a.bc.GetBlock(blockIndex)
	if err != nil {
		return err
	}

	if txIndex < 0 || txIndex >= len(block.Transactions) {
		return fmt.Errorf("transaction index %d out of range (0-%d)", txIndex, len(block.Transactions)-1)
	}

	tree := merkle.BuildTree(block.Transactions)
	if tree == nil {
		return fmt.Errorf("failed to build merkle tree")
	}

	proof, err := tree.GetProof(txIndex)
	if err != nil {
		return fmt.Errorf("failed to get proof: %w", err)
	}

	tx := block.Transactions[txIndex]
	txHash := merkle.HashTransaction(&tx)

	isValid := merkle.VerifyProof(txHash, proof, block.MerkleRoot)

	fmt.Printf("=== SPV Verification ===\n")
	fmt.Printf("Block: #%d\n", blockIndex)
	fmt.Printf("Transaction #%d: %s (%s)\n", txIndex, tx.FullName, tx.Subject)
	fmt.Printf("TX Hash: %s...\n", txHash[:32])
	fmt.Printf("Merkle Root: %s...\n", block.MerkleRoot[:32])
	fmt.Printf("Proof size: %d hashes\n\n", len(proof.Hashes))

	if isValid {
		fmt.Println("✓ Transaction VERIFIED successfully!")
		fmt.Println("✓ Transaction is included in the block")
		fmt.Println("✓ SPV verification passed")
	} else {
		fmt.Println("✗ Verification FAILED!")
		fmt.Println("✗ Transaction is NOT in the block or proof is invalid")
	}

	return nil
}

func PrintBlock(b *blockchain.Block) {
	fmt.Printf("========== Block #%d ==========\n", b.Index)
	fmt.Printf("Timestamp:    %d\n", b.Timestamp)
	fmt.Printf("Transactions: %d\n", len(b.Transactions))

	for i, tx := range b.Transactions {
		fmt.Printf("  [%d] %s - %s (Grade: %d)\n", i, tx.FullName, tx.Subject, tx.Grade)
	}

	fmt.Printf("MerkleRoot:   %s...\n", b.MerkleRoot[:16])
	fmt.Printf("Hash:         %s...\n", b.Hash[:16])
	fmt.Printf("PreviousHash: %s...\n", b.PreviousHash[:16])
	fmt.Printf("Nonce:        %d\n", b.Nonce)
	fmt.Println()
}
