package image

import (
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/VladPetriv/tg_scanner/config"
)

type firebaseStore struct {
	cfg *config.Config
}

func NewFirebase(cfg *config.Config) Store {
	return &firebaseStore{
		cfg: cfg,
	}
}

func (f firebaseStore) Send(ctx context.Context, pathToImage string, objectName string) (string, error) {
	defaultURL := fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/default.jpg?alt=media",
		f.cfg.StorageBucket,
	)
	if objectName == "" || pathToImage == "" {
		return defaultURL, nil
	}

	bucket, err := f.createFirebaseBucket(ctx)
	if err != nil {
		return "", fmt.Errorf("create firebase storage bucket: %w", err)
	}

	storageWriter := bucket.Object(objectName).NewWriter(ctx)

	image, err := os.Open(pathToImage)
	if err != nil {
		return "", fmt.Errorf("open image: %w", err)
	}

	if _, err = io.Copy(storageWriter, image); err != nil {
		return "", fmt.Errorf("copy local image to storage writer: %w", err)
	}

	if err = storageWriter.Close(); err != nil {
		return "", fmt.Errorf("close storage writer: %w", err)
	}

	err = os.Remove(pathToImage)
	if err != nil {
		return "", fmt.Errorf("delete image: %w", err)
	}

	return fmt.Sprintf(
		"https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media",
		f.cfg.StorageBucket,
		objectName,
	), nil
}

func (f firebaseStore) createFirebaseBucket(ctx context.Context) (*storage.BucketHandle, error) {
	appConfig := &firebase.Config{
		ProjectID:     f.cfg.ProjectID,
		StorageBucket: f.cfg.StorageBucket,
	}

	opt := option.WithCredentialsFile(f.cfg.SecretPath)

	app, err := firebase.NewApp(ctx, appConfig, opt)
	if err != nil {
		return nil, fmt.Errorf("create firebase app: %w", err)
	}

	appStorage, err := app.Storage(ctx)
	if err != nil {
		return nil, fmt.Errorf("create new instance of storage: %w", err)
	}

	bucket, err := appStorage.DefaultBucket()
	if err != nil {
		return nil, fmt.Errorf("get storage bucket: %w", err)
	}

	return bucket, nil
}
