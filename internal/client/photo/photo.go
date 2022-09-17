package photo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store"
	cerrors "github.com/VladPetriv/tg_scanner/pkg/errors"
	"github.com/VladPetriv/tg_scanner/pkg/logger"
)

var Size int = 1024 * 1024
var errPhotoDataIsNil error = errors.New("photo data is nil")

type tgPhoto struct {
	log   *logger.Logger
	store *store.Store
}

var _ Photo = (*tgPhoto)(nil)

func New(log *logger.Logger) *tgPhoto {
	return &tgPhoto{
		log: log,
	}
}

func (p tgPhoto) DecodePhoto(photo tg.UploadFileClass) (*model.Image, error) {
	if photo == nil {
		p.log.Error().Err(errPhotoDataIsNil)

		return nil, errPhotoDataIsNil
	}

	img := model.Image{}

	encodedData, err := json.Marshal(photo)
	if err != nil {
		return nil, &cerrors.CreateError{Name: "JSON", ErrorValue: err}
	}

	err = json.Unmarshal(encodedData, &img)
	if err != nil {
		return nil, fmt.Errorf("unmarshal JSON error: %w", err)
	}

	return &img, nil
}

func CreatePhoto(img *model.Image, name string) (string, error) {
	if img == nil {
		return "", errPhotoDataIsNil
	}

	path := fmt.Sprintf("./images/%s.jpg", name)
	photo, err := os.Create(path)
	if err != nil {
		return "", &cerrors.CreateError{Name: "photo", ErrorValue: err}
	}

	_, err = photo.Write(img.Bytes)
	if err != nil {
		return "", fmt.Errorf("write file error: %w", err)
	}

	return path, nil
}

func (p tgPhoto) ProcessPhoto(ctx context.Context, photoData tg.UploadFileClass, name string) (string, error) {
	image, err := p.DecodePhoto(photoData)
	if err != nil {
		p.log.Warn().Err(err)
	}

	filename, err := CreatePhoto(image, fmt.Sprint(name))
	if err != nil {
		p.log.Warn().Err(err)
	}

	imageUrl, err := p.store.Image.Send(ctx, filename, fmt.Sprint(name))
	if err != nil {
		return imageUrl, fmt.Errorf("failed to send image into firebase: %w", err)
	}

	return imageUrl, nil
}
