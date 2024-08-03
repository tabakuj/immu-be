package models

type AccountType int

const (
	Sending   AccountType = iota + 1 //  1
	Receiving                        // 2
)

type AccountInfo struct {
	AccountNumber uint         `json:"account_number"`
	AccountName   string       `json:"account_name"`
	Iban          string       `json:"iban"`
	Address       *string      `json:"address"` // this is used for simplicity normally it needs to be a little bit more complex
	Amount        float64      `json:"amount"`  // Maybe It would be better if int64
	Type          *AccountType `json:"type"`
}
