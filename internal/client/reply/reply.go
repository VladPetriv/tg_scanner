package reply

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

type tgReply struct {
	log *logger.Logger
	api *tg.Client
}

var _ Reply = (*tgReply)(nil)

func New(log *logger.Logger, api *tg.Client) Reply {
	return &tgReply{
		log: log,
		api: api,
	}
}

func (r tgReply) GetReplies(ctx context.Context, message model.Message) ([]model.RepliesMessage, error) {
	logger := r.log

	data, err := r.api.MessagesGetReplies(ctx, &tg.MessagesGetRepliesRequest{
		Peer: &tg.InputPeerChannel{
			ChannelID:  message.PeerID.ID,
			AccessHash: message.PeerID.AccessHash,
		},
		MsgID: message.ID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("get replies for message")
		return nil, fmt.Errorf("get replies for message: %w", err)
	}

	modifiedData, ok := data.AsModified()
	if !ok {
		logger.Info().Msg("received unexpected type of replies")
	}

	tgReplies := modifiedData.GetMessages()

	replies := r.parseTelegramReplies(tgReplies)

	return replies, nil
}

func (r tgReply) parseTelegramReplies(tgReplies []tg.MessageClass) []model.RepliesMessage {
	logger := r.log

	parsedReplies := make([]model.RepliesMessage, 0)

	for _, rpl := range tgReplies {
		var reply model.RepliesMessage

		encodedData, err := json.Marshal(rpl)
		if err != nil {
			logger.Warn().Err(err).Msg("marshal reply data")

			continue
		}

		err = json.Unmarshal(encodedData, &reply)
		if err != nil {
			r.log.Warn().Err(err).Msg("unmarshal reply data")

			continue
		}

		parsedReplies = append(parsedReplies, reply)
	}

	return parsedReplies
}

func (r tgReply) GetReplyPhoto(ctx context.Context, reply model.RepliesMessage) (tg.UploadFileClass, error) {
	logger := r.log

	length := len(reply.Media.Photo.Sizes) - 1

	data, err := r.api.UploadGetFile(ctx, &tg.UploadGetFileRequest{
		Location: &tg.InputPhotoFileLocation{
			ID:            reply.Media.Photo.ID,
			AccessHash:    reply.Media.Photo.AccessHash,
			FileReference: reply.Media.Photo.FileReference,
			ThumbSize:     reply.Media.Photo.Sizes[length].GetType(),
		},
		Offset: 0,
		Limit:  photo.Size,
	})
	if err != nil {
		logger.Error().Err(err).Msg("get reply photo")
		return nil, fmt.Errorf("get reply photo error: %w", err)
	}

	return data, nil
}
