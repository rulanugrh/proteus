package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Appconfig struct {
	App struct {
		Port string
		Host string
	}

	Database struct {
		Port     string
		Name     string
		Password string
		User     string
		Host     string
	}
}

var appconfig *Appconfig

func Get() *Appconfig {
	if appconfig == nil {
		appconfig = ConfigInit()
	}
	return appconfig
}

func ConfigInit() *Appconfig {
	err := godotenv.Load()

	configs := Appconfig{}
	if err != nil {
		configs.Database.User = "root"
		configs.Database.Name = "tokoku"
		configs.Database.Password = "12345"
		configs.Database.Host = "localhost"
		configs.Database.Port = "3306"
		configs.App.Host = "localhost"
		configs.App.Port = "8080"
	}

	configs.Database.User = os.Getenv("DB_USER")
	configs.Database.Password = os.Getenv("DB_PASS")
	configs.Database.Host = os.Getenv("DB_HOST")
	configs.Database.Port = os.Getenv("DB_PORT")
	configs.Database.Name = os.Getenv("DB_NAME")
	configs.App.Host = os.Getenv("APP_HOST")
	configs.App.Port = os.Getenv("APP_PORT")

	return &configs
}
