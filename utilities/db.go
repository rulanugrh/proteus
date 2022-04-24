package utilities

import (
	"fmt"

	"github.com/ItsArul/TokoKu/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Dbconn *gorm.DB

func GetConnection() *gorm.DB {
	configs := config.Get()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		configs.Database.User,
		configs.Database.Password,
		configs.Database.Host,
		configs.Database.Port,
		configs.Database.Name)

	db, errDb := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errDb != nil {
		fmt.Println("cant connect database")
	}

	Dbconn = db

	return db
}
