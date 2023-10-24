package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/tickrlytics/tickerlytics-backend/configs"
)

func OCR(content []byte, docID int64) (string, error) {
	payload := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(payload)
	writer, err := multipartWriter.CreateFormFile("file", fmt.Sprintf("%d.pdf", docID))
	if err != nil {
		return "", err
	}
	_, err = writer.Write(content)
	if err != nil {
		return "", err
	}
	err = multipartWriter.Close()
	if err != nil {
		return "", err
	}
	ocrURL := fmt.Sprintf("%s%d", configs.GetEnvWithKey("OCR_URL", "http://127.0.0.1:8080"), docID)
	req, err := http.NewRequest(http.MethodPost, ocrURL, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("accept", "application/json")

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	} else if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var resMAP map[string]string
	err = json.Unmarshal(body, &resMAP)
	if err != nil {
		return "", err
	}
	errorMSG, isError := resMAP["error"]
	if isError {
		return "", errors.New(errorMSG)
	}
	return resMAP["result"], nil
}
