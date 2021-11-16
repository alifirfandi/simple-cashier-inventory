package model

type ProductRequest struct {
	Name     string `json:"name"`
	ImageUrl string `json:"-"`
	Price    int    `json:"price"`
	Stock    int    `json:"stock"`
}

type ProductRequestQuery struct {
	Q    string `query:"q"`
	Page int    `query:"page"`
	Sort string `query:"sort"`
}

type ProductSelectQuery struct {
	Search    string
	SortField string
	SortBy    string
	Start     int
	Limit     int
}
