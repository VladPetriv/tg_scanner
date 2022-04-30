package channel

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/gotd/td/tg"
)

type Channel struct {
	ID         int
	ChannelID  int
	Title      string
	AccessHash int
	Username   string
}

func GetChannelHistory(ctx context.Context, cPeer *tg.InputPeerChannel, api *tg.Client) (tg.MessagesMessagesClass, error) { // nolint
	bInt := big.NewInt(10000) // nolint

	value, _ := rand.Int(rand.Reader, bInt)

	result, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{ // nolint
		Peer: cPeer,
		Hash: value.Int64(),
	})
	if err != nil {
		return nil, fmt.Errorf("getting message from history error: %w", err)
	}

	return result, nil
}

func GetAllChannels(ctx context.Context, api *tg.Client) ([]Channel, error) {
	channels := make([]Channel, 0)

	data, err := api.MessagesGetAllChats(ctx, []int64{})
	if err != nil {
		return nil, fmt.Errorf("getting group error: %w", err)
	}

	var newChannel Channel

	for _, channel := range data.GetChats() {
		encodedData, err := json.Marshal(channel)
		if err != nil {
			continue
		}

		err = json.Unmarshal(encodedData, &newChannel)
		if err != nil {
			continue
		}

		channels = append(channels, newChannel)
	}

	return channels, nil
}
