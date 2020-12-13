package persistance

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatastore ready to use DB and Client
type MongoDatastore struct {
	*mongo.Database
	*mongo.Client
	context.Context
}

var connectOnce sync.Once
var datastore *MongoDatastore

// NewMongoDatastore creates new NewMongoDatastore
func NewMongoDatastore(ctx context.Context, connectionString string) (ds *MongoDatastore, err error) {
	connectOnce.Do(func() {
		if datastore != nil {
			return
		}
		clientOptions := options.Client().ApplyURI(connectionString).SetDirect(true)
		c, err := mongo.NewClient(clientOptions)
		if err != nil {
			err = errors.Wrap(err, "creating mongodb client failed")
			return
		}
		err = c.Connect(ctx)
		if err != nil {
			err = errors.Wrap(err, "connectiong to mongodb failed")
			return
		}

		err = c.Ping(ctx, nil)
		if err != nil {
			err = errors.Wrap(err, "pinging to mongodb failed")
			return
		}

		db := c.Database("cosmosdb-parklakes-viberbot")
		datastore = &MongoDatastore{
			Database: db,
			Client:   c,
			Context:  ctx,
		}
	})
	return datastore, err
}

// GetDatastore returns datastore instance
func GetDatastore() *MongoDatastore {
	return datastore
}

// Ping checks DB avaliability
func (r MongoDatastore) Ping() error {
	return r.Client.Ping(r.Context, nil)
}

// Disconnect disconnects from datastore
func (r MongoDatastore) Disconnect() error {
	return r.Client.Disconnect(r.Context)
}

func (r MongoDatastore) findOne(col *mongo.Collection, id string, e interface{}) error {
	f := col.FindOne(r.Context, bson.M{"_id": id})

	if err := f.Err(); err != nil {
		return errors.Wrapf(err, "getting entity with id %s failed", id)
	}
	if err := f.Decode(e); err != nil {
		return errors.Wrapf(err, "decoding entity with id %s failed", id)
	}
	return nil
}

func (r MongoDatastore) upsert(col *mongo.Collection, id string, e interface{}) error {
	opts := options.Replace().SetUpsert(true)
	if _, err := col.ReplaceOne(r.Context, bson.M{"_id": id}, e, opts); err != nil {
		return errors.Wrapf(err, "upsert failed for entity with id %s", id)
	}
	return nil
}
