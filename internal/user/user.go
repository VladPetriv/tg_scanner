package user

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/VladPetriv/tg_scanner/pkg/utils"
	"github.com/gotd/td/tg"
)

type User struct {
	ID         int
	UserID     int
	FirstName  string
	LastName   string
	Username   string
	Photo      tg.UserProfilePhoto
	AccessHash int
	Image      *UserImage
}

type UserImage struct {
	Bytes []byte
}

func GetUserInfo(ctx context.Context, userID int, messageID int, cPeer *tg.InputPeerChannel, api *tg.Client) (*User, error) {
	var user *User

	data, err := api.UsersGetFullUser(ctx, &tg.InputUserFromMessage{
		Peer:   cPeer,
		UserID: int64(userID),
		MsgID:  messageID,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "user from telegram", ErrorValue: err}
	}

	for _, u := range data.Users {
		notEmptyUser, _ := u.AsNotEmpty()

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

func GetUserPhoto(ctx context.Context, user *User, api *tg.Client) (tg.UploadFileClass, error) {
	var id int
	if user.ID == 0 {
		id = user.UserID
	}

	id = user.ID

	data, err := api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPeerPhotoFileLocation{
			Peer: &tg.InputPeerUser{
				UserID:     int64(id),
				AccessHash: int64(user.AccessHash),
			},
			PhotoID: user.Photo.PhotoID,
		},
		Limit: 1024 * 1024,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "user photo", ErrorValue: err}
	}

	return data, nil
}

func DecodeUserPhoto(photo tg.UploadFileClass) (*UserImage, error) {
	if photo == nil {
		return nil, fmt.Errorf("photo is nil")
	}

	var userImage *UserImage

	js, err := json.Marshal(photo)
	if err != nil {
		return nil, &utils.CreateError{Name: "JSON", ErrorValue: err}
	}

	err = json.Unmarshal(js, &userImage)
	if err != nil {
		return nil, fmt.Errorf("unmarshal JSON error: %w", err)
	}

	return userImage, nil
}

func ProcessUserPhoto(ctx context.Context, user *User, api *tg.Client) (string, error) {
	userPhotoData, err := GetUserPhoto(ctx, user, api)
	if err != nil {
		return "", err
	}

	userImage, err := DecodeUserPhoto(userPhotoData)
	if err != nil {
		return "", fmt.Errorf("decode user photo error: %w", err)
	}

	user.Image = userImage

	fileName, err := CreateUserPhoto(user)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func CreateUserPhoto(user *User) (string, error) {
	if user.Image == nil {
		return "", fmt.Errorf("user image is nil")
	}

	path := fmt.Sprintf("./images/%s.jpg", user.Username)
	photo, err := os.Create(path)
	if err != nil {
		return "", &utils.CreateError{Name: "user photo", ErrorValue: err}
	}

	_, err = photo.Write(user.Image.Bytes)
	if err != nil {
		return "", fmt.Errorf("write file error: %w", err)
	}

	return path, nil
}
