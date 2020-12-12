package persistance

import (
	"go.mongodb.org/mongo-driver/mongo"
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

// FindOne returns car owners data from database by id
func (r CarOwnersRepo) FindOne(id string, e interface{}) error {
	return r.datastore.findOne(r.Collection, id, e)
}

// Upsert inserts or updates car owner data
func (r CarOwnersRepo) Upsert(id string, e interface{}) error {
	return r.datastore.upsert(r.Collection, id, e)
}
