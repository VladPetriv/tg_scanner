package group

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

type tgGroup struct {
	log *logger.Logger
	api *tg.Client
}

var _ TgGroup = (*tgGroup)(nil)

func New(log *logger.Logger, api *tg.Client) TgGroup {
	return &tgGroup{
		log: log,
		api: api,
	}
}

func (g tgGroup) GetGroups(ctx context.Context) ([]model.TgGroup, error) {
	logger := g.log

	groups := make([]model.TgGroup, 0)

	data, err := g.api.MessagesGetAllChats(ctx, []int64{})
	if err != nil {
		logger.Error().Err(err).Msg("get all groups")
		return nil, fmt.Errorf("get all groups: %w", err)
	}

	for _, groupData := range data.GetChats() {
		var group model.TgGroup

		fullGroupInfo, isFull := groupData.AsFull()
		if !isFull {
			logger.Warn().Bool("is full", isFull).Msg("received unexpected type of group")

			continue
		}

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

	logger.Info().Interface("groups", groups).Msg("successfully got telegram groups")
	return groups, nil
}

func (g tgGroup) GetGroupPhoto(ctx context.Context, group model.TgGroup) (tg.UploadFileClass, error) {
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

func (g tgGroup) CreateFilesForGroups(groups []model.TgGroup) error {
	logger := g.log

	for _, group := range groups {
		fileName := fmt.Sprintf("%s.json", group.Username)
		if _, err := os.Stat("./data/" + fileName); err == nil {
			continue
		}

		file, err := os.Create(fileName)
		if err != nil {
			logger.Error().Err(err).Msg("create file")
			return fmt.Errorf("create file error: %w", err)
		}

		_, err = file.WriteString("[ ]")
		if err != nil {
			logger.Error().Err(err).Msg("write to file")
			return fmt.Errorf("write to file error: %w", err)
		}

		err = os.Rename(fileName, fmt.Sprintf("./data/%s", fileName))
		if err != nil {
			logger.Error().Err(err).Msg("rename file")
			return fmt.Errorf("rename file error: %w", err)
		}
	}

	return nil
}
