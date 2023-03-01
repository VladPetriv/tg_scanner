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

func (r tgReply) GetReplies(ctx context.Context, msg model.TgMessage) ([]model.TgRepliesMessage, error) {
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

	parsedResplies := r.parseTelegramReplies(ctx, replies)

	return parsedResplies, nil
}

func (r tgReply) parseTelegramReplies(ctx context.Context, replies tg.MessagesMessagesClass) []model.TgRepliesMessage {
	logger := r.log

	parsedReplies := make([]model.TgRepliesMessage, 0)

	repliesMessages, ok := replies.AsModified()
	if !ok {
		logger.Error().Bool("ok", ok).Msg("received unexpected type of replies")
	}

	for _, reply := range repliesMessages.GetMessages() {
		var replyMessage model.TgRepliesMessage

		encodedData, err := json.Marshal(reply)
		if err != nil {
			logger.Error().Err(err).Msg("marshal reply data")

			continue
		}

		err = json.Unmarshal(encodedData, &replyMessage)
		if err != nil {
			r.log.Error().Err(err).Msg("unmarshal reply data")

			continue
		}

		parsedReplies = append(parsedReplies, replyMessage)
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
