package config

import (
	"os"

	"github.com/joho/godotenv"
)

type App struct {
	Server struct {
		Host   string
		Port   string
		Secret string
		Origin string
	}

	Database struct {
		Host string
		Port string
		Name string
		User string
		Pass string
	}

	RabbitMQ struct {
		Host string
		Port string
	}
}

var app *App

func GetConfig() *App {
	if app == nil {
		app = initConfig()
	}

	return app
}

func initConfig() *App {
	conf := App{}
	err := godotenv.Load()
	if err != nil {
		conf.Server.Host = "localhost"
		conf.Server.Port = "10000"

		conf.Database.Host = "db"
		conf.Database.Port = "5432"
		conf.Database.User = "root"
		conf.Database.Pass = ""

		conf.RabbitMQ.Host = "localhost"
		conf.RabbitMQ.Port = ""

		return &conf
	}

	conf.Server.Host = os.Getenv("SERVER_HOST")
	conf.Server.Port = os.Getenv("SERVER_PORT")
	conf.Server.Secret = os.Getenv("SERVER_SECRET")
	conf.Server.Origin = os.Getenv("SERVER_ORIGIN")

	conf.Database.Host = os.Getenv("POSTGRES_HOST")
	conf.Database.Port = os.Getenv("POSTGRES_PORT")
	conf.Database.User = os.Getenv("POSTGRES_USER")
	conf.Database.Name = os.Getenv("POSTGRES_NAME")
	conf.Database.Pass = os.Getenv("POSTGRES_PASS")

	conf.RabbitMQ.Host = os.Getenv("RABBITMQ_HOST")
	conf.RabbitMQ.Port = os.Getenv("RABBITMQ_PORT")

	return &conf
}
