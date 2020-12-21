package integrations

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/triviy/parklakes-viberbot/infrastructure/logger"
)

// Header that passed to HTTP request
type Header struct {
	Name  string
	Value string
}

// SendPostRequest sends JSON HTTP request
func SendPostRequest(url string, request interface{}, response interface{}, headers ...Header) error {
	b, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "serialization of request failed")
	}
	logger.InfofDetailed(string(b), "---- Sending HTTP Request to %s", url)

	c := &http.Client{}
	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return errors.Wrapf(err, "creating request to %s failed", url)
	}
	r.Header.Add("Content-Type", "application/json")
	for _, h := range headers {
		r.Header.Add(h.Name, h.Value)
	}
	resp, err := c.Do(r)
	if err != nil {
		return errors.Wrapf(err, "sending request to %s failed", url)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return errors.Wrap(err, "deserialization of response failed")
	}
	logger.InfofDetailed(response, "---- Receiving HTTP response from %s", url)
	return nil
}
