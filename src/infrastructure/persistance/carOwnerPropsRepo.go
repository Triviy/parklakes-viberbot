package persistance

import (
	"go.mongodb.org/mongo-driver/mongo"
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

// FindOne returns database props by id
func (r CarOwnerPropsRepo) FindOne(id string, e interface{}) error {
	return r.datastore.findOne(r.Collection, id, e)
}

// Upsert inserts or updates database props
func (r CarOwnerPropsRepo) Upsert(id string, e interface{}) error {
	return r.datastore.upsert(r.Collection, id, e)
}
