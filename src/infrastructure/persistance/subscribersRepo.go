package persistance

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SubscribersRepo implementation
type SubscribersRepo struct {
	datastore *MongoDatastore
	*mongo.Collection
}

// NewSubscribersRepo creates new instance of CarOwnersRepo
func NewSubscribersRepo(ds *MongoDatastore) *SubscribersRepo {
	return &SubscribersRepo{
		ds,
		ds.Database.Collection("parklakes-subscribers"),
	}
}

// GetCollection returns props collection
func (r SubscribersRepo) GetCollection() (*mongo.Collection, context.Context) {
	return r.Collection, datastore.Context
}

// FindOne returns subscribers data from database by id
func (r SubscribersRepo) FindOne(id string, e interface{}, opts ...*options.FindOneOptions) error {
	return r.datastore.findOne(r.Collection, id, e, opts...)
}

// Upsert inserts or updates subscribers data
func (r SubscribersRepo) Upsert(id string, e interface{}) error {
	return r.datastore.upsert(r.Collection, id, e)
}

// UpdateOne updates subscribers properties
func (r SubscribersRepo) UpdateOne(id string, u primitive.M) error {
	return r.datastore.updateOne(r.Collection, id, u)
}

// DeleteOne deletes subscriber by id
func (r SubscribersRepo) DeleteOne(id string) error {
	return r.datastore.deleteOne(r.Collection, id)
}
