package user

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gotd/td/tg"
)

type User struct {
	ID        int
	FirstName string
	LastName  string
	Username  string
	Photo     interface{}
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

	return user, nil
}
