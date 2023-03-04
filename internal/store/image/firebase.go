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

var _ Store = (*firebaseStore)(nil)

func NewFirebase(cfg *config.Config) *firebaseStore {
	return &firebaseStore{cfg: cfg}
}

func (f firebaseStore) Send(ctx context.Context, pathToImage string, objectName string) (string, error) {
	defaultURL := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/default.jpg?alt=media",
		f.cfg.StorageBucket,
	)
	if objectName == "" || pathToImage == "" {
		return defaultURL, nil
	}

	config := &firebase.Config{
		ProjectID:     f.cfg.ProjectID,
		StorageBucket: f.cfg.StorageBucket,
	}

	opt := option.WithCredentialsFile(f.cfg.SecretPath)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		return "", fmt.Errorf("create firebase app: %w", err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return "", fmt.Errorf("get access to app storage: %w", err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return "", fmt.Errorf("get default bucket: %w", err)
	}

	storageWriter := bucket.Object(objectName).NewWriter(ctx)

	image, err := os.Open(pathToImage)
	if err != nil {
		return "", fmt.Errorf("open image: %w", err)
	}

	if _, err = io.Copy(storageWriter, image); err != nil {
		return "", fmt.Errorf("copy image to firebase storage writer: %w", err)
	}

	if err = storageWriter.Close(); err != nil {
		return "", fmt.Errorf("close firebaser storage writer: %w", err)
	}

	url := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media",
		f.cfg.StorageBucket,
		objectName,
	)

	err = os.Remove(pathToImage)
	if err != nil {
		return "", fmt.Errorf("delete image: %w", err)
	}

	return url, nil
}
