package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AppConfig struct {
	Database struct {
		Port     string
		Name     string
		Password string
		User     string
		Host     string
	}
}

var Dbconn *gorm.DB

func GetConnect() *gorm.DB {
	var config AppConfig

	err := godotenv.Load()
	if err != nil {
		config.Database.User = "root"
		config.Database.Name = "TokoKu"
		config.Database.Password = "12345"
		config.Database.Host = "localhost"
		config.Database.Port = "3306"
	}

	config.Database.User = os.Getenv("DB_USER")
	config.Database.Password = os.Getenv("DB_PASS")
	config.Database.Host = os.Getenv("DB_HOST")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.Name = os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name)

	db, errDb := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if errDb != nil {
		fmt.Println("cant connect database")
	}

	Dbconn = db

	return db
}
