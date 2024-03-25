package config

import (
	"fmt"
	"log"

	"github.com/rulanugrh/order/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB
	conf *App
}

func InitializeDB(conf *App) *Postgres {
	return &Postgres{conf: conf}
}

func (p *Postgres) StartConnection() error {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=Asia/Jakarta",
		p.conf.Database.User,
		p.conf.Database.Pass,
		p.conf.Database.Host,
		p.conf.Database.Port,
		p.conf.Database.Name,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Printf("error cant connect to DB %s", err.Error())
		return err
	}

	p.DB = db
	
	log.Println("success connect to database")
	return nil
}

func (p *Postgres) Migrate() error {
	err := p.DB.AutoMigrate(&entity.Product{}, &entity.Order{})
	if err != nil {
		log.Printf("error cannot migrate to DB %s", err.Error())
		return err
	}
	log.Println("success migrate to DB")
	return nil
}