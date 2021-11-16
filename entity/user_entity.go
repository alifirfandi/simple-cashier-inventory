package entity

import (
	"database/sql"
	"time"
)

type User struct {
	Id        int64        `gorm:"column:id;primaryKey;not null;autoIncrement"`
	Name      string       `gorm:"column:name;not null"`
	Email     string       `gorm:"column:email;not null"`
	Password  string       `gorm:"column:password;not null"`
	Role      string       `gorm:"column:role;not null;default=ADMIN"`
	CreatedAt time.Time    `gorm:"column:created_at;not null;autoCreateTime"`
	UpdateAt  time.Time    `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt sql.NullTime `gorm:"column:deleted_at;default=NULL"`
}

func (User) TableName() string {
	return "users"
}
