package domain

const (
	ForkOriginal = iota
	ForkSoft
	ForkHard
)

type StudentRecord struct {
	FullName string
	Zachetka string
	Group    string
	Subject  string
	Course   int
	Grade    int
	Teacher  string
}

type Block struct {
	Index        int
	Timestamp    int64
	Data         StudentRecord
	PreviousHash string
	Hash         string
	Nonce        int
	ForkVersion  int
}

type ForkConfig struct {
	SoftForkHeight int
	HardForkHeight int
	DifficultyOld  string
	DifficultyNew  string
}
