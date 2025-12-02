package models

// ForkVersion
const (
	ForkVersionOriginal = 0
	ForkVersionSoft     = 1
	ForkVersionHard     = 2
)

type StudentRecord struct {
	FullName string
	Zachetka string
	Group    string
	Subject  string
	Course   int
	Grade    int
	Teacher  string `json:"Teacher,omitempty"` // Новое поле для Soft Fork - добавлено с блока #5
}

type Block struct {
	Index        int
	Timestamp    int64
	Data         StudentRecord
	PreviousHash string
	Hash         string
	Nonce        int
	ForkVersion  int `json:"ForkVersion,omitempty"` // 0=original, 1=soft, 2=hard
}
