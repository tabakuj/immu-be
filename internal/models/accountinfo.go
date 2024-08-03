package models

type AccountType int

const (
	Sending   AccountType = iota + 1 //  1
	Receiving                        // 2
)

type AccountInfo struct {
	Id      uint // this name is used so we can have a index on this
	Name    string
	Iban    string
	Address *string // this is used for simplicity normally it needs to be a little bit more complex
	Amount  float64 // Maybe It would be better if int64
	Type    *AccountType
}
