package repository

import (
	"context"
	"time"

	"github.com/rulanugrh/order/internal/config"
	"github.com/rulanugrh/order/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductInterface interface {
	Create(req entity.Product) error
	FindID(id uint) (*entity.Product, error)
	Update(id uint, model entity.Product) error
}

type product struct {
	client *mongo.Collection
}

func ProductRepository(client *config.MongoDB, conf *config.App) ProductInterface {
	return &product{
		client: client.Conn.Database(conf.MongoDB.Name).Collection("product"),
	}
}

func(p *product) Create(req entity.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	_, err := p.client.InsertOne(ctx, &req)
	if err != nil {
		return err
	}

	return nil
}

func(p *product) FindID(id uint) (*entity.Product, error) {
	var response entity.Product
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	err := p.client.FindOne(ctx, bson.M{"id": id}).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (p *product) Update(id uint, model entity.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	err := p.client.FindOneAndUpdate(ctx, bson.M{"id": id}, bson.M{"$set": model}).Err()
	if err != nil {
		return err
	}

	return nil
}

func (p *product) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	_, err := p.client.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	return nil
}