package firebase

import (
	"context"
	"fmt"
	"io"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/VladPetriv/tg_scanner/pkg/config"
	"github.com/VladPetriv/tg_scanner/pkg/utils"
)

func SendImageToStorage(ctx context.Context, cfg *config.Config, path string, objectName string) (string, error) {
	defaultUrl := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/default.jpg?alt=media", cfg.StorageBucket)
	if objectName == "" || path == "" {
		return defaultUrl, nil
	}

	config := &firebase.Config{
		ProjectID:     cfg.ProjectID,
		StorageBucket: cfg.StorageBucket,
	}

	opt := option.WithCredentialsFile(cfg.SecretPath)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return "", &utils.CreateError{Name: "firebase appication", ErrorValue: err}
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return "", &utils.CreateError{Name: "firebase storage", ErrorValue: err}
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return "", &utils.GettingError{Name: "default bucket", ErrorValue: err}
	}

	storageWriter := bucket.Object(objectName).NewWriter(ctx)

	image, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening image error: %w", err)
	}

	if _, err := io.Copy(storageWriter, image); err != nil {
		return "", fmt.Errorf("copying image to firebase storage error: %w", err)
	}

	if err := storageWriter.Close(); err != nil {
		return "", fmt.Errorf("closing firebase storage writer error: %w", err)
	}

	url := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", cfg.StorageBucket, objectName)

	err = os.Remove(path)
	if err != nil {
		return "", fmt.Errorf("error while deleting image")
	}

	return url, nil
}
