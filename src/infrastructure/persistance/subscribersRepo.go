package persistance

import (
	"go.mongodb.org/mongo-driver/mongo"
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

// FindOne returns subscribers data from database by id
func (r SubscribersRepo) FindOne(id string, e interface{}) error {
	return r.datastore.findOne(r.Collection, id, e)
}

// Upsert inserts or updates subscribers data
func (r SubscribersRepo) Upsert(id string, e interface{}) error {
	return r.datastore.upsert(r.Collection, id, e)
}

// DeleteOne deletes subscriber by id
func (r SubscribersRepo) DeleteOne(id string) error {
	return r.datastore.deleteOne(r.Collection, id)
}
