package commands

import (
	"github.com/triviy/parklakes-viberbot/application/integrations/viber"
	"github.com/triviy/parklakes-viberbot/domain/interfaces"
)

// WelcomeCmd instance of viber webhook cmd
type WelcomeCmd struct {
	subscriberRepo interfaces.GenericRepo
}

// NewWelcomeCmd creates new instance of WelcomeCmd
func NewWelcomeCmd() *WelcomeCmd {
	return &WelcomeCmd{}
}

// Execute calls setting Viber callback URLs
func (cmd WelcomeCmd) Execute() *viber.WelcomeResponse {
	return &viber.WelcomeResponse{
		Type: viber.TextType,
		Text: "Вітаю! Введіть повний номер автівки, а я спробую вам допомогти",
	}
}
