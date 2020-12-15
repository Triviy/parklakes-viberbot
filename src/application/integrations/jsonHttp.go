package integrations

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// SendPostRequest sends JSON HTTP request
func SendPostRequest(url string, request interface{}, response interface{}) error {
	logrus.WithField("request", &request).Infof("Request to %s", url)
	bytesRepresentation, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "serialization of request failed")
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return errors.Wrapf(err, "sending request to %s failed", url)
	}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return errors.Wrap(err, "deserialization of response failed")
	}
	logrus.WithField("request", &response).Infof("Response from %s", url)
	return nil
}