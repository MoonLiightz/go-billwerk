package optimize

type ErrorResponse struct {
	Code             int    `json:"code"`
	ErrorMessage     string `json:"error"`
	ErrorDescription string `json:"message"`
	HTTPReason       string `json:"http_reason"`
	HTTPStatus       int    `json:"http_status"`
	Path             string `json:"path"`
	Timestamp        string `json:"timestamp"`
	RequestID        string `json:"request_id"`
	TransactionError string `json:"transaction_error"`
}

func (e ErrorResponse) Error() string {
	message := e.ErrorMessage
	if e.ErrorDescription != "" {
		message += ": " + e.ErrorDescription
	}

	return message
}
