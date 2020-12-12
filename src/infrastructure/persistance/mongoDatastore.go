package persistance

import (
	"context"
	"log"
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

// NewMongoDatastore creates new NewMongoDatastore
func NewMongoDatastore(ctx context.Context, connectionString string) *MongoDatastore {
	var db *mongo.Database
	var c *mongo.Client
	connectOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(connectionString).SetDirect(true)
		session, err := mongo.NewClient(clientOptions)
		if err != nil {
			log.Fatal(err)
		}
		err = session.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}

		err = c.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("unable to connect %v", err)
		}

		db = session.Database("cosmosdb-parklakes-viberbot")
	})

	return &MongoDatastore{
		Database: db,
		Client:   c,
		Context:  ctx,
	}
}

func (r MongoDatastore) findOne(col *mongo.Collection, id string, e interface{}) error {
	f := col.FindOne(r.Context, bson.M{"_id": id})

	if err := f.Err(); err != nil {
		return errors.Wrapf(err, "getting card owner %v failed", id)
	}
	if err := f.Decode(&e); err != nil {
		return errors.Wrapf(err, "decoding card owner %v failed", id)
	}
	return nil
}

func (r MongoDatastore) upsert(col *mongo.Collection, id string, e interface{}) error {
	opts := options.Replace().SetUpsert(true)
	if _, err := col.ReplaceOne(r.Context, bson.M{"_id": id}, e, opts); err != nil {
		return errors.Wrapf(err, "Upsert failed for %s", id)
	}
	return nil
}
