package adapters

import (
	"context"

	"gestrym-nutrition/src/nutrition/domain/interfaces"
	trainingInterfaces "gestrym-nutrition/src/nutrition/domain/interfaces"
)

type StorageServiceAdapterImpl struct {
	LegacyAdapter trainingInterfaces.FileStorageAdapter
}

func NewStorageServiceAdapterImpl(legacy trainingInterfaces.FileStorageAdapter) interfaces.StorageService {
	return &StorageServiceAdapterImpl{
		LegacyAdapter: legacy,
	}
}

func (a *StorageServiceAdapterImpl) UploadFromURL(ctx context.Context, imageURL string, fileName string) (string, error) {
	// We use the existing FileStorageAdapter to upload to our storage microservice
	collectionID, err := a.LegacyAdapter.UploadFromURL(imageURL, "nutrition")
	if err != nil {
		return "", err
	}
	return collectionID, nil
}
