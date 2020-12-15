package commands

import (
	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/application/integrations/viber"
	"github.com/triviy/parklakes-viberbot/domain/interfaces"
	"github.com/triviy/parklakes-viberbot/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UpdateSubscriberCmd instance of viber webhook cmd
type UpdateSubscriberCmd struct {
	subscriberRepo interfaces.GenericRepo
}

// NewUpdateSubscriberCmd creates new instance of UpdateSubscriberCmd
func NewUpdateSubscriberCmd(subscriberRepo interfaces.GenericRepo) *UpdateSubscriberCmd {
	return &UpdateSubscriberCmd{subscriberRepo}
}

// Execute calls setting Viber callback URLs
func (cmd UpdateSubscriberCmd) Execute(user *viber.User, contact *viber.Contact) error {
	var phonesProjection map[string][]string
	opts := options.FindOne()
	opts.Projection = bson.M{"phoneNumbers": 1, "_id": 0}

	if err := cmd.subscriberRepo.FindOne(user.ID, &phonesProjection, opts); err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
	}
	newSub := models.Subscriber{
		ID:      user.ID,
		Name:    user.Name,
		Avatar:  user.Avatar,
		Country: user.Country,
		Active:  true,
	}
	if val, ok := phonesProjection["phoneNumbers"]; ok && len(phonesProjection["phoneNumbers"]) > 0 {
		copy(val, newSub.PhoneNumbers)
	}
	if len(contact.PhoneNumber) > 5 && !contains(newSub.PhoneNumbers, contact.PhoneNumber) {
		newSub.PhoneNumbers = append(newSub.PhoneNumbers, contact.PhoneNumber)
	}

	if err := cmd.subscriberRepo.Upsert(user.ID, &newSub); err != nil {
		return err
	}
	return nil
}

func contains(slice []string, search string) bool {
	for _, v := range slice {
		if v == search {
			return true
		}
	}
	return false
}
