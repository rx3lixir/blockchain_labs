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
		"%d%d%s%d%s%s%s%s%d%s%d",
		block.Index,
		block.Timestamp,
		block.Data.ID,
		block.Data.Grade,
		block.Data.FullName,
		block.Data.Zachetka,
		block.Data.Group,
		block.Data.Subject,
		block.Data.Course,
		block.PreviousHash,
		block.Nonce,
	)

	hash := sha256.Sum256([]byte(record))
	return hex.EncodeToString(hash[:])
}
