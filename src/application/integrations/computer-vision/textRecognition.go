package computervision

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/computervision"
	"github.com/Azure/go-autorest/autorest"
	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/infrastructure/logger"
)

const (
	numberOfCharsInoperationID = 36
	maxRetries                 = 10
)

// ImageTextReader wraper for Computer Vision API
type ImageTextReader struct {
	ctx    context.Context
	client computervision.BaseClient
}

// NewImageTextReader creates new instance of ImageTextReader
func NewImageTextReader(ctx context.Context, apiKey string, apiURL string) *ImageTextReader {
	computerVisionClient := computervision.New(apiURL)
	computerVisionClient.Authorizer = autorest.NewCognitiveServicesAuthorizer(apiKey)
	return &ImageTextReader{ctx, computerVisionClient}
}

// BatchReadFileRemoteImage reads text from image
func (r ImageTextReader) BatchReadFileRemoteImage(imageURL string) ([]string, error) {
	var remoteImage computervision.ImageURL
	remoteImage.URL = &imageURL

	textHeaders, err := r.client.BatchReadFile(r.ctx, remoteImage)
	if err != nil {
		return nil, errors.Wrap(err, "batch file reading failed")
	}
	operationLocation := autorest.ExtractHeaderValue("Operation-Location", textHeaders.Response)

	operationID := string(operationLocation[len(operationLocation)-numberOfCharsInoperationID:])

	readOperationResult, err := r.client.GetReadOperationResult(r.ctx, operationID)
	if err != nil {
		return nil, errors.Wrap(err, "getting read operation results failed")
	}

	i := 0

	for readOperationResult.Status != computervision.Failed &&
		readOperationResult.Status != computervision.Succeeded {
		if i >= maxRetries {
			break
		}
		i++

		logger.Info("Server status: %v, waiting %v seconds...\n", readOperationResult.Status, i)
		time.Sleep(1 * time.Second)

		readOperationResult, err = r.client.GetReadOperationResult(r.ctx, operationID)
		if err != nil {
			return nil, errors.Wrap(err, "getting read operation results failed")
		}
	}

	var results []string
	for _, recResult := range *(readOperationResult.RecognitionResults) {
		for _, line := range *recResult.Lines {
			results = append(results, *line.Text)
		}
	}
	logger.Infof("Computer Vision results for image %s are %v", imageURL, results)
	return results, nil
}
