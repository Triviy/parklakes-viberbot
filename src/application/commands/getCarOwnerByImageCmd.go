package commands

import (
	"fmt"

	"github.com/pkg/errors"
	computervision "github.com/triviy/parklakes-viberbot/application/integrations/computer-vision"
	"github.com/triviy/parklakes-viberbot/application/integrations/viber"
	"github.com/triviy/parklakes-viberbot/domain/interfaces"
	"github.com/triviy/parklakes-viberbot/domain/models"
	"github.com/triviy/parklakes-viberbot/web/config"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetCarOwnerByImageCmd instance of viber webhook cmd
type GetCarOwnerByImageCmd struct {
	config          *config.APIConfig
	carOwnersRepo   interfaces.GenericRepo
	imageTextReader *computervision.ImageTextReader
}

// NewGetCarOwnerByImageCmd creates new instance of GetCarOwnerByImageCmd
func NewGetCarOwnerByImageCmd(config *config.APIConfig, carOwnersRepo interfaces.GenericRepo, imageTextReader *computervision.ImageTextReader) *GetCarOwnerByImageCmd {
	return &GetCarOwnerByImageCmd{config, carOwnersRepo, imageTextReader}
}

// Execute calls setting Viber callback URLs
func (cmd GetCarOwnerByImageCmd) Execute(cm *viber.CallbackMessage, userID string) error {
	if cm == nil {
		return errors.New("viber.CallbackMessage is nil")
	}
	r, err := cmd.imageTextReader.BatchReadFileRemoteImage(cm.Media)
	if err != nil {
		return err
	}
	text, err := cmd.getUserResponse(r)
	if err != nil {
		return err
	}

	request := viber.MessageRequest{
		Receiver:     userID,
		Type:         viber.TextType,
		Text:         text,
		TrackingData: cm.TrackingData,
	}

	return viber.SendMessage(&request, cmd.config.GetViberBaseURL(), cmd.config.GetViberAPIKey())
}

func (cmd GetCarOwnerByImageCmd) getUserResponse(input []string) (text string, err error) {
	var carNumber string
	for _, t := range input {
		if cn, ok := models.ToCarNumber(t); ok {
			carNumber = cn
			break
		}
	}
	if carNumber == "" {
		return "Вибачте, не вдалось розпізнати номера автівки по фото. Спробуйте ще 😉", nil
	}
	var co models.CarOwner
	if err := cmd.carOwnersRepo.FindOne(carNumber, &co); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Sprintf("Вибачте, мені не вдалося знайти автівки з номером %s 😥", carNumber), nil
		}
		return "", err
	}
	return co.ToBotResponse(), nil
}
