package entity

import (
	"database/sql"
	"time"
)

type User struct {
	Id        int64        `gorm:"column:id"`
	Name      string       `gorm:"column:name"`
	Email     string       `gorm:"column:email"`
	Password  string       `gorm:"column:password"`
	Role      string       `gorm:"column:role;default=ADMIN"`
	CreatedAt time.Time    `gorm:"column:created_at"`
	UpdateAt  time.Time    `gorm:"column:updated_at"`
	DeletedAt sql.NullTime `gorm:"column:deleted_at;default=NULL"`
}

func (User) TableName() string {
	return "users"
}
