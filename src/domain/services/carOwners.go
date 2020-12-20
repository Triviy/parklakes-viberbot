package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/triviy/parklakes-viberbot/domain/interfaces"
	"github.com/triviy/parklakes-viberbot/domain/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUsersResponseByCarNumber searches for car owner in DB and returns valid user response
func GetUsersResponseByCarNumber(repo interfaces.GenericRepo, carNumber string) (text string, err error) {
	var co models.CarOwner
	if err := repo.FindOne(carNumber, &co); err != nil {
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
