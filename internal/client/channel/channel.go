package channel

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/file"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

var channelImageSize int = 1024 * 1024

func GetChannelHistory(ctx context.Context, cPeer *tg.InputPeerChannel, api *tg.Client) (tg.MessagesMessagesClass, error) { // nolint
	result, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{ // nolint
		Peer: cPeer,
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
	var id int64
	if channel.ChannelID == 0 {
		id = channel.ID
	} else {
		id = channel.ChannelID
	}

	data, err := api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPeerPhotoFileLocation{
			Peer: &tg.InputPeerChannel{
				ChannelID:  id,
				AccessHash: channel.AccessHash,
			},
			PhotoID: channel.Photo.PhotoID,
		},
		Limit: channelImageSize,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "channel photo", ErrorValue: err}
	}

	return data, nil
}

func ProcessChannelPhoto(ctx context.Context, channel *model.TgChannel, api *tg.Client) (string, error) {
	channelPhotoData, err := GetChannelPhoto(ctx, channel, api)
	if err != nil {
		return "", err
	}

	channelImage, err := file.DecodePhoto(channelPhotoData)
	if err != nil {
		return "", fmt.Errorf("decode channel photo error: %w", err)
	}

	channel.Image = channelImage

	fileName, err := file.CreatePhoto(channelImage, channel.Username)
	if err != nil {
		return "", &utils.CreateError{Name: "channel image", ErrorValue: err}
	}

	return fileName, nil
}
