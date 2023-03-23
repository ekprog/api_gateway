package core

const (
	Success         = "success"
	ServerError     = "server_error"
	NotFound        = "not_found"
	ValidationError = "validation_error"
	Unauthorised    = "unauthorised"
)

type Status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type StatusResponse struct {
	Status Status `json:"status"`
}

type IdResponse struct {
	Status Status `json:"status"`
	Id     int32
}
