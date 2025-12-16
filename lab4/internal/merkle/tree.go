package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/rx3lixir/lab_bc/internal/blockchain"
)

type Node struct {
	Hash  string
	Left  *Node
	Right *Node
	Data  *blockchain.StudentRecord
}

type MerkleTree struct {
	Root   *Node
	Leaves []*Node
}

func BuildTree(transactions []blockchain.StudentRecord) *MerkleTree {
	if len(transactions) == 0 {
		return nil
	}

	leaves := make([]*Node, 0, len(transactions))
	for i := range transactions {
		hash := HashTransaction(&transactions[i])
		node := &Node{
			Hash: hash,
			Data: &transactions[i],
		}
		leaves = append(leaves, node)
	}

	if len(leaves)%2 != 0 {
		lastLeaf := leaves[len(leaves)-1]
		duplicate := &Node{
			Hash: lastLeaf.Hash,
			Data: lastLeaf.Data,
		}
		leaves = append(leaves, duplicate)
	}

	currentLevel := leaves
	for len(currentLevel) > 1 {
		nextLevel := make([]*Node, 0, len(currentLevel)/2)

		for i := 0; i < len(currentLevel); i += 2 {
			left := currentLevel[i]
			right := currentLevel[i+1]

			combinedHash := HashCombined(left.Hash, right.Hash)
			parent := &Node{
				Hash:  combinedHash,
				Left:  left,
				Right: right,
			}
			nextLevel = append(nextLevel, parent)
		}

		currentLevel = nextLevel
	}

	return &MerkleTree{
		Root:   currentLevel[0],
		Leaves: leaves,
	}
}

func HashTransaction(record *blockchain.StudentRecord) string {
	data := fmt.Sprintf(
		"%s%s%s%s%s%d%d",
		record.ID,
		record.FullName,
		record.Zachetka,
		record.Group,
		record.Subject,
		record.Course,
		record.Grade,
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func HashCombined(left, right string) string {
	combined := left + right
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

func (mt *MerkleTree) GetRoot() string {
	if mt.Root == nil {
		return ""
	}
	return mt.Root.Hash
}

func (mt *MerkleTree) GetProof(txIndex int) (*MerkleProof, error) {
	if txIndex < 0 || txIndex >= len(mt.Leaves) {
		return nil, fmt.Errorf("invalid transaction index: %d", txIndex)
	}

	proof := &MerkleProof{
		Index:  txIndex,
		Hashes: []string{},
	}

	currentIndex := txIndex

	allLevels := [][]string{}

	currentLevelHashes := make([]string, len(mt.Leaves))
	for i, leaf := range mt.Leaves {
		currentLevelHashes[i] = leaf.Hash
	}
	allLevels = append(allLevels, currentLevelHashes)

	for len(currentLevelHashes) > 1 {
		nextLevel := make([]string, 0, len(currentLevelHashes)/2)
		for i := 0; i < len(currentLevelHashes); i += 2 {
			combined := HashCombined(currentLevelHashes[i], currentLevelHashes[i+1])
			nextLevel = append(nextLevel, combined)
		}
		allLevels = append(allLevels, nextLevel)
		currentLevelHashes = nextLevel
	}

	currentIndex = txIndex

	for levelIdx := 0; levelIdx < len(allLevels)-1; levelIdx++ {
		level := allLevels[levelIdx]
		isLeft := currentIndex%2 == 0

		var siblingIndex int
		if isLeft {
			siblingIndex = currentIndex + 1
		} else {
			siblingIndex = currentIndex - 1
		}

		if siblingIndex < len(level) {
			proof.Hashes = append(proof.Hashes, level[siblingIndex])
		}

		currentIndex = currentIndex / 2
	}

	return proof, nil
}

func VerifyProof(txHash string, proof *MerkleProof, expectedRoot string) bool {
	currentHash := txHash
	currentIndex := proof.Index

	for _, siblingHash := range proof.Hashes {
		if currentIndex%2 == 0 {
			// слева, сосед справа
			currentHash = HashCombined(currentHash, siblingHash)
		} else {
			// справа, сосед слева
			currentHash = HashCombined(siblingHash, currentHash)
		}
		currentIndex = currentIndex / 2
	}

	return currentHash == expectedRoot
}

type MerkleProof struct {
	Index  int
	Hashes []string
}

func (mt *MerkleTree) PrintTree() {
	if mt.Root == nil {
		fmt.Println("Empty tree")
		return
	}
	fmt.Println("=== Merkle Tree Structure ===")
	mt.printNode(mt.Root, "", true)
	fmt.Printf("\nRoot Hash: %s\n", mt.Root.Hash[:16]+"...")
}

func (mt *MerkleTree) printNode(node *Node, prefix string, isTail bool) {
	if node == nil {
		return
	}

	connector := "└── "
	if !isTail {
		connector = "├── "
	}

	hashPreview := node.Hash[:16] + "..."
	if node.Data != nil {
		fmt.Printf("%s%s%s (Tx: %s)\n", prefix, connector, hashPreview, node.Data.FullName)
	} else {
		fmt.Printf("%s%s%s\n", prefix, connector, hashPreview)
	}

	if node.Left != nil || node.Right != nil {
		extension := "    "
		if !isTail {
			extension = "│   "
		}

		if node.Right != nil {
			mt.printNode(node.Right, prefix+extension, false)
		}
		if node.Left != nil {
			mt.printNode(node.Left, prefix+extension, true)
		}
	}
}
