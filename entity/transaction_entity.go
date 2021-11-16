package entity

import (
	"database/sql"
	"time"
)

type Transaction struct {
	Id                 int64        `gorm:"column:id;primaryKey;not null;autoIncrement"`
	Invoice            string       `gorm:"column:invoice;not null"`
	GrandTotal         int          `gorm:"column:grand_total;not null"`
	CreatedAt          time.Time    `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt          time.Time    `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt          sql.NullTime `gorm:"column:deleted_at;default=NULL"`
	AdminId            int64        `gorm:"column:admin_id;not null"`
	Admin              User         `gorm:"foreignKey:AdminId;references:Id"`
	TransactionDetails []TransactionDetail
}

func (Transaction) TableName() string {
	return "transactions"
}
