package model

type TransactionDetailResponse struct {
	Id        int64  `json:"id"`
	ProductId int64  `json:"product_id"`
	Name      string `json:"name"`
	ImageUrl  string `json:"image_url"`
	Price     int    `json:"price"`
	Qty       int    `json:"qty"`
	SubTotal  int    `json:"sub_total"`
}

type TransactionResponse struct {
	Invoice    string                      `json:"invoice"`
	Details    []TransactionDetailResponse `json:"details"`
	GrandTotal int                         `json:"grand_total"`
	AdminId    int64                       `json:"admin_id"`
	AdminName  string                      `json:"admin_name"`
	CreatedAt  string                      `json:"created_at"`
	UpdatedAt  string                      `json:"updated_at"`
}
