package domain

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

type Hasher interface {
	Hash(block *Block) string
	Algorithm() string
}

type SHA384Hasher struct{}

func (h *SHA384Hasher) Hash(block *Block) string {
	record := fmt.Sprintf(
		"%d%d%d%s%s%s%s%d%s%s%d",
		block.Index, block.Timestamp, block.Data.Grade,
		block.Data.FullName, block.Data.Zachetka, block.Data.Group,
		block.Data.Subject, block.Data.Course, block.Data.Teacher,
		block.PreviousHash, block.Nonce,
	)
	hash := sha512.Sum384([]byte(record))
	return hex.EncodeToString(hash[:])
}

func (h *SHA384Hasher) Algorithm() string {
	return "SHA3-384"
}

type SHA512Hasher struct{}

func (h *SHA512Hasher) Hash(block *Block) string {
	record := fmt.Sprintf(
		"%d%d%d%s%s%s%s%d%s%s%d",
		block.Index, block.Timestamp, block.Data.Grade,
		block.Data.FullName, block.Data.Zachetka, block.Data.Group,
		block.Data.Subject, block.Data.Course, block.Data.Teacher,
		block.PreviousHash, block.Nonce,
	)
	hash := sha512.Sum512([]byte(record))
	return hex.EncodeToString(hash[:])
}

func (h *SHA512Hasher) Algorithm() string {
	return "SHA3-512"
}
