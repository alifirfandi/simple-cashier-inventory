package entity

type TransactionDetail struct {
	Id            int64       `gorm:"column:id;primaryKey;not null;autoIncrement"`
	Qty           int         `gorm:"column:qty;not null"`
	SubTotal      int         `gorm:"column:sub_total;not null"`
	ProductId     int64       `gorm:"column:product_id;not null"`
	TransactionId int64       `gorm:"column:transaction_id;not null"`
	Product       Product     `gorm:"foreignKey:ProductId;references:Id"`
	Transaction   Transaction `gorm:"foreignKey:TransactionId;references:Id"`
}

func (TransactionDetail) TableName() string {
	return "transaction_details"
}
