package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/tickrlytics/tickerlytics-backend/configs"
)

type AzureBlob struct {
	Client *azblob.Client
}

var azureBlob *AzureBlob

func GetAzureBlob() (AzureBlob, error) {
	if azureBlob != nil && azureBlob.Client != nil {
		return *azureBlob, nil
	}
	azureBlob = new(AzureBlob)
	accountName := configs.GetEnvWithKey(configs.KEY_AZURE_STORAGE_ACCOUNT_NAME, "")
	accountKey := configs.GetEnvWithKey(configs.KEY_AZURE_STORAGE_ACCOUNT_KEY, "")
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return *azureBlob, err
	}
	azureBlob.Client, err = azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	return *azureBlob, err
}

func (azureBlob AzureBlob) UploadDoc(fileContent []byte, contentType, fileName, containerName string) (url string, err error) {
	uploadOption := &azblob.UploadBufferOptions{}
	_, err = azureBlob.Client.UploadBuffer(context.TODO(), containerName, fileName, fileContent, uploadOption)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s/%s", azureBlob.Client.URL(), containerName, fileName), nil
}
func (azureBlob AzureBlob) DeleteDoc(fileName, containerName string) {
	_, err := azureBlob.Client.DeleteBlob(context.Background(), containerName, fileName, &azblob.DeleteBlobOptions{})
	if err != nil {
		configs.Logger.Errorf("error in deleting file %s from %s storage container is %s", err.Error())
	}
	configs.Logger.Infof("%s file successfully deleted from %s storage container", fileName, containerName)
}
func (azureBlob AzureBlob) DownloadDoc(url string) []byte {
	//needs to update this function
	content, err := os.ReadFile(url)
	if err != nil {
		configs.Logger.Errorf("error in reading file:%s", err.Error())
	}
	return content
}

func DownloadFile(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
