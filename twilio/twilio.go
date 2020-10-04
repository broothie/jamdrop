package twilio

import (
	"fmt"
	"log"

	"jamdrop/config"
	"jamdrop/model"

	"github.com/kevinburke/twilio-go"
	"github.com/pkg/errors"
)

const fromNumber = "+19382014947"

type Twilio struct {
	*twilio.Client
	Logger *log.Logger
}

func New(cfg *config.Config, logger *log.Logger) *Twilio {
	return &Twilio{
		Client: twilio.NewClient(cfg.TwilioAccountSID, cfg.TwilioAuthToken, nil),
		Logger: logger,
	}
}

func (t *Twilio) SongQueued(user *model.User, songName string) error {
	t.Logger.Println("twilio.SongQueued", user.ID, songName)

	if user.PhoneNumber == "" {
		t.Logger.Println("user does not have a phone number; user_id", user.ID)
		return nil
	}

	body := fmt.Sprintf(`JamDrop: %s dropped "%s" into your queue`, user.DisplayName, songName)
	if _, err := t.Messages.SendMessage(fromNumber, user.PhoneNumber, body, nil); err != nil {
		return errors.Wrapf(err, "failed to send song queued message; user_id: %s, song_name: %s", user.ID, songName)
	}

	return nil
}
