package reply

import (
	"context"
	"encoding/json"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/client/photo"
	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/pkg/errors"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

type tgReply struct {
	log *logger.Logger
	api *tg.Client
}

var _ Reply = (*tgReply)(nil)

func New(log *logger.Logger, api *tg.Client) *tgReply {
	return &tgReply{
		log: log,
		api: api,
	}
}

func (r tgReply) GetReplies(ctx context.Context, message *model.TgMessage, groupPeer *tg.InputPeerChannel) (tg.MessagesMessagesClass, error) {
	replies, err := r.api.MessagesGetReplies(ctx, &tg.MessagesGetRepliesRequest{
		Peer:  groupPeer,
		MsgID: message.ID,
	})
	if err != nil {
		r.log.Error().Err(err)

		return nil, &errors.GetError{Name: "replies", ErrorValue: err}
	}

	return replies, nil
}

func (r tgReply) ProcessReplies(ctx context.Context, replies tg.MessagesMessagesClass, groupPeer *tg.InputPeerChannel) []model.TgRepliesMessage {
	processedReplies := make([]model.TgRepliesMessage, 0)

	modifiedReplies, _ := replies.AsModified()

	for _, rpl := range modifiedReplies.GetMessages() {
		reply := model.TgRepliesMessage{}

		encodedData, err := json.Marshal(rpl)
		if err != nil {
			r.log.Warn().Err(err)

			continue
		}

		err = json.Unmarshal(encodedData, &reply)
		if err != nil {
			r.log.Warn().Err(err)

			continue
		}

		processedReplies = append(processedReplies, reply)
	}

	return processedReplies
}

func (r tgReply) GetRepliePhoto(ctx context.Context, reply model.TgRepliesMessage) (tg.UploadFileClass, error) {
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
		return nil, &errors.GetError{Name: "reply photo", ErrorValue: err}
	}

	return data, nil
}
