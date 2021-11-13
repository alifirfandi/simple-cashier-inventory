package test

import (
	"database/sql"
	"os"

	"github.com/alifirfandi/simple-cashier-inventory/exception"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func dbInit() (gormDb *gorm.DB) {
	var err error
	err = godotenv.Load("../.env.test")
	exception.PanicIfNeeded(err)

	sqlDB, err := sql.Open("mysql", os.Getenv("MYSQL_HOST"))
	exception.PanicIfNeeded(err)

	gormDb, err = gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	exception.PanicIfNeeded(err)

	return
}

func clearDb(mysql *gorm.DB) {
	mysql.Exec("DELETE FROM carts")
	mysql.Exec("DELETE FROM transaction_details")
	mysql.Exec("DELETE FROM transactions")
	mysql.Exec("DELETE FROM products")
	mysql.Exec("DELETE FROM users")
	mysql.Exec(`INSERT INTO users VALUES (
		1,
		'Super Admin',
		'superadmin@golang.com',
		'$2a$14$Jwe3YurEbLcjVZTLqKUFHeLBKnWd2sV2dZd1UIqhmPRO4SqGVbCtS',
		'SUPERADMIN',
		'2021-11-13 07:26:59',
		'2021-11-13 07:26:59',
		NULL
	)`)
}
