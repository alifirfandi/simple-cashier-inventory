package model

type TransactionRequest struct {
	AdminId   int64
	AdminName string
	Details   []CartRequest `json:"details"`
}
