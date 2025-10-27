package helpers

type Response struct {
	Status     int               `json:"status"`
	Message    string            `json:"message"`
	Validation map[string]string `json:"validation"`
	Data       interface{}       `json:"data"`
}

type PaginatedResponse struct {
	Limit int           `json:"limit"`
	Page  int           `json:"page"`
	Total int           `json:"total"`
	List  []interface{} `json:"list"`
}

func NewResponse(status int, message string, validation map[string]string, data interface{}) Response {
	return Response{
		Status:     status,
		Message:    message,
		Validation: validation,
		Data:       data,
	}
}
