package channel

import (
	"context"
	"encoding/json"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

func GetChannelHistory(ctx context.Context, channelPeer *tg.InputPeerChannel, api *tg.Client) (tg.MessagesMessagesClass, error) { // nolint
	result, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{ // nolint
		Peer: channelPeer,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "messages from history", ErrorValue: err}
	}

	return result, nil
}

func GetAllChannels(ctx context.Context, api *tg.Client) ([]model.TgChannel, error) {
	channels := make([]model.TgChannel, 0)

	data, err := api.MessagesGetAllChats(ctx, []int64{})
	if err != nil {
		return nil, &utils.GettingError{Name: "channels", ErrorValue: err}
	}

	for _, channelData := range data.GetChats() {
		channel := model.TgChannel{}

		fullChannel, _ := channelData.AsFull()
		encodedData, err := json.Marshal(fullChannel)
		if err != nil {
			continue
		}

		err = json.Unmarshal(encodedData, &channel)
		if err != nil {
			continue
		}

		channels = append(channels, channel)
	}

	return channels, nil
}

func GetChannelPhoto(ctx context.Context, channel *model.TgChannel, api *tg.Client) (tg.UploadFileClass, error) {
	data, err := api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPeerPhotoFileLocation{
			Peer: &tg.InputPeerChannel{
				ChannelID:  channel.ID,
				AccessHash: channel.AccessHash,
			},
			PhotoID: channel.Photo.PhotoID,
		},
		Limit: photo.Size,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "channel photo", ErrorValue: err}
	}

	return data, nil
}
