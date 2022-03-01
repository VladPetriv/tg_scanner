package channel

import (
	"context"
	"strconv"
	"strings"

	"github.com/gotd/td/tg"
)

func GetAccessHash(ctx context.Context, groupName string, api *tg.Client) (int64, error) {
	var accessHash int
	group, err := api.ContactsResolveUsername(ctx, groupName)
	if err != nil {
		return 0, err
	}
	for _, chat := range group.Chats {
		accessHash, _ = strconv.Atoi(strings.Split(chat.String(), " ")[19][11:])
	}

	return int64(accessHash), nil
}

func GetChannelHistory(ctx context.Context, api *tg.Client, channelPeer tg.InputPeerChannel, limint int) (tg.MessagesMessagesClass, error) {
	result, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer:  &channelPeer,
		Hash:  31243312413321,
		Limit: limint,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
