package model

type TransactionResponse struct {
	Id        int64  `json:"id"`
	ProductId int64  `json:"product_id"`
	Name      string `json:"name"`
	ImageUrl  string `json:"image_url"`
	Price     int    `json:"price"`
	Qty       int    `json:"qty"`
	SubTotal  int    `json:"sub_total"`
}
