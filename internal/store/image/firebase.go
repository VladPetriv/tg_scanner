package image

import (
	"context"
	"fmt"
	"io"
	"os"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/VladPetriv/tg_scanner/config"
)

type firebaseStore struct {
	cfg *config.Config
}

func NewFirebase(cfg *config.Config) Store {
	return &firebaseStore{cfg: cfg}
}

func (f firebaseStore) Send(ctx context.Context, path string, objectName string) (string, error) {
	defaultURL := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/default.jpg?alt=media",
		f.cfg.StorageBucket,
	)
	if objectName == "" || path == "" {
		return defaultURL, nil
	}

	config := &firebase.Config{
		ProjectID:     f.cfg.ProjectID,
		StorageBucket: f.cfg.StorageBucket,
	}

	opt := option.WithCredentialsFile(f.cfg.SecretPath)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return "", fmt.Errorf("create firebase app error: %w", err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf("get access to storage error: %w", err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return "", fmt.Errorf("get default bucket error: %w", err)
	}

	storageWriter := bucket.Object(objectName).NewWriter(ctx)

	image, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening image error: %w", err)
	}

	if _, err = io.Copy(storageWriter, image); err != nil {
		return "", fmt.Errorf("copying image to firebase storage error: %w", err)
	}

	if err = storageWriter.Close(); err != nil {
		return "", fmt.Errorf("closing firebase storage writer error: %w", err)
	}

	url := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", f.cfg.StorageBucket, objectName)

	err = os.Remove(path)
	if err != nil {
		return "", fmt.Errorf("delete image error: %w", err)
	}

	return url, nil
}
