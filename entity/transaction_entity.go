package entity

import (
	"database/sql"
	"time"
)

type Transaction struct {
	Id                int64               `gorm:"column:id;primaryKey;not null;autoIncrement"`
	Invoice           string              `gorm:"column:invoice;not null"`
	Total             int                 `gorm:"column:total;not null"`
	CreatedAt         time.Time           `gorm:"column:created_at;not null;autoCreateTime"`
	UpdateAt          time.Time           `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt         sql.NullTime        `gorm:"column:deleted_at;default=NULL"`
	AdminId           int64               `gorm:"column:admin_id;not null"`
	Admin             User                `gorm:"foreignKey:AdminId;references:Id"`
	TransactionDetail []TransactionDetail `gorm:"foreignKey:Id"`
}

func (Transaction) TableName() string {
	return "transactions"
}
