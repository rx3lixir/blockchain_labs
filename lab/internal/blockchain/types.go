package blockchain

const (
	Difficulty = "00"
)

type StudentRecord struct {
	ID       string
	FullName string
	Zachetka string
	Group    string
	Subject  string
	Course   int
	Grade    int
}

type Block struct {
	Index        int
	Timestamp    int64
	Data         StudentRecord
	PreviousHash string
	Hash         string
	Nonce        int
}
