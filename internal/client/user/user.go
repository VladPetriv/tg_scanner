package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

var userImageSize int = 1024 * 1024

func GetUserInfo(ctx context.Context, userID int64, messageID int, cPeer *tg.InputPeerChannel, api *tg.Client) (*model.TgUser, error) {
	user := &model.TgUser{}

	data, err := api.UsersGetFullUser(ctx, &tg.InputUserFromMessage{
		Peer:   cPeer,
		UserID: userID,
		MsgID:  messageID,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "user from telegram", ErrorValue: err}
	}

	for _, userData := range data.Users {
		notEmptyUser, _ := userData.AsNotEmpty()

		encodedData, err := json.Marshal(notEmptyUser)
		if err != nil {
			return nil, &utils.CreateError{Name: "JSON", ErrorValue: err}
		}

		err = json.Unmarshal(encodedData, &user)
		if err != nil {
			return nil, fmt.Errorf("unmarshal JSON error: %w", err)
		}
	}

	time.Sleep(time.Second * 3)

	return user, nil
}

func GetUserPhoto(ctx context.Context, user *model.TgUser, api *tg.Client) (tg.UploadFileClass, error) {
	var id int64

	if user.ID == 0 {
		id = user.UserID
	} else {
		id = user.ID
	}

	data, err := api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPeerPhotoFileLocation{
			Peer: &tg.InputPeerUser{
				UserID:     id,
				AccessHash: user.AccessHash,
			},
			PhotoID: user.Photo.PhotoID,
		},
		Limit: userImageSize,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "user photo", ErrorValue: err}
	}

	return data, nil
}

func ProcessUserPhoto(ctx context.Context, user *model.TgUser, api *tg.Client) (string, error) {
	userPhotoData, err := GetUserPhoto(ctx, user, api)
	if err != nil {
		return "", err
	}

	userImage, err := file.DecodePhoto(userPhotoData)
	if err != nil {
		return "", fmt.Errorf("decode user photo error: %w", err)
	}

	user.Image = userImage

	fileName, err := file.CreatePhoto(userImage, user.Username)
	if err != nil {
		return "", err
	}

	return fileName, nil
}
