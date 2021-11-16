package entity

import (
	"database/sql"
	"time"
)

type Cart struct {
	Id        int64        `gorm:"column:id;primaryKey;not null;autoIncrement"`
	Qty       int          `gorm:"column:qty;not null"`
	CreatedAt time.Time    `gorm:"column:created_at;not null;autoCreateTime"`
	UpdateAt  time.Time    `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt sql.NullTime `gorm:"column:deleted_at;default=NULL"`
	ProductId int64        `gorm:"foreignKey:ProductId;references:Id"`
	AdminId   int64        `gorm:"foreignKey:AdminId;references:Id"`
}

func (Cart) TableName() string {
	return "carts"
}
