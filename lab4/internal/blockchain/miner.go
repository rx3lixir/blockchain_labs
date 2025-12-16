package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type Miner struct{}

func NewMiner() *Miner {
	return &Miner{}
}

func (m *Miner) Mine(block *Block, difficulty string) time.Duration {
	start := time.Now()
	block.Nonce = 0

	block.MerkleRoot = CalculateMerkleRoot(block.Transactions)

	for {
		block.Hash = CalculateHash(block)
		if strings.HasPrefix(block.Hash, difficulty) {
			return time.Since(start)
		}
		block.Nonce++
	}
}

func CalculateHash(block *Block) string {
	record := fmt.Sprintf(
		"%d%d%s%s%d",
		block.Index,
		block.Timestamp,
		block.MerkleRoot,
		block.PreviousHash,
		block.Nonce,
	)

	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}

func CalculateMerkleRoot(transactions []StudentRecord) string {
	if len(transactions) == 0 {
		return ""
	}

	hashes := make([]string, len(transactions))
	for i, tx := range transactions {
		hashes[i] = hashTransaction(tx)
	}

	if len(hashes)%2 != 0 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	for len(hashes) > 1 {
		nextLevel := make([]string, 0, len(hashes)/2)

		for i := 0; i < len(hashes); i += 2 {
			combined := hashes[i] + hashes[i+1]
			hash := sha256.Sum256([]byte(combined))
			nextLevel = append(nextLevel, hex.EncodeToString(hash[:]))
		}

		hashes = nextLevel
	}

	return hashes[0]
}

func hashTransaction(tx StudentRecord) string {
	data := fmt.Sprintf(
		"%s%s%s%s%s%d%d",
		tx.ID,
		tx.FullName,
		tx.Zachetka,
		tx.Group,
		tx.Subject,
		tx.Course,
		tx.Grade,
	)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
