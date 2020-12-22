package interfaces

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GenericRepo generic interface
type GenericRepo interface {
	GetCollection() (*mongo.Collection, context.Context)
	FindOne(id string, e interface{}, opts ...*options.FindOneOptions) error
	Upsert(id string, e interface{}) error
	UpdateOne(id string, u primitive.M) error
	DeleteOne(id string) error
}
