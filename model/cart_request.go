package model

type CartRequest struct {
	AdminId   int64
	ProductId int64 `json:"product_id"`
	Qty       int   `json:"qty"`
}
