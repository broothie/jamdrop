package twilio

import (
	"context"
	"fmt"
	"jamdrop/config"
	"jamdrop/logger"
	"jamdrop/model"
	"jamdrop/requestid"

	"github.com/kevinburke/twilio-go"
	"github.com/pkg/errors"
)

const fromNumber = "+19382014947"

type Twilio struct {
	*twilio.Client
	Logger *logger.Logger
}

func New(cfg *config.Config, logger *logger.Logger) *Twilio {
	return &Twilio{
		Client: twilio.NewClient(cfg.TwilioAccountSID, cfg.TwilioAuthToken, nil),
		Logger: logger,
	}
}

func (t *Twilio) SongQueued(ctx context.Context, user, targetUser *model.User, songName string) error {
	t.Logger.Debug("twilio.SongQueued", logger.Fields{"user_id": user.ID, "target_user_id": targetUser.ID, "song_name": songName}, requestid.LogContext(ctx))

	if targetUser.PhoneNumber == "" {
		t.Logger.Info("user does not have a phone number", logger.Fields{"target_user_id": targetUser.ID})
		return nil
	}

	body := fmt.Sprintf(`JamDrop: %s dropped "%s" into your queue`, user.DisplayName, songName)
	if _, err := t.Messages.SendMessage(fromNumber, targetUser.PhoneNumber, body, nil); err != nil {
		return errors.Wrapf(err, "failed to send song queued message; user_id: %s, song_name: %s", user.ID, songName)
	}

	return nil
}
