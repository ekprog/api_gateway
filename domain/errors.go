package domain

type Code string

type Status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func ServerError() Status {
	return Status{
		Code:    "server_error",
		Message: "Unknown Error. Try again later...",
	}
}

func UnauthorizedError() Status {
	return Status{
		Code:    "unauthorized",
		Message: "Login first...",
	}
}
