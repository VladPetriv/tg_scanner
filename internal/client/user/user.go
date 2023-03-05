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

const getUserInfoTimeout = 3 * time.Second

type tgUser struct {
	log *logger.Logger
	api *tg.Client
}

var _ User = (*tgUser)(nil)

func New(log *logger.Logger, api *tg.Client) User {
	return &tgUser{
		log: log,
		api: api,
	}
}

func (u tgUser) GetUser(ctx context.Context, entity interface{}, group *model.Group) (*model.User, error) {
	logger := u.log

	userID, entityID := getUserDataFromEntity(entity)

	data, err := u.api.UsersGetFullUser(ctx, &tg.InputUserFromMessage{
		Peer: &tg.InputPeerChannel{
			ChannelID:  group.ID,
			AccessHash: group.AccessHash,
		},
		UserID: userID,
		MsgID:  entityID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("get user data")
		return nil, fmt.Errorf("get user data: %w", err)
	}

	user := u.parseUser(data.Users)

	time.Sleep(getUserInfoTimeout)

	return &user, nil
}

func (u tgUser) parseUser(tgUsers []tg.UserClass) model.User {
	logger := u.log

	var user model.User

	for _, usr := range tgUsers {
		notEmptyUser, ok := usr.AsNotEmpty()
		if !ok {
			logger.Info().Msg("received empty user")

			continue
		}

		encodedData, err := json.Marshal(notEmptyUser)
		if err != nil {
			logger.Error().Err(err).Msg("marshal user data")
		}

		err = json.Unmarshal(encodedData, &user)
		if err != nil {
			logger.Error().Err(err).Msg("unmarshal user data")
		}
	}

	return user
}

func getUserDataFromEntity(data interface{}) (int64, int) {
	switch dataType := data.(type) {
	case model.Message:
		if dataType.FromID.ID != 0 {
			return dataType.FromID.ID, dataType.ID
		}

		return dataType.FromID.UserID, dataType.ID
	case model.RepliesMessage:
		if dataType.FromID.ID != 0 {
			return dataType.FromID.ID, dataType.ID
		}

		return dataType.FromID.UserID, dataType.ID
	default:
		return 0, 0
	}
}

func (u tgUser) GetUserPhoto(ctx context.Context, user model.User) (tg.UploadFileClass, error) {
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
