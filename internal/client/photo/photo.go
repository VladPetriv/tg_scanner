package photo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/gotd/td/tg"

	"github.com/VladPetriv/tg_scanner/internal/model"
	"github.com/VladPetriv/tg_scanner/internal/store/firebase"
	"github.com/VladPetriv/tg_scanner/pkg/config"
	cerrors "github.com/VladPetriv/tg_scanner/pkg/errors"
)

var Size int = 1024 * 1024
var ErrPhotoDataIsNil error = errors.New("photo data is nil")

func DecodePhoto(photo tg.UploadFileClass) (*model.Image, error) {
	if photo == nil {
		return nil, ErrPhotoDataIsNil
	}

	var img *model.Image

	encodedImage, err := json.Marshal(photo)
	if err != nil {
		return nil, &cerrors.CreateError{Name: "JSON", ErrorValue: err}
	}

	err = json.Unmarshal(encodedImage, &img)
	if err != nil {
		return nil, fmt.Errorf("unmarshal JSON error: %w", err)
	}

	return img, nil
}

func CreatePhoto(img *model.Image, name string) (string, error) {
	if img == nil {
		return "", ErrPhotoDataIsNil
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

func ProcessPhoto[T string | int](ctx context.Context, photoData tg.UploadFileClass, name T, cfg *config.Config) (string, error) {
	image, err := DecodePhoto(photoData)
	if err != nil {
		fmt.Println(err)
	}

	filename, err := CreatePhoto(image, fmt.Sprint(name))
	if err != nil {
		fmt.Println(err)
	}

	imageUrl, err := firebase.SendImageToStorage(ctx, cfg, filename, fmt.Sprint(name))
	if err != nil {
		return imageUrl, fmt.Errorf("send image to firebase error: %w", err)
	}

	return imageUrl, nil
}
