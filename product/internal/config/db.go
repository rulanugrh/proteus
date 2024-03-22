package config

import (
	"fmt"
  "log"

	"gorm.io/gorm"
  "gorm.io/driver/postgres"
)

type Database struct {
	DB *gorm.DB
  conf *App
}

func InitializeDB(conf *App) *Database {
  return &Database{conf: conf}
}

func (p *Database) ConnectionDB() {
  dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Jakarta", p.conf.Database.User, p.conf.Database.Pass, p.conf.Database.Host, p.conf.Database.Port, p.conf.Database.Name)

  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

  if err != nil {
    log.Printf("error cant connect to DB %s", err.Error)
  }

  p.DB = db
}

