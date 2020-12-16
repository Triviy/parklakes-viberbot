package commands

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/application/integrations/viber"
	"github.com/triviy/parklakes-viberbot/domain/interfaces"
	"github.com/triviy/parklakes-viberbot/domain/models"
	"github.com/triviy/parklakes-viberbot/web/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetCarOwnerByTextCmd instance of viber webhook cmd
type GetCarOwnerByTextCmd struct {
	config        *config.APIConfig
	carOwnersRepo interfaces.GenericRepo
}

// NewGetCarOwnerByTextCmd creates new instance of GetCarOwnerByTextCmd
func NewGetCarOwnerByTextCmd(config *config.APIConfig, carOwnersRepo interfaces.GenericRepo) *GetCarOwnerByTextCmd {
	return &GetCarOwnerByTextCmd{config, carOwnersRepo}
}

// Execute calls setting Viber callback URLs
func (cmd GetCarOwnerByTextCmd) Execute(cm *viber.CallbackMessage, userID string) error {
	if cm == nil {
		return errors.New("viber.CallbackMessage is nil")
	}
	text, err := cmd.getUsersResponseByText(cm.Text)
	if err != nil {
		return err
	}

	request := viber.MessageRequest{
		AuthToken:    cmd.config.GetViberAPIKey(),
		Receiver:     userID,
		Type:         viber.TextType,
		Text:         text,
		TrackingData: cm.TrackingData,
	}

	return viber.SendMessage(&request, cmd.config.GetViberBaseURL())
}

func (cmd GetCarOwnerByTextCmd) getUsersResponseByText(input string) (text string, err error) {
	carNumber := models.NormalizeCarNumber(input)
	if len(carNumber) < 3 || len(carNumber) > 16 {
		return "Вибачте, отриманий номер автівки замалий або завеликий. Спробуйте ще 😉", nil
	}
	var co models.CarOwner
	if err := cmd.carOwnersRepo.FindOne(carNumber, &co); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "Вибачте, мені не вдалося знайти автівки з таким номером 😥", nil
		}
		return "", err
	}
	text = fmt.Sprintf("Я знайшов 😄\nВласник автівки %s\nНомер телефону: %s", co.Owner, co.Phones[0])
	if len(co.Phones) > 1 {
		text += fmt.Sprintf("Додатковый: %s", co.Phones[1])
	}
	return text, nil
}
