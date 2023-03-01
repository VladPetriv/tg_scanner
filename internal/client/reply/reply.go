package reply

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
	"github.com/gotd/td/tg"
)

type tgReply struct {
	log *logger.Logger
	api *tg.Client
}

var _ TgReply = (*tgReply)(nil)

func New(log *logger.Logger, api *tg.Client) TgReply {
	return &tgReply{
		log: log,
		api: api,
	}
}

func (r tgReply) GetReplies(ctx context.Context, msg model.TgMessage) (tg.MessagesMessagesClass, error) {
	logger := r.log

	replies, err := r.api.MessagesGetReplies(ctx, &tg.MessagesGetRepliesRequest{
		Peer: &tg.InputPeerChannel{
			ChannelID:  msg.PeerID.ID,
			AccessHash: msg.PeerID.AccessHash,
		},
		MsgID: msg.ID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("get replies from message")
		return nil, fmt.Errorf("get replies from message error: %w", err)
	}

	return replies, nil
}

func (r tgReply) ParseTelegramReplies(ctx context.Context, replies tg.MessagesMessagesClass, groupPeer *tg.InputPeerChannel) []model.TgRepliesMessage { //nolint:lll
	logger := r.log

	parsedReplies := make([]model.TgRepliesMessage, 0)
	repliesMessages, _ := replies.AsModified()

	for _, rpl := range repliesMessages.GetMessages() {
		reply := model.TgRepliesMessage{}

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

func (r tgReply) GetReplyPhoto(ctx context.Context, reply model.TgRepliesMessage) (tg.UploadFileClass, error) {
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
