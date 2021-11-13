package model

type ProductRequest struct {
	Name     string `json:"name"`
	ImageUrl string `json:"-"`
	Price    int    `json:"price"`
	Stock    int    `json:"stock"`
}
