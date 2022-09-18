package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/errors"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

var _getUserInfoTimeout = 3 * time.Second

type tgUser struct {
	log *logger.Logger
	api *tg.Client
}

var _ User = (*tgUser)(nil)

func New(log *logger.Logger, api *tg.Client) *tgUser {
	return &tgUser{
		log: log,
		api: api,
	}
}

func (u tgUser) GetUser(ctx context.Context, message *model.TgMessage, groupPeer *tg.InputPeerChannel) (*model.TgUser, error) {
	user := model.TgUser{}

	data, err := u.api.UsersGetFullUser(ctx, &tg.InputUserFromMessage{
		Peer:   groupPeer,
		UserID: message.FromID.ID,
		MsgID:  message.ID,
	})
	if err != nil {
		u.log.Error().Err(err)

		return nil, &errors.GetError{Name: "user info", ErrorValue: err}
	}

	for _, userData := range data.Users {
		notEmptyUser, _ := userData.AsNotEmpty()

		encodedData, err := json.Marshal(notEmptyUser)
		if err != nil {
			u.log.Warn().Err(err)

			return nil, &errors.CreateError{Name: "JSON", ErrorValue: err}
		}

		err = json.Unmarshal(encodedData, &user)
		if err != nil {
			u.log.Warn().Err(err)

			return nil, fmt.Errorf("unmarshal JSON error: %w", err)
		}
	}

	time.Sleep(_getUserInfoTimeout)

	return &user, nil
}

func (u tgUser) GetUserPhoto(ctx context.Context, user model.TgUser) (tg.UploadFileClass, error) {
	data, err := u.api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPeerPhotoFileLocation{
			Peer: &tg.InputPeerUser{
				UserID:     user.ID,
				AccessHash: user.AccessHash,
			},
			PhotoID: user.Photo.PhotoID,
		},
		Limit: photo.Size,
	})
	if err != nil {
		u.log.Error().Err(err)

		return nil, &errors.GetError{Name: "user photo", ErrorValue: err}
	}

	return data, nil
}
