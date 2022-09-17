package group

import (
	"context"
	"encoding/json"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/errors"
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
	groups := make([]model.TgGroup, 5)

	data, err := g.api.MessagesGetAllChats(ctx, []int64{})
	if err != nil {
		g.log.Error().Err(err)
		return nil, &errors.GetError{Name: "groups", ErrorValue: err}
	}

	for _, groupData := range data.GetChats() {
		var group model.TgGroup

		fullGroupInfo, _ := groupData.AsFull()

		encodedData, err := json.Marshal(fullGroupInfo)
		if err != nil {
			g.log.Warn().Err(err)
			continue
		}

		err = json.Unmarshal(encodedData, &group)
		if err != nil {
			g.log.Warn().Err(err)
			continue
		}

		groups = append(groups, group)
	}

	return groups, nil
}

func (g tgGroup) GetMessagesFromGroupHistory(ctx context.Context, groupPeer *tg.InputPeerChannel) (tg.MessagesMessagesClass, error) {
	groupHistory, err := g.api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
		Peer: groupPeer,
	})
	if err != nil {
		g.log.Error().Err(err)
		return nil, &errors.GetError{Name: "messages from history", ErrorValue: err}
	}

	return groupHistory, nil
}

func (g tgGroup) GetGroupPhoto(ctx context.Context, group *model.TgGroup) (tg.UploadFileClass, error) {
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
		g.log.Error().Err(err)
		return nil, &errors.GetError{Name: "group photo", ErrorValue: err}
	}

	return data, nil
}
