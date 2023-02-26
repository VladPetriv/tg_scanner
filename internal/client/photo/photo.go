package photo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

const Size = 1024 * 1024

type tgPhoto struct {
	log   *logger.Logger
	store *store.Store
}

var _ TgPhoto = (*tgPhoto)(nil)

func New(log *logger.Logger, store *store.Store) TgPhoto {
	return &tgPhoto{
		log:   log,
		store: store,
	}
}

func (p tgPhoto) ProcessPhoto(ctx context.Context, photoData tg.UploadFileClass, name string) (string, error) {
	logger := p.log

	image, err := p.decodePhoto(photoData)
	if err != nil {
		logger.Warn().Err(err).Msg("decode photo")
	}

	filename, err := p.createPhoto(image, fmt.Sprint(name))
	if err != nil {
		logger.Warn().Err(err).Msg("create photo file")
	}

	imageURL, err := p.store.Image.Send(ctx, filename, fmt.Sprint(name))
	if err != nil {
		logger.Error().Err(err).Msg("send image into firebase")
		return imageURL, fmt.Errorf("send image into firebase error: %w", err)
	}

	return imageURL, nil
}

func (p tgPhoto) decodePhoto(photo tg.UploadFileClass) (*model.Image, error) {
	logger := p.log

	var img model.Image

	encodedData, err := json.Marshal(photo)
	if err != nil {
		logger.Error().Err(err).Msg("marshal photo data")
		return nil, fmt.Errorf("marshal photo data error: %w", err)
	}

	err = json.Unmarshal(encodedData, &img)
	if err != nil {
		logger.Error().Err(err).Msg("unmarshal photo data")
		return nil, fmt.Errorf("unmarshal photo data error: %w", err)
	}

	return &img, nil
}

func (p tgPhoto) createPhoto(img *model.Image, name string) (string, error) {
	logger := p.log

	path := fmt.Sprintf("./images/%s.jpg", name)
	photo, err := os.Create(path)
	if err != nil {
		logger.Error().Err(err).Msg("create file")
		return "", fmt.Errorf("create file error: %w", err)
	}

	_, err = photo.Write(img.Bytes)
	if err != nil {
		logger.Error().Err(err).Msg("wirte data to file")
		return "", fmt.Errorf("write data to file error: %w", err)
	}

	return path, nil
}
