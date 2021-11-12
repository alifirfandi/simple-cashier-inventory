package entity

type User struct {
	Id       int64  `gorm:"column:id"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

func (User) TableName() string {
	return "users"
}
