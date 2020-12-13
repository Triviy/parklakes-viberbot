package handlers

type okResponse struct {
	Message string `json:"message"`
}

func createOkResponse() *okResponse {
	return &okResponse{"OK"}
}
