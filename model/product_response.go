package model

type ProductResponse struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	ImageUrl  string `json:"image_url"`
	Price     int    `json:"price"`
	Stock     int    `json:"stock"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type ProductListResponse struct {
	TotalData    int               `json:"total_data"`
	TotalPage    int               `json:"total_page"`
	CurrentPage  int               `json:"current_page"`
	LimitPerPage int               `json:"limit_per_page"`
	Products     []ProductResponse `json:"products"`
}
