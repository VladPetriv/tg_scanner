package channel

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/VladPetriv/tg_scanner/pkg/utils"
	"github.com/gotd/td/tg"
)

type Channel struct {
	ID         int
	ChannelID  int
	Title      string
	AccessHash int
	Username   string
	Photo      tg.ChatPhoto
	Image      *ChannelImage
}

type ChannelImage struct {
	Bytes []byte
}

func GetChannelHistory(ctx context.Context, cPeer *tg.InputPeerChannel, api *tg.Client) (tg.MessagesMessagesClass, error) { // nolint
	bInt := big.NewInt(10000) // nolint

	value, _ := rand.Int(rand.Reader, bInt)

	result, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{ // nolint
		Peer: cPeer,
		Hash: value.Int64(),
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "messages from history", ErrorValue: err}
	}

	return result, nil
}

func GetAllChannels(ctx context.Context, api *tg.Client) ([]Channel, error) {
	channels := make([]Channel, 0)

	data, err := api.MessagesGetAllChats(ctx, []int64{})
	if err != nil {
		return nil, &utils.GettingError{Name: "channels", ErrorValue: err}
	}

	var newChannel Channel

	for _, channel := range data.GetChats() {
		fullChannel, _ := channel.AsFull()
		encodedData, err := json.Marshal(fullChannel)
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

func GetChannelPhoto(ctx context.Context, channel *Channel, api *tg.Client) (tg.UploadFileClass, error) {
	var id int
	if channel.ChannelID == 0 {
		id = channel.ID
	} else {
		id = channel.ChannelID
	}

	data, err := api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPeerPhotoFileLocation{
			Peer: &tg.InputPeerChannel{
				ChannelID:  int64(id),
				AccessHash: int64(channel.AccessHash),
			},
			PhotoID: channel.Photo.PhotoID,
		},
		Limit: 1024 * 1024,
	})
	if err != nil {
		return nil, &utils.GettingError{Name: "channel photo", ErrorValue: err}
	}

	return data, nil
}

func DecodeChannelPhoto(photo tg.UploadFileClass) (*ChannelImage, error) {
	if photo == nil {
		return nil, fmt.Errorf("photo is nil")
	}

	var channelImage *ChannelImage

	js, err := json.Marshal(photo)
	if err != nil {
		return nil, &utils.CreateError{Name: "JSON", ErrorValue: err}
	}

	err = json.Unmarshal(js, &channelImage)
	if err != nil {
		return nil, fmt.Errorf("unmarshal JSON error: %w", err)
	}

	return channelImage, nil
}

func ProcessChannelPhoto(ctx context.Context, channel *Channel, api *tg.Client) (string, error) {
	channelPhotoData, err := GetChannelPhoto(ctx, channel, api)
	if err != nil {
		return "", err
	}

	channelImage, err := DecodeChannelPhoto(channelPhotoData)
	if err != nil {
		return "", fmt.Errorf("decode channel photo error: %w", err)
	}

	channel.Image = channelImage

	fileName, err := CreateChannelImage(channel)
	if err != nil {
		return "", &utils.CreateError{Name: "channel image", ErrorValue: err}
	}

	return fileName, nil
}

func CreateChannelImage(channel *Channel) (string, error) {
	if channel.Image == nil {
		return "", fmt.Errorf("channel image is nil")
	}

	path := fmt.Sprintf("./images/%s.jpg", channel.Username)
	image, err := os.Create(path)
	if err != nil {
		return "", &utils.CreateError{Name: "file", ErrorValue: err}
	}

	_, err = image.Write(channel.Image.Bytes)
	if err != nil {
		return "", fmt.Errorf("write file error: %w", err)
	}

	return path, nil
}
