package model

type HistoryRequestQuery struct {
	Q         string `query:"q"`
	Page      int    `query:"page"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
}

type HistorySelectQuery struct {
	Search    string
	Start     int
	Limit     int
	StartDate string
	EndDate   string
}
