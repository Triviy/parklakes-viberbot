package commands

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/domain/interfaces"
	"github.com/triviy/parklakes-viberbot/domain/services"
	"github.com/triviy/parklakes-viberbot/models"
	"github.com/triviy/parklakes-viberbot/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	defaultLastMigrationTime = "01.01.2017 00:00:00"
	lastMigrationTimeProp    = "last-migration-time"
)

// MigrateCarOwnersCmd runs migration
type MigrateCarOwnersCmd struct {
	carOwnersRepo     interfaces.Repo
	carOwnerPropsRepo interfaces.Repo
}

// NewMigrateCarOwnersCmd creates new instance of MigrateCarOwnersCmd
func NewMigrateCarOwnersCmd(cor interfaces.Repo, copr interfaces.Repo) *MigrateCarOwnersCmd {
	return &MigrateCarOwnersCmd{
		cor,
		copr,
	}
}

// MigrateCardOwners migrates car owners to database
func (c MigrateCarOwnersCmd) MigrateCardOwners(kt time.Time, cos map[string]models.CarOwner) error {
	prop := getLastMigrationProp(ctx, c)
	lastMigrationTime, err := utils.ToKyivTime(prop.Value)
	if err != nil {
		return errors.Wrapf(err, "coverting %s to Kyiv time failed", lastMigrationTimeProp)
	}
	opts := options.Replace().SetUpsert(true)
	coCol := c.Database(database).Collection(coCollection)

	for _, co := range cos {
		ktCreated, err := utils.ToKyivTime(co.Created)
		if err != nil {
			log.Printf("Failed to convert created time of %s to Kyiv time: %v\n", co.ID, err)
		}
		if ktCreated.Before(lastMigrationTime) {
			continue
		}
		if _, err := coCol.ReplaceOne(ctx, bson.M{"_id": co.ID}, co, opts); err != nil {
			return errors.Wrapf(err, "Upsert failed for %s", co.ID)
		}
		log.Printf("%v migrated\n", co.ID)
	}
	setLastMigrationProp(ctx, c, kt)
	return nil
}

func (c MigrateCarOwnersCmd) getLastMigrationProp() (*models.CarOwnerProp, error) {
	var prop models.CarOwnerProp
	err := c.carOwnerPropsRepo.FindOne(lastMigrationTimeProp, prop)
	return prop, err
}

func setLastMigrationProp(ctx context.Context, c *mongo.Client, kt time.Time) error {
	prop := models.CarOwnerProp{
		ID:    lastMigrationTimeProp,
		Value: services.ToKyivFormat(kt),
	}
	err := c.carOwnerPropsRepo.Upsert(lastMigrationTimeProp, prop)
	log.Printf("Last migration time is set to %s\n", kt)
	return cop, err
}
