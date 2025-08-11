package filestorage

import (
	"context"
	"log"
)

func FileStorageProvider(location string, localBasePath string) Storage {
	var storage Storage
	var err error

	switch location {
	case "s3":
		storage, err = NewS3StorageFromEnv(context.Background())
		if err != nil {
			log.Fatalf("failed init s3: %v", err)
		}
	default:
		storage, err = NewLocalStorage(localBasePath)
		if err != nil {
			log.Fatalf("failed init local storage: %v", err)
		}
	}

	return storage
}
