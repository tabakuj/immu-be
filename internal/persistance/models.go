package persistance

import "immudb/internal/models"

type CreateResponse struct {
	DocumentID    string `json:"documentId"`
	TransactionID string `json:"transactionId"`
}

type GetAllRequest struct {
	Query   *Query `json:"query"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
}
type FieldComparison struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// Expression represents a set of field comparisons.
type Expression struct {
	FieldComparisons []FieldComparison `json:"fieldComparisons"`
}

// OrderBy represents the sorting order for fields.
type OrderBy struct {
	Desc  bool   `json:"desc"`
	Field string `json:"field"`
}

// Query represents the query structure with expressions and order by options.
type Query struct {
	Expressions []Expression `json:"expressions"`
	Limit       int          `json:"limit"`
	OrderBy     []OrderBy    `json:"orderBy"`
}

type GetAllResponse struct {
	Page      int         `json:"page"`
	PerPage   int         `json:"perPage"`
	Revisions []Revisions `json:"revisions"`
	SearchID  string      `json:"searchId"`
}
type VaultMd struct {
	Creator string `json:"creator"`
	Ts      int    `json:"ts"`
}

type Revisions struct {
	Document      models.AccountInfo `json:"document"` // if you are taking more than one document this needs to be either interface or Generic type
	Revision      string             `json:"revision"`
	TransactionID string             `json:"transactionId"`
}
