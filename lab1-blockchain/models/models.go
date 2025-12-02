package models

type StudentRecord struct {
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
