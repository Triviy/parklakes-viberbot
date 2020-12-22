package persistance

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CarOwnerPropsRepo implementation
type CarOwnerPropsRepo struct {
	datastore *MongoDatastore
	*mongo.Collection
}

// NewCarOwnerPropsRepo creates new instance of CarOwnerPropsRepo
func NewCarOwnerPropsRepo(ds *MongoDatastore) *CarOwnerPropsRepo {
	return &CarOwnerPropsRepo{
		ds,
		ds.Database.Collection("parklakes-car-owners-props"),
	}
}

// GetCollection returns props collection
func (r CarOwnerPropsRepo) GetCollection() (*mongo.Collection, context.Context) {
	return r.Collection, datastore.Context
}

// FindOne returns database props by id
func (r CarOwnerPropsRepo) FindOne(id string, e interface{}, opts ...*options.FindOneOptions) error {
	return r.datastore.findOne(r.Collection, id, e, opts...)
}

// Upsert inserts or updates database props
func (r CarOwnerPropsRepo) Upsert(id string, e interface{}) error {
	return r.datastore.upsert(r.Collection, id, e)
}

// UpdateOne updates car owners properties
func (r CarOwnerPropsRepo) UpdateOne(id string, u primitive.M) error {
	return r.datastore.updateOne(r.Collection, id, u)
}

// DeleteOne deletes props by id
func (r CarOwnerPropsRepo) DeleteOne(id string) error {
	return r.datastore.deleteOne(r.Collection, id)
}
