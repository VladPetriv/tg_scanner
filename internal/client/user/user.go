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
)

func GetUserInfo(ctx context.Context, userID int64, messageID int, cPeer *tg.InputPeerChannel, api *tg.Client) (*model.TgUser, error) {
	user := &model.TgUser{}

	data, err := api.UsersGetFullUser(ctx, &tg.InputUserFromMessage{
		Peer:   cPeer,
		UserID: userID,
		MsgID:  messageID,
	})
	if err != nil {
		return nil, &errors.GettingError{Name: "user from telegram", ErrorValue: err}
	}

	for _, userData := range data.Users {
		notEmptyUser, _ := userData.AsNotEmpty()

		encodedData, err := json.Marshal(notEmptyUser)
		if err != nil {
			return nil, &errors.CreateError{Name: "JSON", ErrorValue: err}
		}

		err = json.Unmarshal(encodedData, &user)
		if err != nil {
			return nil, fmt.Errorf("unmarshal JSON error: %w", err)
		}
	}

	time.Sleep(time.Second * 3)

	return user, nil
}

func GetUserPhoto(ctx context.Context, user model.TgUser, api *tg.Client) (tg.UploadFileClass, error) {
	data, err := api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
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
		return nil, &errors.GettingError{Name: "user photo", ErrorValue: err}
	}

	return data, nil
}
