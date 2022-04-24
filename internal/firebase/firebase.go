package firebase

import (
	"context"
	"fmt"
	"io"
	"os"

	firebase "firebase.google.com/go"
	"github.com/VladPetriv/tg_scanner/config"
	"google.golang.org/api/option"
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
		return "", fmt.Errorf("createing firebase app error: %w", err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf("createing firebase storage error: %w", err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return "", fmt.Errorf("getting default bucket error: %w", err)
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

	url := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s.jpg?alt=media", cfg.StorageBucket, objectName)
	return url, nil
}
