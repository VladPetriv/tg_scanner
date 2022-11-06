package group

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

type tgGroup struct {
	log *logger.Logger
	api *tg.Client
}

var _ Group = (*tgGroup)(nil)

func New(log *logger.Logger, api *tg.Client) *tgGroup {
	return &tgGroup{
		log: log,
		api: api,
	}
}

func (g tgGroup) GetGroups(ctx context.Context) ([]model.TgGroup, error) {
	logger := g.log

	groups := make([]model.TgGroup, 5)

	data, err := g.api.MessagesGetAllChats(ctx, []int64{})
	if err != nil {
		logger.Error().Err(err).Msg("get all groups")
		return nil, fmt.Errorf("get all groups error: %w", err)
	}

	for _, groupData := range data.GetChats() {
		var group model.TgGroup

		fullGroupInfo, _ := groupData.AsFull()

		encodedData, err := json.Marshal(fullGroupInfo)
		if err != nil {
			logger.Warn().Err(err).Msg("marshal group data")

			continue
		}

		err = json.Unmarshal(encodedData, &group)
		if err != nil {
			logger.Warn().Err(err).Msg("Unmarshal group data")

			continue
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (g tgGroup) GetMessagesFromGroupHistory(ctx context.Context, groupPeer *tg.InputPeerChannel) (tg.MessagesMessagesClass, error) {
	logger := g.log

	groupHistory, err := g.api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer: groupPeer,
	})
	if err != nil {
		logger.Error().Err(err).Msg("get messages from group history")
		return nil, fmt.Errorf("get messages from group history error: %w", err)
	}

	return groupHistory, nil
}

func (g tgGroup) GetGroupPhoto(ctx context.Context, group *model.TgGroup) (tg.UploadFileClass, error) {
	logger := g.log

	data, err := g.api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPeerPhotoFileLocation{
			Peer: &tg.InputPeerChannel{
				ChannelID:  group.ID,
				AccessHash: group.AccessHash,
			},
			PhotoID: group.Photo.PhotoID,
		},
		Limit: photo.Size,
	})
	if err != nil {
		logger.Error().Err(err).Msg("get group photo")
		return nil, fmt.Errorf("get group photo error: %w", err)
	}

	return data, nil
}
