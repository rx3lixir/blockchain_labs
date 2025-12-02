package domain

import (
	"strings"
	"time"
)

type Miner struct {
	hasher Hasher
}

func NewMiner(hasher Hasher) *Miner {
	return &Miner{hasher: hasher}
}

func (m *Miner) Mine(block *Block, difficulty string) time.Duration {
	start := time.Now()
	block.Nonce = 0

	for {
		block.Hash = m.hasher.Hash(block)
		if strings.HasPrefix(block.Hash, difficulty) {
			return time.Since(start)
		}
		block.Nonce++
	}
}
