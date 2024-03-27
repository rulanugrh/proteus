package config

import (
	"os"

	"github.com/joho/godotenv"
)

type App struct {
	Server struct {
		HTTP   string
		GRPC   string
		Host   string
		Secret string
	}

	Database struct {
		Name string
		Host string
		Port string
		User string
		Pass string
	}

	RabbitMQ struct {
		Host string
		Port string
		User string
		Pass string
	}
	
	Xendit struct {
		SuccessURL string
		FailureURL string
		CancelURL string
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
		conf.Server.HTTP = "3000"
		conf.Server.GRPC = "8000"
		conf.Server.Host = "localhost"
		conf.Server.Secret = "secret"

		conf.Database.Host = "localhost"
		conf.Database.Port = "5431"
		conf.Database.User = "user"
		conf.Database.Pass = "user"
		conf.Database.Name = "order.db"
		
		conf.RabbitMQ.Host = "localhost"
		conf.RabbitMQ.Pass = "password"
		conf.RabbitMQ.User = "user"
		conf.RabbitMQ.Port = ""

		return &conf
	}

	conf.Server.GRPC = os.Getenv("GRPC_PORT")
	conf.Server.HTTP = os.Getenv("HTTP_PORT")
	conf.Server.Host = os.Getenv("APP_HOST")
	conf.Server.Secret = os.Getenv("APP_SECRET")

	conf.Database.Host = os.Getenv("POSTGRES_HOST")
	conf.Database.Pass = os.Getenv("POSTGRES_PASS")
	conf.Database.User = os.Getenv("POSTGRES_USER")
	conf.Database.Port = os.Getenv("POSTGRES_PORT")
	conf.Database.Name = os.Getenv("POSTGRESDB_NAME")

	conf.RabbitMQ.Host = os.Getenv("RABBITMQ_HOST")
	conf.RabbitMQ.Port = os.Getenv("RABBITMQ_PORT")
	conf.RabbitMQ.User = os.Getenv("RABBITMQ_USER")
	conf.RabbitMQ.Pass = os.Getenv("RABBITMQ_PASS")

	conf.Xendit.SuccessURL = os.Getenv("XENDIT_SUCCESS_URL")
	conf.Xendit.FailureURL = os.Getenv("XENDIT_FAILURE_URL")
	conf.Xendit.CancelURL = os.Getenv("XENDIT_CANCEL_URL")

	return &conf
}