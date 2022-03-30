package channel

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/gotd/td/tg"
	"github.com/sirupsen/logrus"
)

type Group struct {
	ID         int
	Title      string
	AccessHash int
	Username   string
}

func GetChannelHistory(ctx context.Context, limit int, cPeer tg.InputPeerChannel, api *tg.Client) (tg.MessagesMessagesClass, error) { // nolint
	bInt := big.NewInt(10000) // nolint

	value, err := rand.Int(rand.Reader, bInt)
	if err != nil {
		logrus.Errorf("ERROR_WHILE_GENERATE_RANDOM_INT:%s", err)
	}

	result, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{ // nolint
		Peer: &cPeer,
		Hash: value.Int64(),
	})
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_GETTING_MESSAGES:%w", err)
	}

	return result, nil
}

func GetAllGroups(ctx context.Context, api *tg.Client) ([]Group, error) {
	groups := make([]Group, 0)

	data, err := api.MessagesGetAllChats(ctx, []int64{})
	if err != nil {
		return nil, fmt.Errorf("ERROR_WHILE_GETTING_GROUPS:%w", err)
	}

	var newGroup Group

	for _, group := range data.GetChats() {
		encodedData, err := json.Marshal(group)
		if err != nil {
			continue
		}

		err = json.Unmarshal(encodedData, &newGroup)
		if err != nil {
			continue
		}

		groups = append(groups, newGroup)
	}

	return groups, nil
}
