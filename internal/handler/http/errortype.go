package http

// ErrorResponce возврат ошибки
type ErrorResponce struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
