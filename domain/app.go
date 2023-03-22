package domain

type Code string

const (
	Success string = "success"
)

type Status struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type StatusResponse struct {
	Status Status
}

func ServerError() Status {
	return Status{
		Code:    "server_error",
		Message: "Unknown Error. Try again later...",
	}
}

func ValidationError() Status {
	return Status{
		Code:    "validation_error",
		Message: "Validation error",
	}
}

func NotFoundError() Status {
	return Status{
		Code:    "not_found",
		Message: "Resource was not found",
	}
}

func UnauthorizedError() Status {
	return Status{
		Code:    "unauthorized",
		Message: "Login first...",
	}
}

type AccessRole int

const (
	RoleGuest      AccessRole = 0
	RoleUser       AccessRole = 1
	RoleSuperAdmin AccessRole = 10
)
