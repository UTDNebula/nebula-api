package responses

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    string `json:"error"`
}
