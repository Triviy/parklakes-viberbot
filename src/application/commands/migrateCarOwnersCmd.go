package commands

import (
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/application/integrations"
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
	carOwnersRepo        interfaces.Repo
	carOwnerPropsRepo    interfaces.Repo
	carOwnersSpreadsheet integrations.GoogleSpreadsheet
}

// NewMigrateCarOwnersCmd creates new instance of MigrateCarOwnersCmd
func NewMigrateCarOwnersCmd(
	carOwnersRepo interfaces.Repo,
	carOwnerPropsRepo interfaces.Repo,
	carOwnersSpreadsheet integrations.GoogleSpreadsheet,
) *MigrateCarOwnersCmd {
	return &MigrateCarOwnersCmd{
		carOwnersRepo,
		carOwnerPropsRepo,
		carOwnersSpreadsheet,
	}
}

// Migrate gets car owners data from Google SpreadSheet and stores it in DB
func (cmd MigrateCarOwnersCmd) Migrate() error {
	migrationTime, err := services.GetKyivTime()
	if err != nil {
		return err
	}
	lastMigrationTime, err := cmd.getLastMigration()
	if err != nil {
		return err
	}

	data, err := cmd.carOwnersSpreadsheet.ReadSpreadsheetRange("A2:E")
	if err != nil {
		return err
	}
	log.Printf("Got %v records:\n", len(data))
	cos, err := prepareCarOwnersForSave(lastMigrationTime, data)
	if err != nil {
		return err
	}
	log.Printf("Got %v new car owners\n", len(cos))

	if err := cmd.runDBMigration(cos); err != nil {
		return err
	}

	if err := cmd.setLastMigration(migrationTime); err != nil {
		return err
	}
	log.Printf("Last migration time is set to %s\n", migrationTime)
	return nil
}

func prepareCarOwnersForSave(lastMigrationTime time.Time, data [][]interface{}) (cos map[string]models.CarOwner, err error) {
	cos = make(map[string]models.CarOwner)
	for _, row := range data {
		co, err := createCarOwnerFromRecord(row)
		if err != nil {
			return nil, err
		}
		created, err := services.ToKyivTime(co.Created)
		if err != nil {
			log.Printf("Failed to convert created time of %s to Kyiv time: %v\n", co.ID, err)
		}
		if created.After(lastMigrationTime) {
			cos[co.ID] = *co
		}
	}
	return
}

func createCarOwnerFromRecord(record []interface{}) (co *models.CarOwner, err error) {
	carNumber := strings.TrimSpace(record[1].(string))
	if carNumber == "" {
		err = errors.New("Car number is empty")
		return
	}
	carOwner := models.CarOwner{
		ID:        models.NormalizeCarNumber(carNumber),
		CarNumber: carNumber,
		Created:   record[0].(string),
		Owner:     record[2].(string),
	}
	firstPhone := strings.TrimSpace(record[3].(string))
	if firstPhone != "" {
		carOwner.Phones = append(carOwner.Phones, firstPhone)
	}
	if len(record) == 5 {
		secondPhone := strings.TrimSpace(record[4].(string))
		if secondPhone != "" {
			carOwner.Phones = append(carOwner.Phones, secondPhone)
		}
	}
	co = &carOwner
	return
}

func (cmd MigrateCarOwnersCmd) runDBMigration(cos map[string]models.CarOwner) error {
	for _, co := range cos {
		if err := cmd.carOwnersRepo.Upsert(co.ID, co); err != nil {
			return errors.Wrapf(err, "Upsert failed for %s", co.ID)
		}
		log.Printf("%v migrated\n", co.ID)
	}
	return nil
}

func (cmd MigrateCarOwnersCmd) getLastMigration() (time.Time, error) {
	var prop models.CarOwnerProp
	err := cmd.carOwnerPropsRepo.FindOne(lastMigrationTimeProp, prop)
	if err != nil {
		return services.ToKyivTime(defaultLastMigrationTime)
	}
	return services.ToKyivTime(prop.Value)
}

func (cmd MigrateCarOwnersCmd) setLastMigration(kt time.Time) error {
	prop := models.CarOwnerProp{
		ID:    lastMigrationTimeProp,
		Value: services.ToKyivFormat(kt),
	}
	if err := cmd.carOwnerPropsRepo.Upsert(lastMigrationTimeProp, prop); err != nil {
		return errors.Wrapf(err, "upserting %s failed", lastMigrationTimeProp)
	}
	return nil
}
