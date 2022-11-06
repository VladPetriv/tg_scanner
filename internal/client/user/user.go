package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
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

func (u tgUser) GetUser(ctx context.Context, data interface{}, groupPeer *tg.InputPeerChannel) (*model.TgUser, error) {
	logger := u.log

	userID, modelID := getUserDataFromEntity(data)
	user := model.TgUser{}

	fullUser, err := u.api.UsersGetFullUser(ctx, &tg.InputUserFromMessage{
		Peer:   groupPeer,
		UserID: userID,
		MsgID:  modelID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("get full user")
		return nil, fmt.Errorf("get full user error: %w", err)
	}

	for _, userData := range fullUser.Users {
		notEmptyUser, _ := userData.AsNotEmpty()

		encodedData, err := json.Marshal(notEmptyUser)
		if err != nil {
			logger.Error().Err(err).Msg("marshal user data")
			return nil, fmt.Errorf("marshal user data error: %w", err)
		}

		err = json.Unmarshal(encodedData, &user)
		if err != nil {
			logger.Error().Err(err).Msg("unmarshal user data")
			return nil, fmt.Errorf("unmarshal user data error: %w", err)
		}
	}

	// here we sleep to avoid timeout error from telegram API
	time.Sleep(_getUserInfoTimeout)

	return &user, nil
}

func getUserDataFromEntity(data interface{}) (int64, int) {
	var userID int64
	var modelID int

	switch dataType := data.(type) {
	case model.TgMessage:
		if dataType.FromID.ID != 0 {
			userID = dataType.FromID.ID
		} else {
			userID = dataType.FromID.UserID
		}

		modelID = dataType.ID
	case model.TgRepliesMessage:
		if dataType.FromID.ID != 0 {
			userID = dataType.FromID.ID
		} else {
			userID = dataType.FromID.UserID
		}

		modelID = dataType.ID
	}

	return userID, modelID
}

func (u tgUser) GetUserPhoto(ctx context.Context, user model.TgUser) (tg.UploadFileClass, error) {
	logger := u.log

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
		logger.Error().Err(err).Msg("get user photo")
		return nil, fmt.Errorf("get user photo error: %w", err)
	}

	return data, nil
}
