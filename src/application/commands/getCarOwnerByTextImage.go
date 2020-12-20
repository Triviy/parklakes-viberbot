package commands

import (
	"fmt"
	"strings"

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
	cmd.imageTextReader.BatchReadFileRemoteImage(cm.Media)
	return nil
	// text, err := cmd.getUsersResponseByText(cm.Text)
	// if err != nil {
	// 	return err
	// }

	// request := viber.MessageRequest{
	// 	Receiver:     userID,
	// 	Type:         viber.TextType,
	// 	Text:         text,
	// 	TrackingData: cm.TrackingData,
	// }

	// return viber.SendMessage(&request, cmd.config.GetViberBaseURL(), cmd.config.GetViberAPIKey())
}

func (cmd GetCarOwnerByImageCmd) getUsersResponseByText(input string) (text string, err error) {
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
	var sb strings.Builder
	sb.WriteString("Я знайшов 😄")
	if co.Owner != "" {
		sb.WriteString(fmt.Sprintf("\nВласник автівки %s", co.Owner))
	}
	sb.WriteString(fmt.Sprintf("\nНомер телефону: %s", co.Phones[0]))
	if len(co.Phones) > 1 {
		sb.WriteString(fmt.Sprintf("\nДодатковый: %s", co.Phones[1]))
	}
	return sb.String(), nil
}
