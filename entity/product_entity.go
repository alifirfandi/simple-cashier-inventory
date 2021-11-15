package entity

import (
	"database/sql"
	"time"
)

type Product struct {
	Id        int64        `gorm:"column:id;primaryKey;not null;autoIncrement"`
	Name      string       `gorm:"column:name;not null"`
	ImageUrl  string       `gorm:"column:image_url"`
	Price     int          `gorm:"column:price;not null"`
	Stock     int          `gorm:"column:stock;not null"`
	CreatedAt time.Time    `gorm:"column:created_at;not null;autoCreateTime"`
	UpdateAt  time.Time    `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt sql.NullTime `gorm:"column:deleted_at;default=NULL"`
}

func (Product) TableName() string {
	return "products"
}
