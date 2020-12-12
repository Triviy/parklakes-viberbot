package commands

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/domain/interfaces"
	"github.com/triviy/parklakes-viberbot/domain/models"
	"github.com/triviy/parklakes-viberbot/domain/services"
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
func (cmd MigrateCarOwnersCmd) MigrateCardOwners(migrationTime time.Time, cos []models.CarOwner) error {
	prop, err := cmd.getLastMigrationProp()
	if err != nil {
		return errors.Wrap(err, "getting last migration time failed")
	}

	lastMigrationTime, err := services.ToKyivTime(prop.Value)
	if err != nil {
		return errors.Wrapf(err, "coverting %s to Kyiv time failed", lastMigrationTimeProp)
	}

	if err := cmd.migrate(lastMigrationTime, cos); err != nil {
		return errors.Wrap(err, "migrating of car owners failed")
	}

	if err := cmd.setLastMigrationProp(migrationTime); err != nil {
		return errors.Wrap(err, "updating last migration time failed")
	}

	log.Printf("Last migration time is set to %s\n", migrationTime)
	return nil
}

func (cmd MigrateCarOwnersCmd) migrate(lastMigrationTime time.Time, cos []models.CarOwner) error {
	for _, co := range cos {
		ktCreated, err := services.ToKyivTime(co.Created)
		if err != nil {
			log.Printf("Failed to convert created time of %s to Kyiv time: %v\n", co.ID, err)
		}
		if ktCreated.Before(lastMigrationTime) {
			continue
		}
		if err := cmd.carOwnersRepo.Upsert(co.ID, co); err != nil {
			return errors.Wrapf(err, "Upsert failed for %s", co.ID)
		}
		log.Printf("%v migrated\n", co.ID)
	}
	return nil
}

func (cmd MigrateCarOwnersCmd) getLastMigrationProp() (*models.CarOwnerProp, error) {
	var prop models.CarOwnerProp
	err := cmd.carOwnerPropsRepo.FindOne(lastMigrationTimeProp, prop)
	return &prop, err
}

func (cmd MigrateCarOwnersCmd) setLastMigrationProp(kt time.Time) error {
	prop := models.CarOwnerProp{
		ID:    lastMigrationTimeProp,
		Value: services.ToKyivFormat(kt),
	}
	return cmd.carOwnerPropsRepo.Upsert(lastMigrationTimeProp, prop)
}
