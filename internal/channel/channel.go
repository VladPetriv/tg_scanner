package channel

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gotd/td/tg"
)

type Group struct {
	ID         int
	Title      string
	AccessHash int
	Username   string
}

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

func GetAllGroups(ctx context.Context, api *tg.Client) ([]Group, error) {
	var groups []Group
	data, err := api.MessagesGetAllChats(ctx, []int64{})
	if err != nil {
		return nil, err
	}
	var g Group
	for _, group := range data.GetChats() {
		encodedData, err := json.Marshal(group)
		if err != nil {
			continue
		}
		err = json.Unmarshal(encodedData, &g)
		if err != nil {
			continue
		}
		groups = append(groups, g)
	}

	return groups, nil
}
