package handlers

import (
	"encoding/json"

	"github.com/labstack/echo"
)

type okResponse struct {
	Message string `json:"message"`
}

func ok(r *echo.Response) error {
	r.Status = 200
	r.Header().Add(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	resp := okResponse{"OK"}
	message, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if _, err := r.Write([]byte(message)); err != nil {
		return err
	}
	return nil
}
