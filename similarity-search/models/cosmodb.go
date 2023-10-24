package models

import (
	"context"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/tickrlytics/tickerlytics-backend/configs"
)

var lock = &sync.Mutex{}

type azCosmos struct {
	client *azcosmos.Client
}
type CosmoContainer struct {
	client *azcosmos.ContainerClient
}

var azCosmosInstance *azCosmos

func createazCosmosClient(endpoint, key string) (*azcosmos.Client, error) {
	cred, err := azcosmos.NewKeyCredential(key)
	if err != nil {
		return nil, err
	}
	return azcosmos.NewClientWithKey(endpoint, cred, nil)
}
func GetCosmoDBInstance() (*azCosmos, error) {
	if azCosmosInstance != nil && azCosmosInstance.client != nil {
		return azCosmosInstance, nil
	}
	lock.Lock()
	defer lock.Unlock()
	var err error
	if azCosmosInstance != nil && azCosmosInstance.client != nil {
		return azCosmosInstance, nil
	}
	if azCosmosInstance == nil {
		azCosmosInstance = new(azCosmos)
	}
	azCosmosInstance.client, err = createazCosmosClient(configs.GetEnvWithKey("AZ_COSMOS_DB_ENDPOINT", ""), configs.GetEnvWithKey("AZ_COSMOS_DB_KEY", ""))
	return azCosmosInstance, err
}

func (azCosmos *azCosmos) UpsertItem(data []byte, itemId string, databaseID, containerID string) error {
	containerClient, err := azCosmos.client.NewContainer(databaseID, containerID)
	if err != nil {
		return err
	}
	_, err = containerClient.UpsertItem(context.Background(), azcosmos.NewPartitionKeyString(itemId), data, nil)
	return err
}

func (azCosmos *azCosmos) GetItem(itemId, databaseID, containerID string) ([]byte, error) {
	containerClient, err := azCosmos.client.NewContainer(databaseID, containerID)
	if err != nil {
		return nil, err
	}
	itemResponse, err := containerClient.ReadItem(context.Background(), azcosmos.NewPartitionKeyString(itemId), itemId, nil)
	if err != nil {
		return nil, err
	}
	return itemResponse.Value, nil
}
