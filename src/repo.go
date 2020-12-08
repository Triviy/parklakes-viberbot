package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultLastMigrationTime = "01.01.2017 00:00:00"
	database                 = "cosmosdb-parklakes-viberbot"
	coCollection             = "parklakes-car-owners"
	propsCollection          = "parklakes-car-owners-props"
	lastMigrationTimeProp    = "last-migration-time"
)

// GetCardOwner returns card owners data from database by id
func GetCardOwner(id string) *CarOwner {
	c := connect(time.Minute)
	ctx := context.Background()
	defer c.Disconnect(ctx)

	col := c.Database(database).Collection(coCollection)
	f := col.FindOne(ctx, bson.M{"_id": id})

	var co CarOwner
	if err := f.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Could not get car owner %s: %v\n", id, err)
			return nil
		}
		log.Fatalf("Error occured while getting card owner %v", err)
	}
	f.Decode(&co)
	log.Printf("Car owner found: %v\n", co)
	return &co
}

// MigrateCardOwners migrates car owners to database
func MigrateCardOwners(cos map[string]CarOwner) {
	c := connect(time.Minute * 10)
	ctx := context.Background()
	defer c.Disconnect(ctx)

	prop := getLastMigrationProp(ctx, c)
	lastMigrationTime := ToKyivTime(prop.Value)
	opts := options.Replace().SetUpsert(true)
	coCol := c.Database(database).Collection(coCollection)

	for _, co := range cos {
		if ToKyivTime(co.Created).Before(lastMigrationTime) {
			continue
		}
		_, err := coCol.ReplaceOne(ctx, bson.M{"_id": co.ID}, co, opts)
		if err != nil {
			log.Fatalf("Failed to migrate car owner: %v\n", err)
		}
		log.Printf("%v migrated\n", co.ID)
	}
	setLastMigrationProp(ctx, c)
}

func getLastMigrationProp(ctx context.Context, c *mongo.Client) *CarOwnerProp {
	propsCol := c.Database(database).Collection(propsCollection)
	prop := CarOwnerProp{
		ID: lastMigrationTimeProp,
	}
	f := propsCol.FindOne(ctx, bson.M{"_id": prop.ID})
	if f.Err() != nil {
		prop.Value = defaultLastMigrationTime
	} else {
		f.Decode(&prop)
	}
	log.Printf("Last migration time is set to %s\n", prop.Value)
	return &prop
}

func setLastMigrationProp(ctx context.Context, c *mongo.Client) {
	propsCol := c.Database(database).Collection(propsCollection)
	prop := CarOwnerProp{
		ID:    lastMigrationTimeProp,
		Value: GetKyivTime().Format(kyivFormat),
	}
	opts := options.Replace().SetUpsert(true)
	_, err := propsCol.ReplaceOne(ctx, bson.M{"_id": prop.ID}, prop, opts)
	if err != nil {
		log.Fatalf("Failed replace migration time: %v\n", err)
	}
	log.Printf("Migration props updated %v\n", prop)
}

func connect(timeout time.Duration) *mongo.Client {
	mongoDBConnectionString := GetDBConnectionString()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoDBConnectionString).SetDirect(true)
	c, err := mongo.NewClient(clientOptions)

	err = c.Connect(ctx)

	if err != nil {
		log.Fatalf("unable to initialize connection %v", err)
	}
	err = c.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("unable to connect %v", err)
	}
	return c
}
