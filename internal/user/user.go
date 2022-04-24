package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	Image      UserImage
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
		return nil, fmt.Errorf("getting user error: %w", err)
	}

	for _, u := range data.Users {
		notEmptyUser, _ := u.AsNotEmpty()

		encodedData, err := json.Marshal(notEmptyUser)
		if err != nil {
			return nil, fmt.Errorf("creating JSON error: %w", err)
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
		return nil, fmt.Errorf("getting user photo error: %w", err)
	}

	return data, nil
}

func DecodeUserPhoto(photo tg.UploadFileClass) (*UserImage, error) {
	var userImage *UserImage

	js, err := json.Marshal(photo)
	if err != nil {
		return nil, fmt.Errorf("createing JSON error: %w", err)
	}
	err = json.Unmarshal(js, &userImage)
	if err != nil {
		return nil, fmt.Errorf("unmarshal JSON error: %w", err)
	}

	return userImage, nil
}
