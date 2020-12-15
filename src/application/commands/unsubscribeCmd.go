package commands

import (
	"github.com/triviy/parklakes-viberbot/domain/interfaces"
	"go.mongodb.org/mongo-driver/bson"
)

// UnsubscribeCmd instance of viber webhook cmd
type UnsubscribeCmd struct {
	subscriberRepo interfaces.GenericRepo
}

// NewUnsubscribeCmd creates new instance of UnsubscribeCmd
func NewUnsubscribeCmd(carOwnersRepo interfaces.GenericRepo) *UnsubscribeCmd {
	return &UnsubscribeCmd{carOwnersRepo}
}

// Execute calls setting Viber callback URLs
func (cmd UnsubscribeCmd) Execute(userID string) error {
	if err := cmd.subscriberRepo.UpdateOne(userID, bson.M{"active": false}); err != nil {
		return err
	}
	return nil
}
