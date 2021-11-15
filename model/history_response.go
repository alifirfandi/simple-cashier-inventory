package model

type HistoryResponse struct {
	Invoice    string                `json:"invoice"`
	Details    []TransactionResponse `json:"details"`
	GrandTotal int                   `json:"grand_total"`
	AdminId    int64                 `json:"admin_id"`
	AdminName  string                `json:"admin_name"`
	CreatedAt  string                `json:"created_at"`
	UpdatedAt  string                `json:"updated_at"`
}

type HistoryListResponse struct {
	TotalData    int               `json:"total_data"`
	TotalPage    int               `json:"total_page"`
	CurrentPage  int               `json:"current_page"`
	LimitPerPage int               `json:"limit_per_page"`
	Histories    []HistoryResponse `json:"histories"`
}
