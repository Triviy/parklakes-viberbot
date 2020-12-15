package persistance

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CarOwnersRepo implementation
type CarOwnersRepo struct {
	datastore *MongoDatastore
	*mongo.Collection
}

// NewCarOwnersRepo creates new instance of CarOwnersRepo
func NewCarOwnersRepo(ds *MongoDatastore) *CarOwnersRepo {
	return &CarOwnersRepo{
		ds,
		ds.Database.Collection("parklakes-car-owners"),
	}
}

// GetCollection returns car owners collection
func (r CarOwnersRepo) GetCollection() (*mongo.Collection, context.Context) {
	return r.Collection, datastore.Context
}

// FindOne returns car owners data from database by id
func (r CarOwnersRepo) FindOne(id string, e interface{}, opts ...*options.FindOneOptions) error {
	return r.datastore.findOne(r.Collection, id, e)
}

// Upsert inserts or updates car owner data
func (r CarOwnersRepo) Upsert(id string, e interface{}) error {
	return r.datastore.upsert(r.Collection, id, e)
}

// UpdateOne updates props properties
func (r CarOwnersRepo) UpdateOne(id string, u primitive.M) error {
	return r.datastore.updateOne(r.Collection, id, u)
}

// DeleteOne deletes car owner by id
func (r CarOwnersRepo) DeleteOne(id string) error {
	return r.datastore.deleteOne(r.Collection, id)
}
